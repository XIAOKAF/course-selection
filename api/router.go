package api

import (
	"github.com/gin-gonic/gin"
)

func InitEngine() {
	engine := gin.Default()
	engine.Use(Authorization)
	//高级管理员
	administratorGroup := engine.Group("/administrator")
	{
		administratorGroup.Use(parseToken)                    //解析token
		administratorGroup.POST("/login", administratorLogin) //管理员登录
		administratorGroup.POST("/remember", rememberStatus)  //记住登录状态
	}
	//学生
	//教师
	engine.Run(":8080")
}
