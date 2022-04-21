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
	teachingClass := ctx.PostForm("teachingClass")
	courseGrade := ctx.PostForm("courseGrade")
	duration := ctx.PostForm("duration")

	//必要字段为空
	if courseNumber == "" || courseName == "" || courseDepartment == "" || courseCredit == "" || courseType == "" || teachingClass == "" || courseGrade == "" || duration == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}

	//将课程编号和课程姓名以键值对的形式存储到redis之中
	err := service.HashSet("courseHash", courseNumber, courseName)
	tool.DealWithErr(ctx, err, "将课程编号和课程名称存入redis错误")

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
		courseDetails.Duration = result["duration"]
		courseDetailsArr = append(courseDetailsArr, courseDetails)
	}
	tool.Success(ctx, 200, courseDetailsArr)
}

func getSpecificCourse(ctx *gin.Context) {
	//模糊搜索
	keyWords := ctx.PostForm("keyWords")
	val := service.SScan("courseName", 0, "*"+keyWords+"*", 10)
	tool.Success(ctx, 200, val)
}

func chooseCourse(ctx *gin.Context) {
	//中间件验证请求头是否携带token且token存在并合格
	//从token中获取id方便后续插入数据的操作
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	tool.DealWithErr(ctx, err, "token解析出错")
	courseNumber := ctx.PostForm("courseNumber")
	teachingClass := ctx.PostForm("teachingClass")
	if courseNumber == "" || teachingClass == "" {
		tool.Failure(ctx, 400, "关键字段不能为空哦")
		return
	}
	//在redis课程编号集合里面查找该课程编号是否存在
	err, flag := service.SIsMember("courseNumber", courseNumber)
	tool.DealWithErr(ctx, err, "在redis中查询该课程编号出错")
	if !flag {
		tool.Failure(ctx, 400, "课程不存在")
		return
	}
	choice := model.Choice{
		TeachingClass: teachingClass,
		UnifiedCode:   tokenClaims.UserId,
	}
	//将学生选课信息存入MySQL
	err = service.ChooseCourse(choice)
	tool.DealWithErr(ctx, err, "将选课信息存入MySQL出错")
	//将学生选课信息存入redis
	err = service.SetAdd(courseNumber+teachingClass, tokenClaims.UserId)
	tool.DealWithErr(ctx, err, "将选课信息存入redis出错")
	tool.Success(ctx, 200, "选课成功")
}
