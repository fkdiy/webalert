package config

import (
	"fmt"
	"os"
)

func Load() {
	dat, err := os.ReadFile("config.yaml")

	if err != nil {
		fmt.Println("Error: Could not locate config.yaml")
		return
	}

	fmt.Print(string(dat))
}
