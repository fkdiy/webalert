package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fkdiy/webalert/internal/config"
	"github.com/fkdiy/webalert/internal/monitor"
)

func main() {
	log.SetFlags(0)

	// Locate and parse config.yaml
	conf, err := config.Load()

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
	}
}
