package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fkdiy/webalert/internal/config"
	"github.com/fkdiy/webalert/internal/email"
	"github.com/fkdiy/webalert/internal/jitter"
	"github.com/fkdiy/webalert/internal/monitor"
)

var version = "dev"

func main() {
	log.SetFlags(0)

	configPath := flag.String(
		"config",
		"webalert.config.yaml",
		"Path to the configuration file",
	)

	versionFlag := flag.Bool(
		"version",
		false,
		"Print version information",
	)

	flag.Parse()

	if *versionFlag {
		fmt.Printf("webalert %s\n", version)
		return
	}

	// Locate and parse config.yaml
	conf, err := config.Load(*configPath)

	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	mon := monitor.New(conf.Targets, version)

	// Set initial state
	_, err = mon.Check()

	if err != nil {
		log.Printf("Some targets could not be initialized:\n\n%v\n\n", err)
	}

	// Calculate interval with jitter
	interval := time.Duration(conf.Interval+jitter.GetJitter(conf.Jitter)) * time.Second

	// Initialize timer
	timer := time.NewTimer(interval)
	defer timer.Stop()

	fmt.Printf("Checking targets every %v seconds with %v seconds of jitter ...\n\n", conf.Interval, conf.Jitter)

	// Run infinite check loop
	for range timer.C {
		changes, err := mon.Check()

		if err != nil {
			log.Printf("Some targets could not be checked:\n\n%v\n\n", err)
		}

		if len(changes) == 0 {
			timer.Reset(time.Duration(conf.Interval+jitter.GetJitter(conf.Jitter)) * time.Second)
			continue
		}

		// Print change alert to console
		for _, change := range changes {
			fmt.Printf(
				"%v: Change registered in selector '%v'\nPrevious: %v\nCurrent: %v\nAlerting: %v\n\n",
				change.Target.URL,
				change.Target.Selector,
				change.Previous,
				change.Current,
				strings.Join(conf.EMail.Recipients, ", "),
			)
		}

		switch conf.EMail.Mode {
		case config.EMailModePerChange:
			// Send an e-mail for every change
			for _, change := range changes {
				err = email.Send(conf.EMail, change)

				if err != nil {
					log.Printf("Cant't send e-mail: %v\n\n", err)
				}
			}

		case config.EMailModeDigest:
			// Send one e-mail containing all changes
			err := email.SendDigest(conf.EMail, changes)

			if err != nil {
				log.Printf("Could not send email: %v", err)
			}
		}

		timer.Reset(time.Duration(conf.Interval+jitter.GetJitter(conf.Jitter)) * time.Second)
	}
}
