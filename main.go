package main

import (
	"learn_redis/backend/config"
	"learn_redis/backend/middleware"
	"learn_redis/backend/router"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()

	r := gin.Default()
	r.Use(middleware.GlobalInterceptor())
	router := router.Router{}
	// use reflect to improve
	router.RouteBlog(r)
	router.RouteShopType(r)
	router.RouteUser(r)

	r.Run(":8081")

}
