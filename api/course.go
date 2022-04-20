package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

//插入新的课程信息
func insertCourse(ctx *gin.Context) {
	//检验登录状态
	//确认登录状态
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		fmt.Println("token解析失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	_, flag, err := service.Get(tokenClaims.UserId)
	if err != nil {
		fmt.Println("查询统一验证码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if !flag {
		tool.Failure(ctx, 400, "token不存在")
		return
	}

	courseNumber := ctx.PostForm("courseNumber")
	courseName := ctx.PostForm("courseName")
	courseDepartment := ctx.PostForm("courseDepartment")
	courseCredit := ctx.PostForm("courseCredit")
	courseType := ctx.PostForm("courseType")
	teacher := ctx.PostForm("teacher")
	teachingClass := ctx.PostForm("teachingClass")
	courseGrade := ctx.PostForm("courseGrade")
	setTime := ctx.PostForm("setTime")
	duration := ctx.PostForm("duration")

	if courseNumber == "" || courseName == "" || courseDepartment == "" || courseCredit == "" || courseType == "" || teacher == "" || teachingClass == "" || courseGrade == "" || setTime == "" || duration == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}

	classCredit, err := strconv.ParseFloat(courseCredit, 32)
	tool.DealWithErr(ctx, err, "课程学分string转float64错误")
	classType, err := strconv.Atoi(courseType)
	tool.DealWithErr(ctx, err, "课程类型string转int错误")
	course := model.Course{
		CourseNumber:     courseNumber,
		CourseName:       courseName,
		CourseDepartment: courseDepartment,
		CourseCredit:     classCredit,
		CourseType:       classType,
		SetTime:          setTime,
		Duration:         duration,
	}

	//将课程信息放入MySQL
	err = service.InsertCourse(course)
	tool.DealWithErr(ctx, err, "课程信息存入MySQL出错")
	//将信息存入redis
	err = service.RInsertCourse(course)
	tool.DealWithErr(ctx, err, "课程信息存入redis出错")
}

func getCourseInfo(ctx *gin.Context) {

}
