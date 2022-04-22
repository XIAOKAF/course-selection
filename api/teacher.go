package api

import (
	"course-selection/service"
	"course-selection/tool"
	"github.com/gin-gonic/gin"
)

func teacherLogin(ctx *gin.Context) {
	workNumber := ctx.PostForm("workNumber")
	password := ctx.PostForm("pwd")
	if workNumber == "" || password == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	//查询该教师是否存在
	err, flag := service.HExists("teacher", workNumber)
	tool.DealWithErr(ctx, err, "查询教师工号是否存在错误")
	if !flag {
		tool.Failure(ctx, 400, "该教师不存在")
		return
	}
	//密码是否正确
	pwd, err := service.HashGet("teacher", workNumber)
	tool.DealWithErr(ctx, err, "查询教师密码错误")
	if password != pwd {
		tool.Failure(ctx, 400, "密码错误")
		return
	}
	//记住登录状态（24h
	err, token := service.CreateToken(workNumber, 2)
	tool.DealWithErr(ctx, err, "创建token错误")
	err = service.HashSet("token", workNumber, token)
	tool.DealWithErr(ctx, err, "存储token错误")
	tool.Success(ctx, 200, token)
}
