package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"runtime"
)

func Logging(c *gin.Context) {
	log.Printf("#%v", runtime.NumGoroutine())
	log.Printf("%v %v", c.Request.Method, c.Request.RequestURI)
	c.Next()
}

func (chat *Chat) Router() *gin.Engine {
	router := gin.New()
	router.Use(Logging)
	router.GET("/ws", chat.serveWs)
	router.POST("/", chat.createRoom)
	router.GET("/", chat.getRooms)
	router.GET("/online", chat.getOnline)
	return router
}
