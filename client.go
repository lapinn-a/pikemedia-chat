package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Client struct {
	conn     *websocket.Conn
	name     string
	toSocket chan string
	toHub    chan Message
}

type Message struct {
	client  Client
	message string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (client Client) writePump() {
	for message := range client.toSocket {
		err := client.conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (client Client) readPump() {
	for {
		_, p, err := client.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		client.toHub <- Message{client, string(p)}
	}
}

func (hub Hub) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	var client Client
	conn.SetCloseHandler(func(code int, text string) error {
		hub.unregister <- client
		return nil
	})

	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}

	client.conn = conn
	client.name = string(p)
	client.toSocket = make(chan string, 64)
	client.toHub = hub.broadcast
	hub.register <- client

	go client.writePump()
	go client.readPump()
}
