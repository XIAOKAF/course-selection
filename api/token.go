package api

import (
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
)

func parseToken(ctx *gin.Context) {
	token := ctx.GetHeader("token")
	if token == "" {
		tool.Failure(ctx, 401, "请先登录哦")
		ctx.Abort()
		return
	}
	tokenClaims, err := service.ParseToken(token)

	if err != nil {
		if err.Error() == "fail to parse token" {
			tool.Failure(ctx, 500, "服务器错误")
			log.Fatal("token解析失败", err)
			return
		}
		tool.Failure(ctx, 400, "请先登录哦")
		ctx.Abort()
		return
	}

	if tokenClaims.Variety == "refreshToken" {
		tool.Failure(ctx, 400, "token类型错误")
		ctx.Abort()
		return
	}

	//查询token是否存在
	tokenString, err := service.HashGet(tokenClaims.Identify, "token")
	if err != nil {
		if err == redis.Nil {
			fmt.Println("储存token可能出现错误")
			tool.Failure(ctx, 400, "token错误")
			ctx.Abort()
			return
		}
		tool.Failure(ctx, 500, "服务器错误")
		log.Fatal("查询token失败", err)
		return
	}

	if tokenString != token {
		tool.Failure(ctx, 400, "token错误")
		ctx.Abort()
		return
	}
}
