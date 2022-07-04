package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func Logging(c *gin.Context) {
	log.Printf("%v %v", c.Request.Method, c.Request.RequestURI)
	c.Next()
}

func (chat *Chat) Router() *gin.Engine {
	router := gin.New()
	router.Use(Logging)
	router.GET("/ws", chat.serveWs)
	router.GET("/", chat.getRooms)
	return router
}
