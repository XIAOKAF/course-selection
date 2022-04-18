package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("www..xyz.com")

func parseToken(ctx *gin.Context) {
	tokenClaims := &model.TokenClaims{}
	tokenString := ctx.Request.Header.Get("token")
	if tokenString == "" {
		tool.Success(ctx, 200, "先登录哦")
		ctx.Abort()
		return
	}
	token, err := jwt.ParseWithClaims(tokenString, tokenClaims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	if err != nil {
		fmt.Println("解析token错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		ctx.Abort()
		return
	}
	result, flag, err := service.Get(tokenClaims.UserId)
	if err != nil {
		fmt.Println("从redis中获取token错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		ctx.Abort()
		return
	}
	//未查询到请求头中携带的token所对应的用户id
	if !flag {
		tool.Failure(ctx, 400, "token已经永久过期，请重新登录")
		ctx.Abort()
		return
	}
	//请求头中携带的token与在redis中查询到的token不一致
	if result != tokenString {
		tool.Failure(ctx, 400, "token错误")
		ctx.Abort()
		return
	}
	//token已经过期，但是在redis之中还可以查询到该键值，证明可以进行刷新token的操作
	if !token.Valid {
		err, newToken := service.CreateToken(tokenClaims.UserId, 2)
		if err != nil {
			fmt.Println("刷新token错误", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		tool.Success(ctx, 200, "你的新token是:"+newToken)
		return
	}
	tool.Success(ctx, 200, tokenClaims.UserId+"你好")
}
