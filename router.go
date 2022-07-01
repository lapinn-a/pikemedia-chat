package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func Logging(c *gin.Context) {
	log.Printf("%v %v", c.Request.Method, c.Request.RequestURI)
	c.Next()
}

func (hub Hub) Router() *gin.Engine {
	router := gin.New()
	router.Use(Logging)
	router.GET("/ws", hub.serveWs)
	router.StaticFile("/", "./static/index.html")
	return router
}
