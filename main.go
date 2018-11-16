package main

import (
	"goWanlu/routes"
	"goWanlu/server"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func welcome(c *gin.Context) {
	c.String(http.StatusOK, "Hello from Wanlu!")
}

func main() {
	env := server.Inst()
	defer env.Drop()
	log.Println("Starting....")

	env.Gin.GET("/", welcome)
	routes.RegisterApiRoutes(env)
	routes.RegisterWebRoutes(env)
	env.Gin.Run(env.ListenAddr)
}
