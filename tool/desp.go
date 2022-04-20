package tool

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, code int, info interface{}) {
	ctx.JSON(code, gin.H{
		"code": code,
		"info": info,
	})
}

func Failure(ctx *gin.Context, code int, info interface{}) {
	ctx.JSON(code, gin.H{
		"code": code,
		"info": info,
	})
}

func DealWithErr(ctx *gin.Context, err error, info string) {
	if err != nil {
		fmt.Println(info, err)
		ctx.JSON(500, gin.H{
			"code": 500,
			"info": "服务器错误",
		})
		ctx.Abort()
	}
}
