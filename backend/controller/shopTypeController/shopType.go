package shopTypeController

import (
	"learn_redis/backend/common"
	"learn_redis/backend/config"
	"learn_redis/backend/model"

	"github.com/gin-gonic/gin"
)

func GetShopTypeList(ctx *gin.Context) {
	var shop_list []model.ShopType
	config.MySQL.Order("sort").Find(&shop_list)
	ctx.JSON(200, common.Result{Success: true, Data: shop_list})
}
