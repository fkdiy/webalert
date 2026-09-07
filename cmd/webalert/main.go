package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fkdiy/webalert/internal/config"
)

func main() {
	conf, err := config.Load()

	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}

	ticker := time.NewTicker(time.Duration(conf.Interval) * time.Second)

	defer ticker.Stop()

	fmt.Printf("Checking targets every %v seconds ...\n\n", conf.Interval)

	for range ticker.C {
		fmt.Println("Tick")
	}
}
