package router

import (
	"learn_redis/backend/controller/shopTypeController"

	"github.com/gin-gonic/gin"
)

func (Router) RouteShopType(r *gin.Engine) {
	g := r.Group("/shop-type")
	g.GET("/list", shopTypeController.GetShopTypeList)
}
