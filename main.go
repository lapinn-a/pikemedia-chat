package main

import (
	"log"
)

func main() {
	chat := NewChat()
	hub := NewHub()
	go hub.run()
	chat.rooms["Общая"] = hub

	err := chat.Router().Run(":81")
	if err != nil {
		log.Fatalf("FATAL: Error starting server: %s\n", err)
	}
}
