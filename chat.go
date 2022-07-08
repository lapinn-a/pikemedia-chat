package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"runtime/pprof"
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
	//go hub.run()
	go func() {
		labels := pprof.Labels("func", "run")
		pprof.Do(context.Background(), labels, func(_ context.Context) {
			hub.run()
		})
	}()
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
		room, ok := chat.rooms[c.Query("room")]

		if !ok {
			c.String(http.StatusNotFound, "Room not found")
			return
		}
		online := room.CountOnline()
		c.JSON(http.StatusOK, gin.H{"online": online})
	} else {
		online := make(map[string]int)
		acc := 0
		for key, value := range chat.rooms {
			roomOnline := value.CountOnline()
			online[key] = roomOnline
			acc += roomOnline
		}
		c.JSON(http.StatusOK, gin.H{"rooms": online, "overall": acc})
	}
}
