package main

import (
	"flag"
	"goWeb/routes"
	"goWeb/server"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func welcome(c *gin.Context) {
	c.String(http.StatusOK, "Hello from goWeb!")
}

// StartWorkers start starsWorker by goroutine.
func StartWorkers() {
	//go workers.StatsWorker()
}

// StartGin Start Gin Server
func StartGin() {
	server.SetConfig(configPath)
	env := server.Inst()
	defer env.Drop()

	log.Println("Starting....[...]:" + env.Port)
	log.Println("startGin " + env.Path)
	env.Gin.GET("/", welcome)
	routes.RegisterApiRoutes(env)
	routes.RegisterWebRoutes(env)
	env.Gin.Run(env.ListenAddr)
}

var configPath string

func init() {
	flag.StringVar(&configPath, "p", "/root/", "goWeb file path")
	flag.Parse()
}

func main() {
	StartWorkers()
	StartGin()
}
