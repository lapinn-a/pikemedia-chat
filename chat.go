package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
)

type Chat struct {
	rooms map[string]*Hub
}

func NewChat() *Chat {
	return &Chat{
		rooms: make(map[string]*Hub),
	}
}

func (chat *Chat) createRoom(c *gin.Context) {
	_, ok := chat.rooms[c.PostForm("room")]

	if ok {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	hub := NewHub()
	go hub.run()
	chat.rooms[c.PostForm("room")] = hub
	c.Redirect(http.StatusSeeOther, "/")
}

func (chat *Chat) getRooms(c *gin.Context) {
	tmpl, _ := template.ParseFiles("static/index.html")
	err := tmpl.Execute(c.Writer, chat.rooms)
	if err != nil {
		log.Println(err)
		return
	}
}

func (chat *Chat) getOnline(c *gin.Context) {
	if c.Query("room") != "" {
		online := chat.rooms[c.Query("room")].CountOnline()
		c.JSON(http.StatusOK, gin.H{"online": online})
	} else {
		online := make(map[string]int)
		acc := 0
		for key, value := range chat.rooms {
			online[key] = value.CountOnline()
		}
		c.JSON(http.StatusOK, gin.H{"rooms": online, "overall": acc})
	}
}