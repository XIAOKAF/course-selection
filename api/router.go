package api

import (
	"github.com/gin-gonic/gin"
)

func InitEngine() {
	engine := gin.Default()

	engine.Use(Authorization)

	//短信接口
	SmsGroup := engine.Group("/Sms")
	{
		SmsGroup.POST("/sendSms", sendSms)   //发送短信验证码
		SmsGroup.POST("/checkSms", checkSms) //校验短信验证码
	}

	//管理员接口
	administratorGroup := engine.Group("/administrator")
	{
		administratorGroup.Use(parseToken)                    //解析token
		administratorGroup.POST("/login", administratorLogin) //管理员登录
		administratorGroup.POST("/cancel", cancel)            //注销学生账号
	}

	//学生接口
	engine.POST("/studentRegister", studentRegister) //学生注册
	engine.POST("/studentLogin", studentLogin)       //密码登录
	studentGroup := engine.Group("/student")
	{
		studentGroup.Use(parseToken)                                 //解析token
		studentGroup.POST("/changePassword", changePwdByOldPwd)      //通过旧密码修改密码
		studentGroup.POST("/updateMobile", updateMobile)             //更新电话号码
		studentGroup.POST("/checkCodeForUpdate", checkCodeForUpdate) //更新电话号码时校验验证码
		studentGroup.POST("/updateAvatar", updateAvatar)             //更新头像
		studentGroup.GET("/selectInfo", selectInfo)                  //查询学生信息
		studentGroup.GET("/selection", selection)                    //学生查询自己的选课信息
	}

	//教师接口
	engine.POST("/teacherLogin", teacherLogin) //教师登录
	teacherGroup := engine.Group("/teacher")
	{
		teacherGroup.GET("/getTeachingClass", getTeachingClass)  //获取所有教学班
		teacherGroup.GET("/studentSelection", studentsSelection) //获取选择教学班的同学信息
	}

	//课程接口
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
