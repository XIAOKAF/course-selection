package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"database/sql"
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
		if err == sql.ErrNoRows {
			tool.Failure(ctx, 400, "账号错误")
			return
		}
		fmt.Println("查询管理员密码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if pwd != password {
		tool.Failure(ctx, 400, "密码居然错了┭┮﹏┭┮")
		return
	}

	err, token := service.CreateToken(administratorId, 200)
	if err != nil {
		fmt.Println("创建token失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet(administratorId, "token", token)
	if err != nil {
		fmt.Println("储存token失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	if auth == "" {
		//不授权长时间免密登录
		err, refreshToken := service.CreateToken(administratorId, 500)
		if err != nil {
			fmt.Println("创建refreshToken失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		err = service.HashSet(administratorId, "refreshToken", refreshToken)
		if err != nil {
			fmt.Println("储存refreshToken失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		tool.Success(ctx, 200, token)
		return
	}

	err, refreshToken := service.RememberStatus(administratorId, 1000)
	if err != nil {
		fmt.Println("创建refreshToken失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet(administratorId, "refreshToken", refreshToken)
	if err != nil {
		fmt.Println("储存refreshToken失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	tool.Success(ctx, 200, token)
}

//注销学生账号
func cancel(ctx *gin.Context) {
	studentId := ctx.PostForm("studentId")
	if studentId == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	//查询该学生是否存在
	_, err := service.HashGet(studentId, "mobile")
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "该学生不存在")
			return
		}
		fmt.Println("查询学生姓名失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//删除MySQL中的信息
	err = service.Cancel(studentId)
	if err != nil {
		fmt.Println("删除学生信息错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//删除redis中的信息
	err, keysArr := service.HKeys(studentId)
	if err != nil {
		fmt.Println("删除学生信息错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HDel(studentId, keysArr)
	if err != nil {
		fmt.Println("删除学生信息错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	tool.Success(ctx, http.StatusOK, "已经将该学生删除")
}

func inviteTeacher(ctx *gin.Context) {
	teacherNumber := ctx.PostForm("teacherNumber")
	teacherId := ctx.PostForm("teacherId")
	teacherName := ctx.PostForm("teacherName")
	if teacherName == "" || teacherId == "" || teacherNumber == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	err := service.HashSet(teacherNumber, "teacherName", teacherName)
	if err != nil {
		fmt.Println("存储教师信息失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet(teacherNumber, "teacherId", teacherId)
	if err != nil {
		fmt.Println("存储教师信息失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	tool.Success(ctx, 200, "successfully!")
}
