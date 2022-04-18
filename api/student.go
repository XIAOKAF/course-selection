package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
)

func studentRegister(ctx *gin.Context) {
	unifiedCode := ctx.PostForm("unifiedCode")
	if unifiedCode == "" {
		tool.Failure(ctx, 400, "还没有输入统一验证码哦")
		return
	}
	//查询该学生是否是本校学生，是则返回true，不是则返回false
	flag, err := service.SelectStudentByUnifiedCode(unifiedCode)
	if err != nil {
		fmt.Println("查询统一验证码错误", err)
		tool.Failure(ctx, 400, "服务器错误")
		return
	}
	if !flag {
		tool.Failure(ctx, 400, "是本校学生(⊙o⊙)吗？")
		return
	}
	tool.Success(ctx, 200, "亲爱的"+unifiedCode+"，你已经成功激活账户啦！o(*￣▽￣*)ブ")
}

func changePwdByOldPwd(ctx *gin.Context) {
	unifiedCode := ctx.PostForm("unifiedCode")
	oldPwd := ctx.PostForm("oldPwd")
	newPwd := ctx.PostForm("newPwd")
	if unifiedCode == "" {
		tool.Failure(ctx, 400, "统一验证码不能为空哦")
		return
	}
	if oldPwd == "" {
		tool.Failure(ctx, 400, "密码不能为空哦，悄悄提醒你，初始验证码为姓名拼音哦")
		return
	}
	if newPwd == "" {
		tool.Failure(ctx, 400, "你还要不要改密码了")
		return
	}
	result, flag, err := service.Get(unifiedCode)
	if err != nil {
		fmt.Println("查询统一验证码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if !flag {
		tool.Failure(ctx, 400, "该统一验证码不存在")
		return
	}
	if result != oldPwd {
		tool.Failure(ctx, 400, "原来的密码不正确哦")
		return
	}
	student := model.Student{
		UnifiedCode: unifiedCode,
		Password:    newPwd,
	}
	err = service.UpdatePassword(student)
	if err != nil {
		fmt.Println("跟新密码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	tool.Success(ctx, 200, "成功♪(^∇^*)")
}
