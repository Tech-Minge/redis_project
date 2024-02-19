package middleware

import (
	"context"
	"learn_redis/backend/common"
	"learn_redis/backend/config"
	"learn_redis/backend/service/userService"
	"log"

	"github.com/gin-gonic/gin"
)

func GlobalInterceptor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("authorization")
		tokenKey := common.LOGIN_TOKEN_PREFIX + token

		var user userService.UserDTO
		cmd := config.Redis.HGetAll(context.Background(), tokenKey)
		if mp, _ := cmd.Result(); len(mp) > 0 {
			if err := cmd.Scan(&user); err != nil {
				panic(err)
			}
			if _, exist := ctx.Get("User"); exist {
				panic("Unexpected")
			}
			ctx.Set("User", user)
			config.Redis.Expire(context.Background(), tokenKey, common.LOGIN_TOKEN_EXPIRE)
		}
		ctx.Next()
	}
}

func LoginInterceptor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, exist := ctx.Get("User"); !exist {
			log.Printf("用户未登录，拦截请求 %s", ctx.Request.URL)
			ctx.AbortWithStatusJSON(401, common.Result{ErrorMsg: "用户未登录"})
		} else {
			ctx.Next()
		}
	}
}
