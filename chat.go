package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
)

type Chat struct {
	rooms map[string]*Hub
}

func NewChat() *Chat {
	return &Chat{
		rooms: make(map[string]*Hub),
	}
}

func (chat *Chat) getRooms(c *gin.Context) {
	tmpl, _ := template.ParseFiles("static/index.html")
	err := tmpl.Execute(c.Writer, chat.rooms)
	if err != nil {
		log.Println(err)
		return
	}
}
