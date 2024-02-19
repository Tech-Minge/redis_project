package router

import (
	"learn_redis/backend/controller/blogController"
	"learn_redis/backend/middleware"

	"github.com/gin-gonic/gin"
)

func (Router) RouteBlog(r *gin.Engine) {
	g := r.Group("/blog")
	g.GET("/hot", blogController.GetHotBlogs)

	// requests below must login first
	g.Use(middleware.LoginInterceptor())
	g.GET("/of/me", blogController.GetMyBlogs)
}
