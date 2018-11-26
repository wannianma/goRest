package routes

import (
	"goWeb/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterWebRoutes(env *server.Env) {
	router := env.Gin
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "resources/static")

	router.GET("/nav", func(c *gin.Context) {
		c.HTML(http.StatusOK, "top-nav.tmpl.html", gin.H{
			"title": "Main website",
		})
	})
	// Default HTML page (client-side routing implemented via Vue.js)
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"title": "Main website",
		})
	})
}
