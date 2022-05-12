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
	for _, v := range teachingArr {
		teachingClassInfo.TeachingClassNumber = "class1"
		//课程id
		teachingClassInfo.CourseNumber = v
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
		teachingClassInfo.SetTime, err = service.HashGet(teachingClassInfo.TeachingClassNumber, "setTime")
		if err != nil {
			fmt.Println("查询教学班开设时间失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}

		//选课人数
		err, studentArr := service.HKeys(teachingClassInfo.TeachingClassNumber)
		if err != nil {
			fmt.Println("查询选课学生失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		teachingClassInfo.StudentSum = len(studentArr) - 3
		infoArr = append(infoArr, teachingClassInfo)
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
	//教师账号
	teacherId, err := service.HashGet(tokenClaims.Identify, "teacherId")
	if err != nil {
		fmt.Println("查询教师账号失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//课程编号
	_, err = service.HashGet(teacherId, teachingClass)
	if err != nil {
		fmt.Println("查询课程编号失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	var studentIdMap map[string]string
	studentIdMap = make(map[string]string)
	err, studentIdMap = service.HashGetAll(teachingClass)
	if err != nil {
		fmt.Println("查询选课学生失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	student := model.Student{}
	var studentArr []model.Student
	flag := true
	i := 0
	for k := range studentIdMap {
		if k[0:1] == "c" {
		} else {
			//学生姓名
			student.StudentName, err = service.HashGet(k, "studentName")
			if err != nil {
				fmt.Println("获取学生姓名错误", err)
				tool.Failure(ctx, 500, "服务器错误")
				flag = false
			}
			//学生性别
			student.Gender, err = service.HashGet(k, "gender")
			if err != nil {
				fmt.Println("获取学生性别错误", err)
				tool.Failure(ctx, 500, "服务器错误")
				flag = false
			}
			//学生院系
			student.Department, err = service.HashGet(k, "department")
			if err != nil {
				fmt.Println("获取学生院系错误", err)
				tool.Failure(ctx, 500, "服务器错误")
				flag = false
			}
			//学生专业
			student.Major, err = service.HashGet(k, "major")
			if err != nil {
				fmt.Println("获取学生专业错误", err)
				tool.Failure(ctx, 500, "服务器错误")
				flag = false
			}
			//学生年级
			student.Grade, err = service.HashGet(k, "grade")
			if err != nil {
				fmt.Println("查询学生年级失败", err)
				tool.Failure(ctx, 500, "服务器错误")
				flag = false
			}
			studentArr[i] = student
		}
		if !flag {
			break
		}
		i++
	}
	if !flag {
		return
	}

	tool.Success(ctx, http.StatusOK, studentArr)
}
