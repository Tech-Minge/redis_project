package blogController

import (
	"learn_redis/backend/service/blogService"
	"learn_redis/backend/service/userService"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetHotBlogs(ctx *gin.Context) {
	current, err := strconv.Atoi(ctx.DefaultQuery("current", "1"))
	if err != nil {
		panic(err)
	}
	ctx.JSON(200, blogService.GetHotBlogs(current))
}

func GetMyBlogs(ctx *gin.Context) {
	user := ctx.MustGet("User").(userService.UserDTO)
	current, err := strconv.Atoi(ctx.DefaultQuery("current", "1"))
	if err != nil {
		panic(err)
	}
	ctx.JSON(200, blogService.GetUserBlogs(user, current))
}
