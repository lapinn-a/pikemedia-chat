package main

import (
	"fmt"
	"log"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client, 128),
		unregister: make(chan *Client, 128),
		broadcast:  make(chan Message, 128),
	}
}

func (hub *Hub) CountOnline() int {
	return len(hub.clients)
}

func (hub *Hub) disconnect(client *Client) {
	_, ok := hub.clients[client]
	if ok {
		delete(hub.clients, client)
		close(client.toSocket)
	}
}

func (hub *Hub) run() {
	i := 0
	for {
		i++
		select {
		case res := <-hub.register:
			log.Printf("i:%d act:register %v", i, res)
			hub.clients[res] = true
		case res := <-hub.unregister:
			log.Printf("i:%d act:unregister %v", i, res)
			hub.disconnect(res)
		case res := <-hub.broadcast:
			log.Printf("i:%d act:broadcast %v", i, res)
			for client := range hub.clients {
				select {
				case client.toSocket <- func() string {
					if res.action == true {
						return fmt.Sprintf("[%s %s]", res.client.name, res.message)
					} else {
						return fmt.Sprintf("%s: %s", res.client.name, res.message)
					}
				}():
				default:
					log.Printf("Channel full. Disconnect")
					hub.disconnect(client)
				}
			}
		}
	}
}
