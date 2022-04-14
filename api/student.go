package api

import (
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
