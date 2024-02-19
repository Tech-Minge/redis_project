package userController

import (
	"learn_redis/backend/common"
	"learn_redis/backend/service/userService"

	"github.com/gin-gonic/gin"
)

type UserLoginInfo struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func AboutMe(ctx *gin.Context) {
	// pass login interceptor, so must get
	user := ctx.MustGet("User").(userService.UserDTO)
	ctx.JSON(200, common.Result{Success: true, Data: user})
}

func GetUserInfo(ctx *gin.Context) {
	ctx.JSON(200, common.Result{Success: true})
}

func SendCode(ctx *gin.Context) {
	phone := ctx.Query("phone")
	ctx.JSON(200, userService.SendCode(phone))
}

func Login(ctx *gin.Context) {
	var userInfo UserLoginInfo
	if err := ctx.ShouldBindJSON(&userInfo); err != nil {
		panic(err)
	}
	ctx.JSON(200, userService.Login(userInfo.Phone, userInfo.Code))
}
