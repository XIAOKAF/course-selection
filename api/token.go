package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
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
	if err != nil {
		fmt.Println("解析token失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	claims, ok := token.Claims.(*model.TokenClaims)
	if !ok {
		tool.Failure(ctx, 500, "服务器错误")
		log.Fatal("token解析失败")
		return
	}
	result, err := service.HashGet("token", claims.UserId)
	if err != nil {
		fmt.Println("查询token失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
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
		if err != nil {
			fmt.Println("删除token失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		//查询refreshToken
		result, err := service.HashGet("refreshToken", claims.UserId)
		if err != nil {
			if err == redis.Nil {
				tool.Failure(ctx, 400, "token已经失效")
				return
			}
			if err != nil {
				fmt.Println("查询refreshToken失败", err)
				tool.Failure(ctx, 500, "服务器错误")
				return
			}
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
			if err != nil {
				fmt.Println("删除refreshToken失败", err)
				tool.Failure(ctx, 500, "服务器错误")
				return
			}
			tool.Failure(ctx, 400, "请先登录")
		}
		//创建新的token
		err, newToken := service.CreateToken(claims.UserId, 2)
		if err != nil {
			fmt.Println("创建token失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		tool.Success(ctx, 200, newToken)
	}
	tool.Success(ctx, 200, tokenClaims.UserId+"你好")
}
