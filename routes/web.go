package routes

import (
	"goWeb/server"
	"html/template"
	"net/http"
	"time"

	"github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
)

func RegisterWebRoutes(env *server.Env) {
	router := env.Gin
	// router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "resources/static")

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

	//=========== Instant Message ===========//

	//new template middleware
	newYear := router.Group("/ny", gintemplate.NewMiddleware(gintemplate.TemplateConfig{
		Root:         "views",
		Extension:    ".tpl.html",
		Master:       "/layouts/blank",
		Partials:     []string{},
		DisableCache: true,
	}))

	newYear.GET("/", func(ctx *gin.Context) {
		// With the middleware, `HTML()` can detect the valid TemplateEngine.
		gintemplate.HTML(ctx, http.StatusOK, "happy", gin.H{
			"title": "Backend title!",
		})
	})

	// // Default HTML page (client-side routing implemented via Vue.js)
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusOK, "404 Not Found!")
	})
}
