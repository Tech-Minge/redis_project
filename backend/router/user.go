package router

import (
	"learn_redis/backend/controller/userController"
	"learn_redis/backend/middleware"

	"github.com/gin-gonic/gin"
)

func (Router) RouteUser(r *gin.Engine) {
	g := r.Group("/user")
	g.POST("/code", userController.SendCode)
	g.POST("/login", userController.Login)

	g.Use(middleware.LoginInterceptor())
	g.GET("/me", userController.AboutMe)
	g.GET("/info/:id", userController.GetUserInfo)

}
