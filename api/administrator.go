package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
)

//高级管理员登录
func administratorLogin(ctx *gin.Context) {
	administratorId := ctx.PostForm("administratorId")
	password := ctx.PostForm("pwd")
	auth := ctx.PostForm("auth")
	if administratorId == "" {
		tool.Failure(ctx, 400, "你冒充管理员啦！！！")
		return
	}
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
	if auth == "" {
		err, token := service.CreateToken(administratorId, 2)
		tool.DealWithErr(ctx, err, "创建token错误")
		err = service.HashSet("token", administratorId, token)
		tool.DealWithErr(ctx, err, "存储token错误")
		tool.Success(ctx, 200, token)
		return
	}
	err, token := service.RememberStatus(administratorId, 5)
	tool.DealWithErr(ctx, err, "创建token错误")
	err = service.HashSet("token", administratorId, token)
	tool.DealWithErr(ctx, err, "存储token错误")
	tool.Success(ctx, 200, token)
}

func cancel(ctx *gin.Context) {
	unifiedCode := ctx.PostForm("unifiedCode")
	if unifiedCode == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	//查询该学生是否存在
	_, err := service.HashGet(unifiedCode, "studentName")
	if err == redis.Nil {
		tool.Failure(ctx, 400, "该学生不存在")
		return
	}
	//删除MySQL中的信息
	err = service.Cancel(unifiedCode)
	tool.DealWithErr(ctx, err, "删除MySQL中的学生信息错误")
	//删除redis中的信息
	err, keysArr := service.HKeys(unifiedCode)
	tool.DealWithErr(ctx, err, "查询学生信息出错")
	err = service.HDel(unifiedCode, keysArr)
	tool.DealWithErr(ctx, err, "删除redis中的学生信息错误")
	tool.Success(ctx, http.StatusOK, "已经将该学生删除")
}
