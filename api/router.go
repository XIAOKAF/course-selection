package api

import (
	"github.com/gin-gonic/gin"
)

func InitEngine() {
	engine := gin.Default()

	//跨域
	engine.Use(Cors)
	//鉴权
	engine.Use(Authorization)

	//短信接口
	SmsGroup := engine.Group("/Sms")
	{
		SmsGroup.POST("/sendSms", sendSms)   //发送短信验证码
		SmsGroup.POST("/checkSms", checkSms) //校验短信验证码
	}

	//管理员接口
	engine.POST("/administratorLogin", administratorLogin) //管理员登录
	administratorGroup := engine.Group("/administrator")
	{
		administratorGroup.Use(parseToken)                             //解析token
		administratorGroup.POST("/spiderMan", spiderMan)               //导入学生信息
		administratorGroup.POST("/createCurriculum", createCurriculum) //开设新的课程
		administratorGroup.POST("/detailCurriculum", detailCurriculum) //开设教学班
		administratorGroup.DELETE("/cancel", cancel)                   //注销学生账号
	}

	//学生接口
	engine.POST("/studentRegister", studentRegister)   //学生注册
	engine.POST("/loginByStudentId", loginByStudentId) //密码登录
	studentGroup := engine.Group("/student")
	{
		studentGroup.Use(parseToken)                           //解析token
		studentGroup.POST("/changePwdByCode", changePwdByCode) //验证码修改密码
		studentGroup.POST("/updateAvatar", updateAvatar)       //更新头像
		studentGroup.GET("/selectInfo", selectInfo)            //查询学生信息
		studentGroup.GET("/getAvatar", getAvatar)              //获取学生头像
		studentGroup.POST("/chooseCourse", chooseCourse)       //选课
		studentGroup.GET("/selection", selection)              //学生查询自己的选课信息
		studentGroup.DELETE("/quit", quit)                     //学生退出班级
	}

	//教师接口
	engine.POST("/teacherLogin", teacherLogin) //教师登录
	teacherGroup := engine.Group("/teacher")
	{
		teacherGroup.Use(parseToken)
		teacherGroup.GET("/getTeachingClass", getTeachingClass)  //获取所有教学班
		teacherGroup.GET("/studentSelection", studentsSelection) //获取选择教学班的同学信息
	}

	//课程接口
	courseGroup := engine.Group("/course")
	{
		courseGroup.Use(parseToken)
		courseGroup.GET("/getAllCourse", getAllCourse)           //获取所有课程详情
		courseGroup.GET("/getSpecificCourse", getSpecificCourse) //模糊搜索
	}

	engine.Run(":8080")
}
