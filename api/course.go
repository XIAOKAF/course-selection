package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"github.com/gin-gonic/gin"
	"strconv"
)

//插入新的课程信息
func createCurriculum(ctx *gin.Context) {
	courseNumber := ctx.PostForm("courseNumber")
	courseName := ctx.PostForm("courseName")
	courseDepartment := ctx.PostForm("courseDepartment")
	courseCredit := ctx.PostForm("courseCredit")
	courseType := ctx.PostForm("courseType")
	courseGrade := ctx.PostForm("courseGrade")
	duration := ctx.PostForm("duration")

	//必要字段为空
	if courseNumber == "" || courseName == "" || courseDepartment == "" || courseCredit == "" || courseType == "" || courseGrade == "" || duration == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	//在redis中查询该课程是否已经存在
	//存在则不允许再创建
	//课程存在返回true，反之则false
	err, flag := service.SIsMember("course", courseNumber)
	tool.DealWithErr(ctx, err, "查询课程编号是否存在失败")
	if flag {
		tool.Failure(ctx, 400, "该课程已经存在")
		return
	}

	classCredit, err := strconv.ParseFloat(courseCredit, 32)
	tool.DealWithErr(ctx, err, "课程学分string转float64错误")
	//课程类型1表示选修，2表示必修
	classType, err := strconv.Atoi(courseType)
	tool.DealWithErr(ctx, err, "课程类型string转int错误")
	course := model.Course{
		CourseNumber:     courseNumber,
		CourseName:       courseName,
		CourseDepartment: courseDepartment,
		CourseCredit:     classCredit,
		CourseType:       classType,
		Duration:         duration,
	}

	//将课程信息放入MySQL
	err = service.CreateCourse(course)
	tool.DealWithErr(ctx, err, "课程信息存入MySQL出错")
	//将信息存入redis
	err = service.SetAdd("course", courseNumber)
	tool.DealWithErr(ctx, err, "课程信息存入redis失败")
	err = service.RCreateCourse(course)
	tool.DealWithErr(ctx, err, "课程信息存入redis出错")
	tool.Success(ctx, 200, "成功创建课程")
}

//开设教学班
//仅将数据插入redis
func detailCurriculum(ctx *gin.Context) {
	courseNumber := ctx.PostForm("courseNumber")
	teachingClass := ctx.PostForm("teachingClass")
	teacherNumber := ctx.PostForm("teacherNumber")
	setTime := ctx.PostForm("setTime")
	if courseNumber == "" || teachingClass == "" || teacherNumber == "" || setTime == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
	}
	//在redis中查询该课程是否存在
	err, flag := service.SelectCourse(courseNumber)
	tool.DealWithErr(ctx, err, "从redis中查询课程编号出错")
	if !flag {
		tool.Failure(ctx, 400, "该课程不存在")
		return
	}
	//在MySQL中查询该教师是否存在
	flag, err = service.SelectTeacher(teacherNumber)
	tool.DealWithErr(ctx, err, "从MySQL中查询教师编号出错")
	if !flag {
		tool.Failure(ctx, 400, "教师不存在")
		return
	}
	teaching := model.Teaching{
		CourseNumber:  courseNumber,
		TeachingClass: teachingClass,
		SetTime:       setTime,
	}
	//将数据存入redis
	err = service.SetAdd(teacherNumber, teachingClass)
	tool.DealWithErr(ctx, err, "将教学信息存入redis失败")
	err = service.RDetailsCourse(teaching)
	tool.DealWithErr(ctx, err, "将教学信息存入redis出错")
	tool.Success(ctx, 200, "教学信息设置成功")
}

//展示所有的课程信息
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
		courseDetails.Duration = result["duration"]
		courseDetailsArr = append(courseDetailsArr, courseDetails)
	}
	tool.Success(ctx, 200, courseDetailsArr)
}

//模糊搜索课程
func getSpecificCourse(ctx *gin.Context) {
	//模糊搜索
	keyWords := ctx.PostForm("keyWords")
	val := service.SScan("courseName", 0, "*"+keyWords+"*", 10)
	tool.Success(ctx, 200, val)
}
