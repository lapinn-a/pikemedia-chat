package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
)

type Client struct {
	conn     *websocket.Conn
	name     string
	toSocket chan string
	toHub    chan Message
}

type Message struct {
	client  *Client
	message string
	action  bool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (client *Client) writePump() {
	for message := range client.toSocket {
		err := client.conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (client *Client) readPump() {
	for {
		_, p, err := client.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		client.toHub <- Message{client, string(p), false}
	}
}

func (chat *Chat) serveWs(c *gin.Context) {
	log.Printf("serveWS start")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}

	hub, ok := chat.rooms[string(p)]

	if !ok {
		log.Printf("not ok")
		err = conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	var client Client
	conn.SetCloseHandler(func(code int, text string) error {
		hub.unregister <- &client
		hub.broadcast <- Message{&client, "выходит из чата", true}
		return nil
	})

	_, p, err = conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}

	client.conn = conn
	client.name = fmt.Sprintf("%s %d", string(p), rand.Int())
	client.toSocket = make(chan string, 64)
	client.toHub = hub.broadcast
	hub.register <- &client
	hub.broadcast <- Message{&client, "заходит в чат", true}

	go client.writePump()
	go client.readPump()
	log.Printf("serveWS end")
}
