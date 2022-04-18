package api

import (
	"github.com/gin-gonic/gin"
)

func InitEngine() {
	engine := gin.Default()
	engine.Use(Authorization)
	engine.POST("/checkSms", checkSms) //检验验证码是否正确
	//高级管理员
	administratorGroup := engine.Group("/administrator")
	{
		administratorGroup.Use(parseToken)                    //解析token
		administratorGroup.POST("/login", administratorLogin) //管理员登录
		administratorGroup.POST("/remember", rememberStatus)  //记住登录状态
	}
	//学生
	engine.POST("/studentRegister", studentRegister) //学生注册（类似于学生第一天报道后激活官方账号
	studentGroup := engine.Group("/student")
	{
		studentGroup.Use(parseToken)                    //解析token
		studentGroup.POST("loginByVerifyCode", sendSms) //短信登录
	}
	//教师
	engine.Run(":8080")
}
