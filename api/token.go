package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
)

var jwtKey = []byte("www..xyz.com")

func parseToken(ctx *gin.Context) {
	tokenClaims := &model.TokenClaims{}
	tokenString := ctx.Request.Header.Get("token")
	if tokenString == "" {
		tool.Failure(ctx, 200, "请先登录")
		ctx.Abort()
		return
	}
	token, err := jwt.ParseWithClaims(tokenString, tokenClaims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	tool.DealWithErr(ctx, err, "解析token错误")
	claims, ok := token.Claims.(*model.TokenClaims)
	if !ok {
		tool.Failure(ctx, 500, "服务器错误")
		log.Fatal("token解析失败")
		return
	}
	result, err := service.HashGet("token", claims.UserId)
	tool.DealWithErr(ctx, err, "查询token出错")
	if result != tokenString {
		tool.Failure(ctx, 400, "token错误")
		return
	}

	//token已经过期
	if !token.Valid {
		//删除redis中的token
		var filedName []string
		filedName = append(filedName, claims.UserId)
		err := service.HDel("token", filedName)
		tool.DealWithErr(ctx, err, "删除token出错")
		//查询refreshToken
		result, err := service.HashGet("refreshToken", claims.UserId)
		if err != nil {
			if err == redis.Nil {
				tool.Failure(ctx, 400, "token已经失效")
				return
			}
			tool.DealWithErr(ctx, err, "查询refreshToken出错")
		}
		//解析refreshToken
		token, err := jwt.ParseWithClaims(result, tokenClaims, func(token *jwt.Token) (i interface{}, err error) {
			return jwtKey, nil
		})
		//refreshToken过期
		if !token.Valid {
			//删除refreshToken
			var refreshToken []string
			refreshToken = append(refreshToken, result)
			err := service.HDel("refreshToken", refreshToken)
			tool.DealWithErr(ctx, err, "删除refreshToken出错")
			tool.Failure(ctx, 400, "请先登录")
		}
		//创建新的token
		err, newToken := service.CreateToken(claims.UserId, 2)
		tool.DealWithErr(ctx, err, "创建token出错")
		tool.Success(ctx, 200, newToken)
	}

	tool.Success(ctx, 200, tokenClaims.UserId+"你好")
}
