package main

import (
	"log"

	"tableflip.dev/buoy/pkg/commands"
)

func main() {
	if err := commands.New().Execute(); err != nil {
		log.Fatalf("error during command execution: %v", err)
	}

}
