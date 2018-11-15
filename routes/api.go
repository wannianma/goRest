package routes

import (
	"goWanlu/appenv"
	"goWanlu/handler/api/account"
)

func RegisterApiRoutes(env *appenv.Env) {
	router := env.Gin
	// JSON-REST API Version 1
	indoor := router.Group("/indoor/v1")
	{
		indoor.GET("login", account.Login)
	}
}
