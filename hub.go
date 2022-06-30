package main

import (
	"fmt"
)

type Hub struct {
	clients    map[Client]bool
	register   chan Client
	unregister chan Client
	broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[Client]bool),
		register:   make(chan Client),
		unregister: make(chan Client),
		broadcast:  make(chan Message),
	}
}

func (hub Hub) run() {
	for {
		select {
		case res := <-hub.register:
			fmt.Println("Response register ", res)
			hub.clients[res] = true
		case res := <-hub.unregister:
			fmt.Println("Response unregister ", res)
			delete(hub.clients, res)
		case res := <-hub.broadcast:
			fmt.Println("Response broadcast ", res)
			for client := range hub.clients {
				client.toSocket <- fmt.Sprintf("%s: %s", res.client.name, res.message)
			}
		}
	}
}
