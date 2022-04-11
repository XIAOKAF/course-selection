package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

//高级管理员登录
func administratorLogin(ctx *gin.Context) {
	administratorId := ctx.PostForm("administratorId")
	if administratorId == "" {
		tool.Failure(ctx, 400, "你冒充管理员啦！！！")
		return
	}
	password := ctx.PostForm("password")
	if password == "" {
		tool.Failure(ctx, 400, "密码怎么是空的呀(((φ(◎ロ◎;)φ)))")
		return
	}
	administrator := model.Administrator{
		AdministratorId: administratorId,
	}
	err, pwd := service.AdministratorLogin(administrator)
	if err != nil {
		fmt.Println("高级管理", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if pwd != password {
		tool.Failure(ctx, 400, "密码居然错了┭┮﹏┭┮")
		return
	}
}

func RememberStatus(ctx *gin.Context) {
	administratorId := ctx.PostForm("administratorId")
	auth := ctx.PostForm("auth")
	if auth == "" {
		return
	}
	a, err := strconv.Atoi(auth)
	if err != nil {
		fmt.Println("string转int错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//0表示拒绝记住登陆状态，除0之外表示同意7天内免密登录
	if a == 0 {
		return
	}
	//生成token
	err, token := service.CreateToken(administratorId, 2)
	if err != nil {
		fmt.Println("生成token错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//将token存到redis之中
	err = service.Set(administratorId, token, 2)
	if err != nil {
		fmt.Println("存储token错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//将token返回给前端
	tool.Success(ctx, 200, token)
	//中间件的形式检验并解析token
	//若当前token已经过期但在最长保质期内，刷新token，用状态码提示前端token已经刷新
	//反之则用状态码提示token已经永久失效，需要重新登录
}
