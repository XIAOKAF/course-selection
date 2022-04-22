package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
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
