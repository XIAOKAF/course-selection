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
	}
	//学生
	engine.POST("/studentRegister", studentRegister) //学生注册（类似于学生第一天报道后激活官方账号
	studentGroup := engine.Group("/student")
	{
		studentGroup.Use(parseToken)                                 //解析token
		studentGroup.POST("/loginByVerifyCode", sendSms)             //短信登录
		studentGroup.POST("/studentLogin", studentLogin)             //密码登录
		studentGroup.POST("/changePassword", changePwdByOldPwd)      //通过旧密码修改密码
		studentGroup.POST("/updateMobile", updateMobile)             //更新电话号码
		studentGroup.POST("/checkCodeForUpdate", checkCodeForUpdate) //更新电话号码时校验验证码
		studentGroup.POST("/updateAvatar", updateAvatar)             //更新头像
		studentGroup.GET("/selectInfo", selectInfo)                  //查询学生信息
	}
	//教师
	teacherGroup := engine.Group("/teacher")
	{
		teacherGroup.POST("login", teacherLogin) //教师登录
	}
	//课程
	courseGroup := engine.Group("/course")
	{
		courseGroup.POST("/insertCourse", createCurriculum)      //开设新的课程
		courseGroup.POST("/detailCourse", detailCurriculum)      //开设教学班
		courseGroup.GET("/getAllCourse", getAllCourse)           //获取所有课程详情
		courseGroup.GET("/getSpecificCourse", getSpecificCourse) //模糊搜索
		courseGroup.POST("/chooseCourse", chooseCourse)          //选课
	}
	engine.Run(":8080")
}
