package main

import (
	"log"

	"github.com/fkdiy/webalert/internal/config"
)

func main() {
	conf, err := config.Load()

	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}

}
