package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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
	if courseNumber[0:1] != "c" {
		tool.Failure(ctx, 400, "课程编号格式错误")
		return
	}

	//查询课程是否已经存在
	_, err := service.HashGet("course", courseNumber)
	if err != nil && err != redis.Nil {
		fmt.Println("查询课程编号失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if err == redis.Nil {
		//课程类型1表示选修，2表示必修
		classCredit, err := strconv.ParseFloat(courseCredit, 32)
		if err != nil {
			fmt.Println("课程编号转换数据类型错误", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		classType, err := strconv.Atoi(courseType)
		if err != nil {
			fmt.Println("课程类型转换数据类型失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}

		course := model.Course{
			CourseNumber:     courseNumber,
			CourseName:       courseName,
			CourseGrade:      courseGrade,
			CourseDepartment: courseDepartment,
			CourseCredit:     classCredit,
			CourseType:       classType,
			Duration:         duration,
		}

		//将课程信息放入MySQL
		err = service.CreateCourse(course)
		if err != nil {
			fmt.Println("将课程信息存入MySQL失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		//将信息存入redis
		err = service.HashSet("course", courseNumber, courseName)
		if err != nil {
			fmt.Println("将课程编号存入redis失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}

		err = service.HashSet(courseNumber, "courseName", course.CourseName)
		if err != nil {
			fmt.Println("将课程信息存入redis失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		err = service.HashSet(courseNumber, "courseGrade", course.CourseGrade)
		if err != nil {
			fmt.Println("将课程信息存入redis失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		err = service.HashSet(courseNumber, "courseDepartment", course.CourseDepartment)
		if err != nil {
			fmt.Println("将课程信息存入redis失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		variety := strconv.Itoa(course.CourseType)
		err = service.HashSet(courseNumber, "courseType", variety)
		if err != nil {
			fmt.Println("将课程信息存入redis失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		err = service.HashSet(courseNumber, "duration", course.Duration)
		if err != nil {
			fmt.Println("将课程信息存入redis失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		err = service.HashSet(courseNumber, "courseCredit", courseCredit)
		if err != nil {
			fmt.Println("将课程信息存入redis失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}

		tool.Success(ctx, 200, "成功创建课程")
		return
	}

	tool.Failure(ctx, 400, "课程已经创建")
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
		return
	}
	//查询课程是否存在
	_, err := service.HashGet("course", courseNumber)
	if err != nil && err != redis.Nil {
		fmt.Println("查询课程编号失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	//查询教师是否存在
	teacherId, err := service.HashGet(teacherNumber, "teacherId")
	if err != nil && err != redis.Nil {
		fmt.Println("查询教师工号失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	//查询教学班是否存在
	err, classArr := service.HVals("course")
	flag := service.IsClassExist(classArr, teachingClass)
	if flag {
		tool.Failure(ctx, 400, "教学班已经存在")
		return
	}

	//将数据存入redis
	err = service.HashSet(courseNumber, teachingClass, teacherNumber)
	if err != nil {
		fmt.Println("将教学班信息存入课程信息中失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet(teachingClass, "courseNumber", courseNumber)
	if err != nil {
		fmt.Println("将课程编号存入教学班信息失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet(teacherId, teachingClass, courseNumber)
	if err != nil {
		fmt.Println("将教学班信息存入教师信息失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet(teachingClass, "setTime", setTime)
	if err != nil {
		fmt.Println("将开课时间存入教学班信息失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet(teachingClass, "workNumber", teacherNumber)
	if err != nil {
		fmt.Println("将教师信息存入教学班信息失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	tool.Success(ctx, 200, "教学信息设置成功")
}

//展示所有的课程信息
func getAllCourse(ctx *gin.Context) {
	//获取课程编号
	err, courseArr := service.HKeys("course")
	if err != nil {
		fmt.Println("从redis中获取课程编号失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	var courseDetailsArr []map[string]string
	for _, v := range courseArr {
		//获取课程对应的所有教学班编号及对应教师编号
		err, teachingClassArr := service.HashGetAll(v)
		if err != nil {
			fmt.Println("查询课程详情失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		teachingClassArr["courseNumber"] = v
		courseDetailsArr = append(courseDetailsArr, teachingClassArr)
	}
	tool.Success(ctx, 200, courseDetailsArr)
}

/*
//模糊搜索课程
func getSpecificCourse(ctx *gin.Context) {
	//模糊搜索
	keyWords := ctx.PostForm("keyWords")
	val := service.SScan("courseName", 0, "*"+keyWords+"*", 10)
	tool.Success(ctx, 200, val)
}
*/
