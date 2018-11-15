package routes

import (
	"goWanlu/appenv"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterWebRoutes(env *appenv.Env) {
	router := env.Gin
	router.LoadHTMLGlob("templates/*")

	// Default HTML page (client-side routing implemented via Vue.js)
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})
}
