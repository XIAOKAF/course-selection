package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func teacherLogin(ctx *gin.Context) {
	workNumber := ctx.PostForm("workNumber")
	password := ctx.PostForm("pwd")
	auth := ctx.PostForm("auth")
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
	if auth == "" {
		//记住登录状态（24h
		err, token := service.CreateToken(workNumber, 2)
		tool.DealWithErr(ctx, err, "创建token错误")
		err = service.HashSet("token", workNumber, token)
		tool.DealWithErr(ctx, err, "存储token错误")
		tool.Success(ctx, 200, token)
		return
	}
	err, token := service.RememberStatus(workNumber, 5)
	tool.DealWithErr(ctx, err, "创建token错误")
	err = service.HashSet("token", workNumber, token)
	tool.DealWithErr(ctx, err, "存储token错误")
	tool.Success(ctx, 200, token)
}

//获取老师所带的教学班
func getTeachingClass(ctx *gin.Context) {
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		fmt.Println("token解析失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err, classArr := service.SetGet(tokenClaims.UserId)
	tool.DealWithErr(ctx, err, "查询该教师所带教学班错误")
	teaching := model.Teaching{}
	var teachingArr []model.Teaching
	for i, v := range classArr {
		//获取教学班编号对应的课程编号
		teaching.CourseNumber, err = service.HashGet(tokenClaims.UserId+"teaching", v)
		tool.DealWithErr(ctx, err, "查询课程编号错误")
		//获取教学班开设的时间
		teaching.SetTime, err = service.HashGet(v+"teaching", "setTime")
		tool.DealWithErr(ctx, err, "获取教学班那开设时间错误")
		teachingArr[i] = teaching
	}
	tool.Success(ctx, http.StatusOK, teachingArr)
}

//获取选课学生信息
func studentsSelection(ctx *gin.Context) {
	tokenString := ctx.Request.Header.Get("token")
	teachingClass := ctx.PostForm("teachingClass")
	if teachingClass == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		fmt.Println("token解析失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err, flag := service.SIsMember(tokenClaims.UserId, teachingClass)
	tool.DealWithErr(ctx, err, "检索老师是否教授此班错误")
	if !flag {
		tool.Failure(ctx, 400, "保留彼此的空间，不打扰是你的温柔")
		return
	}
	//获取选课学生
	err, codeArr := service.HKeys(teachingClass)
	tool.DealWithErr(ctx, err, "查询所有选课学生出错")
	student := model.Student{}
	var studentArr []model.Student
	for i, v := range codeArr {
		student.StudentName, err = service.HashGet(v, "studentName")
		tool.DealWithErr(ctx, err, "获取学生姓名错误")
		sex, err := service.HashGet(v, "gender")
		tool.DealWithErr(ctx, err, "获取学生性别错误")
		student.Gender, err = strconv.Atoi(sex)
		tool.DealWithErr(ctx, err, "string转int错误")
		student.Department, err = service.HashGet(v, "department")
		tool.DealWithErr(ctx, err, "获取学生院系错误")
		student.Major, err = service.HashGet(v, "major")
		tool.DealWithErr(ctx, err, "获取学生专业错误")
		studentArr[i] = student
	}
	tool.Success(ctx, http.StatusOK, studentArr)
}
