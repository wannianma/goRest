package routes

import (
	"goWanlu/handler/api/indoor"
	"goWanlu/server"
)

func RegisterApiRoutes(env *server.Env) {
	router := env.Gin
	// JSON-REST API Version 1
	v1 := router.Group("/v1")
	{
		v1.GET("login", account.Login)
		v1.POST("register", account.Register)
	}
}
