package main

import (
	"log"
)

func main() {
	hub := NewHub()
	go hub.run()

	err := hub.Router().Run(":81")
	if err != nil {
		log.Fatalf("FATAL: Error starting server: %s\n", err)
	}
}
