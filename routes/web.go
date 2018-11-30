package routes

import (
	"goWeb/server"
	"net/http"
	"time"

	"html/template"

	"github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
)

func RegisterWebRoutes(env *server.Env) {
	router := env.Gin
	//new template engine
	router.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:      "views",
		Extension: ".tpl.html",
		Master:    "layouts/master",
		// Partials:  []string{"partials/ad"},
		Funcs: template.FuncMap{
			"copy": func() string {
				return time.Now().Format("2006")
			},
		},
		DisableCache: true,
	})
	// router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "resources/static")

	router.GET("/nav", func(c *gin.Context) {
		c.HTML(http.StatusOK, "top-nav.tmpl.html", gin.H{
			"title": "Main website",
		})
	})
	// Default HTML page (client-side routing implemented via Vue.js)
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{
			"title": "Main website",
		})
	})
}
