package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fkdiy/webalert/internal/config"
	"github.com/fkdiy/webalert/internal/monitor"
	email "github.com/fkdiy/webalert/internal/notifier"
)

func main() {
	log.SetFlags(0)

	configPath := flag.String(
		"config",
		"webalert.config.yaml",
		"Path to the configuration file",
	)

	flag.Parse()

	// Locate and parse config.yaml
	conf, err := config.Load(*configPath)

	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}

	mon := monitor.New(conf.Targets)

	// Set initial state
	_, err = mon.Check()

	if err != nil {
		log.Printf("Some targets could not be initialized:\n\n%v\n\n", err)
	}

	// Initialize ticker with interval set in config
	ticker := time.NewTicker(time.Duration(conf.Interval) * time.Second)
	defer ticker.Stop()

	fmt.Printf("Checking targets every %v seconds ...\n\n", conf.Interval)

	// Run infinite check loop
	for range ticker.C {
		changes, err := mon.Check()

		if err != nil {
			log.Printf("Some targets could not be checked:\n\n%v\n\n", err)
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
	}
}
