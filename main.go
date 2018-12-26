package main

import (
	"goWeb/routes"
	"goWeb/server"
	"goWeb/workers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func welcome(c *gin.Context) {
	c.String(http.StatusOK, "Hello from goWeb!")
}

// StartWorkers start starsWorker by goroutine.
func StartWorkers() {
	go workers.StatsWorker()
}

// StartGin Start Gin Server
func StartGin() {
	gin.SetMode(gin.ReleaseMode)

	env := server.Inst()
	defer env.Drop()

	log.Println("Starting....[...]:" + env.Port)

	env.Gin.GET("/", welcome)
	routes.RegisterApiRoutes(env)
	routes.RegisterWebRoutes(env)
	env.Gin.Run(env.ListenAddr)
}

func main() {
	StartWorkers()
	StartGin()
}
