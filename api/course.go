package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"github.com/gin-gonic/gin"
	"strconv"
)

//插入新的课程信息
func insertCourse(ctx *gin.Context) {
	courseNumber := ctx.PostForm("courseNumber")
	//课程编号单独存入redis
	err := service.SetAdd("course", courseNumber)
	tool.DealWithErr(ctx, err, "插入课程编号出错")
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

func getAllCourse(ctx *gin.Context) {
	err, members := service.SetGet("course")
	tool.DealWithErr(ctx, err, "从redis中获取课程编号错误")
	var courseDetails model.Course
	var courseDetailsArr []model.Course
	for _, val := range members {
		err, result := service.HashGetAll(val)
		tool.DealWithErr(ctx, err, "查询课程详情出错")
		courseDetails.CourseNumber = result["courseNumber"]
		courseDetails.CourseName = result["courseName"]
		courseDetails.CourseDepartment = result["courseDepartment"]
		credit, err := strconv.ParseFloat(result["courseCredit"], 64)
		tool.DealWithErr(ctx, err, "学分string转为float出错")
		courseDetails.CourseCredit = credit
		classType, err := strconv.Atoi(result["courseType"])
		tool.DealWithErr(ctx, err, "课程类型string转为int出错")
		courseDetails.CourseType = classType
		courseDetails.SetTime = result["setTime"]
		courseDetails.Duration = result["duration"]
		courseDetailsArr = append(courseDetailsArr, courseDetails)
	}
	tool.Success(ctx, 200, courseDetailsArr)
}
