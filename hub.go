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
		register:   make(chan Client, 64),
		unregister: make(chan Client, 64),
		broadcast:  make(chan Message, 64),
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
				select {
				case client.toSocket <- fmt.Sprintf("%s: %s", res.client.name, res.message):
				default:
					fmt.Println("Channel full. Disconnecting")
					hub.unregister <- client
				}
			}
		}
	}
}
