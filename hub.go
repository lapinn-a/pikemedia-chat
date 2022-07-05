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
		register:   make(chan *Client, 64),
		unregister: make(chan *Client, 64),
		broadcast:  make(chan Message, 64),
	}
}

func (hub *Hub) CountOnline() int {
	return len(hub.clients)
}

func (hub *Hub) run() {
	i := 0
	for {
		i++
		select {
		case res := <-hub.register:
			log.Printf("i:%d act:register %v", i, res)
			//fmt.Println("Response register ", res)
			hub.clients[res] = true
			//select {
			//case hub.broadcast <- Message{res, "заходит в чат", true}:
			//default:
			//}
		case res := <-hub.unregister:
			log.Printf("i:%d act:unregister %v", i, res)
			//fmt.Println("Response unregister ", res)
			//close(res.toSocket)
			//delete(hub.clients, res)
			//hub.broadcast <- Message{res, "выходит из чата", true}
		case res := <-hub.broadcast:
			log.Printf("i:%d act:broadcast %v", i, res)
			//fmt.Println("Response broadcast ", res)
			for client := range hub.clients {
				log.Printf("    %v", client)
				select {
				case client.toSocket <- func() string {
					if res.action == true {
						return fmt.Sprintf("[%s %s]", res.client.name, res.message)
					} else {
						return fmt.Sprintf("%s: %s", res.client.name, res.message)
					}
				}():
				default:
					//fmt.Println("Channel full. Disconnecting")
					//hub.unregister <- client
					close(client.toSocket)
					delete(hub.clients, client)
					select {
					case hub.broadcast <- Message{client, "выходит из чата", true}:
					default:
					}
				}
			}
		}
	}
}
