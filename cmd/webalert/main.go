package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fkdiy/webalert/internal/config"
	"github.com/fkdiy/webalert/internal/monitor"
)

func main() {
	// Locate and parse config.yaml
	conf, err := config.Load()

	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}

	mon := monitor.New(conf.Targets)

	// Set initial state
	mon.Check()

	ticker := time.NewTicker(time.Duration(conf.Interval) * time.Second)
	defer ticker.Stop()

	fmt.Printf("Checking targets every %v seconds ...\n\n", conf.Interval)

	for range ticker.C {
		mon.Check()
	}
}
