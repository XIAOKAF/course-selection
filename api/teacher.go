package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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
	_, err := service.HashGet("teacher", workNumber)
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "教师工号不存在")
			return
		}
		fmt.Println("查询教师是否存在失败")
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//密码是否正确
	pwd, err := service.HashGet(workNumber, "password")
	if err != nil {
		fmt.Println("查询教师密码失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if password != pwd {
		tool.Failure(ctx, 400, "密码错误")
		return
	}
	err, token := service.CreateToken(workNumber, 200)
	if err != nil {
		fmt.Println("创建token失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet(workNumber, "token", token)
	if err != nil {
		fmt.Println("储存token失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	if auth == "" {
		//不授权长时间免密登录
		err, refreshToken := service.CreateToken(workNumber, 500)
		if err != nil {
			fmt.Println("创建refreshToken失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		err = service.HashSet(workNumber, "refreshToken", refreshToken)
		if err != nil {
			fmt.Println("储存refreshToken失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		tool.Success(ctx, 200, token)
		return
	}

	err, refreshToken := service.RememberStatus(workNumber, 1000)
	if err != nil {
		fmt.Println("创建refreshToken失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet(workNumber, "refreshToken", refreshToken)
	if err != nil {
		fmt.Println("储存refreshToken失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

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
	teacherId, err := service.HashGet("teacher", tokenClaims.Identify)
	if err != nil {
		fmt.Println("查询教师账号失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err, teachingArr := service.HVals(teacherId)
	if err != nil {
		fmt.Println("查询教师所带教学班失败")
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	var teachingClassInfo model.TeachingClassInfo
	var infoArr []model.TeachingClassInfo
	for i, v := range teachingArr {
		teachingClassInfo.TeachingClassNumber = v
		//课程id
		teachingClassInfo.Course.CourseNumber, err = service.HashGet(v, "courseNumber")
		if err != nil {
			fmt.Println("查询课程编号失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		//课程名称
		teachingClassInfo.Course.CourseName, err = service.HashGet(teachingClassInfo.CourseNumber, "courseName")
		if err != nil {
			fmt.Println("查询课程名称错误", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		//课程类型
		variety, err := service.HashGet(teachingClassInfo.CourseNumber, "courseType")
		if err != nil {
			fmt.Println("查询课程类型失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		teachingClassInfo.CourseType, err = strconv.Atoi(variety)
		if err != nil {
			fmt.Println("类型转换失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		//课程院系
		teachingClassInfo.CourseDepartment, err = service.HashGet(teachingClassInfo.CourseNumber, "courseDepartment")
		if err != nil {
			fmt.Println("查询课程所属院系失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		//课程年级
		teachingClassInfo.CourseGrade, err = service.HashGet(teachingClassInfo.CourseNumber, "courseGrade")
		if err != nil {
			fmt.Println("查询课程所属年级失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		//课程时间
		teachingClassInfo.Duration, err = service.HashGet(teachingClassInfo.CourseNumber, "duration")
		if err != nil {
			fmt.Println("查询课程时长失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		//教学班开设时间
		teachingClassInfo.SetTime, err = service.HashGet(v, "setTime")
		if err != nil {
			fmt.Println("查询教学班开设时间失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		//选课人数
		err, studentArr := service.HKeys(v)
		if err != nil {
			fmt.Println("查询选课学生失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		teachingClassInfo.StudentSum = len(studentArr) - 1
		infoArr[i] = teachingClassInfo
	}
	tool.Success(ctx, http.StatusOK, infoArr)
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
	if err != nil {
		fmt.Println("检索教师是否教该教学班失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if !flag {
		tool.Failure(ctx, 400, "保留彼此的空间，不打扰是你的温柔")
		return
	}
	//获取选课学生
	err, codeArr := service.HKeys(teachingClass)
	if err != nil {
		fmt.Println("查询所有选课学生失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	student := model.Student{}
	var studentArr []model.Student
	for i, v := range codeArr {
		student.StudentName, err = service.HashGet(v, "studentName")
		if err != nil {
			fmt.Println("获取学生姓名错误", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		sex, err := service.HashGet(v, "gender")
		if err != nil {
			fmt.Println("获取学生性别错误", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		student.Gender, err = strconv.Atoi(sex)
		if err != nil {
			fmt.Println("string转int错误", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		student.Department, err = service.HashGet(v, "department")
		if err != nil {
			fmt.Println("获取学生院系错误", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		student.Major, err = service.HashGet(v, "major")
		if err != nil {
			fmt.Println("获取学生专业错误", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		studentArr[i] = student
	}
	tool.Success(ctx, http.StatusOK, studentArr)
}
