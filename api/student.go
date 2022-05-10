package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"strings"
)

func studentRegister(ctx *gin.Context) {
	studentId := ctx.PostForm("userId")
	password := ctx.PostForm("pwd")
	if studentId == "" || password == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	//查询该学生是否是本校学生，是则返回true，不是则返回false
	mobile, err := service.HashGet("student", studentId)
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "学号错误")
			return
		}
		fmt.Println("查询学生是否存在错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	pwd, err := service.HashGet(studentId, "password")
	if err != nil {
		fmt.Println("查询密码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if pwd != password {
		tool.Failure(ctx, 400, "密码错误（提示一下哦，初始密码是姓名拼音")
		return
	}
	err = service.HashSet(mobile, "studentId", studentId)
	if err != nil {
		fmt.Println("储存学生信息失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	tool.Success(ctx, 200, "注册成功")
}

//学号密码登录
func loginByStudentId(ctx *gin.Context) {
	studentId := ctx.PostForm("studentId")
	password := ctx.PostForm("pwd")
	auth := ctx.PostForm("auth")
	if studentId == "" || password == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	pwd, err := service.HashGet(studentId, "password")
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "学号错误")
			return
		}
		fmt.Println("查询密码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	if password != pwd {
		tool.Failure(ctx, 400, "密码错误")
		return
	}

	err, token := service.CreateToken(studentId, 200)
	if err != nil {
		fmt.Println("创建token失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet(studentId, "token", token)
	if err != nil {
		fmt.Println("储存token失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	if auth == "" {
		//不授权长时间免密登录
		err, refreshToken := service.CreateToken(studentId, 500)
		if err != nil {
			fmt.Println("创建refreshToken失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		err = service.HashSet(studentId, "refreshToken", refreshToken)
		if err != nil {
			fmt.Println("储存refreshToken失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		tool.Success(ctx, 200, token)
		return
	}

	err, refreshToken := service.RememberStatus(studentId, 1000)
	if err != nil {
		fmt.Println("创建refreshToken失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet(studentId, "refreshToken", refreshToken)
	if err != nil {
		fmt.Println("储存refreshToken失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	tool.Success(ctx, 200, token)
}

//短信验证码修改密码
func changePwdByCode(ctx *gin.Context) {
	mobile := ctx.PostForm("mobile")
	code := ctx.PostForm("code")
	pwd := ctx.PostForm("newPwd")
	confirmPwd := ctx.PostForm("confirmPwd")
	if mobile == "" || code == "" || pwd == "" || confirmPwd == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	if pwd != confirmPwd {
		tool.Failure(ctx, 400, "两次密码输入不一致")
		return
	}
	rightCode, err := service.Get(mobile + "code")
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "验证码已过期")
			return
		}
		fmt.Println("查询验证码失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if rightCode != code {
		tool.Failure(ctx, 400, "验证码错误")
		return
	}

	//redis更新
	studentId, err := service.HashGet(mobile, "studentId")
	if err != nil {
		fmt.Println("查询学号失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet(studentId, "password", pwd)
	if err != nil {
		fmt.Println("重置密码失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	tool.Success(ctx, 200, "密码更新成功♪(^∇^*)")
}

//更新头像
func updateAvatar(ctx *gin.Context) {
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		fmt.Println("解析token失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	bucket, err := service.ParseBucket()
	if err != nil {
		fmt.Println("解析储存桶配置文件错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	u, err := url.Parse(bucket.Url)
	if err != nil {
		fmt.Println("解析url错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  bucket.SecretId,
			SecretKey: bucket.SecretKey,
		},
	})
	filePath := ctx.PostForm("filePath")
	if filePath == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	_, err = c.Object.PutFromFile(ctx, tokenClaims.Identify, filePath, nil)
	if err != nil {
		fmt.Println("上传头像失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	tool.Success(ctx, 200, "头像更新成功")
}

//查询个人信息
func selectInfo(ctx *gin.Context) {
	//确认登录状态
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		fmt.Println("token解析失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	//从redis中获取具体的信息
	studentName, err := service.HashGet(tokenClaims.Identify, "studentName")
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "token错误")
			return
		}
		fmt.Println("获取学生姓名失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	g, err := service.HashGet(tokenClaims.Identify, "gender")
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "token错误")
			return
		}
		fmt.Println("获取学生性别失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	grade, err := service.HashGet(tokenClaims.Identify, "grade")
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "token错误")
			return
		}
		fmt.Println("获取学生年级失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	class, err := service.HashGet(tokenClaims.Identify, "class")
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "token错误")
			return
		}
		fmt.Println("获取学生班级失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	department, err := service.HashGet(tokenClaims.Identify, "department")
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "token错误")
			return
		}
		fmt.Println("获取学生院系失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	major, err := service.HashGet(tokenClaims.Identify, "major")
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "token错误")
			return
		}
		fmt.Println("获取学生专业失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	student := model.Student{
		StudentId:   tokenClaims.Identify,
		StudentName: studentName,
		Gender:      g,
		Grade:       grade,
		Class:       class,
		Department:  department,
		Major:       major,
		RuleId:      "student",
	}
	tool.Success(ctx, 200, student)
}

//获取头像
func getAvatar(ctx *gin.Context) {
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		fmt.Println("token解析失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	bucket, err := service.ParseBucket()
	if err != nil {
		fmt.Println("解析存储桶错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	u, _ := url.Parse(bucket.Url)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  bucket.SecretId,
			SecretKey: bucket.SecretKey,
		},
	})
	avatar := client.Object.GetObjectURL(tokenClaims.Identify)

	tool.Success(ctx, 200, avatar.Scheme+"://"+avatar.Host+avatar.Path)
}

//选课
func chooseCourse(ctx *gin.Context) {
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		fmt.Println("token解析失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	courseNumber := ctx.PostForm("courseNumber")
	teachingClass := ctx.PostForm("teachingClass")
	if courseNumber == "" || teachingClass == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}

	//查询课程
	_, err = service.HashGet("course", courseNumber)
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "课程不存在")
			return
		}
		fmt.Println("查询课程编号是否存在失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	//查询教学班
	_, err = service.HashGet(courseNumber, teachingClass)
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "教学班不存在")
			return
		}
		fmt.Println("查询教学班是否存在失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	//查询学生是否已经选择过该课程
	mobile, err := service.HashGet(tokenClaims.Identify, "mobile")
	if err != nil {
		fmt.Println("查询电话号码失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	_, err = service.HashGet(mobile, courseNumber)
	if err != nil && err != redis.Nil {
		fmt.Println("查询学生是否已经选则该课程失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	if err == redis.Nil {
		//判断选课时间是否冲突
		err, selectedTeachingClassArr := service.HVals(mobile)
		if err != nil {
			fmt.Println("获取学生已加入教学班失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}

		if len(selectedTeachingClassArr) > 1 {
			var selectedTimeArr []string
			for _, v := range selectedTeachingClassArr {
				if v != tokenClaims.Identify {
					selectedTime, err := service.HashGet(v, "setTime")
					if err != nil {
						fmt.Println("查询课程开设时间失败", err)
						tool.Failure(ctx, 500, "服务器错误")
						return
					}
					timeArr := strings.Split(selectedTime, ",")
					for _, val := range timeArr {
						selectedTimeArr = append(selectedTimeArr, val)
					}
				}
				//查询当前所选课程时间
				setTime, err := service.HashGet(teachingClass, "setTime")
				if err != nil {
					fmt.Println("查询当前所选课程开设时间出错", err)
					tool.Failure(ctx, 500, "服务器错误")
					return
				}
				setTimeArr := strings.Split(setTime, ",")

				//判断课程是否存在时间冲突
				ok := service.JudgeTimeConflict(selectedTimeArr, setTimeArr)
				if ok {
					tool.Failure(ctx, 400, "课程存在时间冲突")
					return
				}
			}
		}

		//将选课信息存入学生信息
		err = service.HashSet(mobile, courseNumber, teachingClass)
		if err != nil {
			fmt.Println("将选课信息存入学生信息失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}

		//将选课信息存入课程信息
		err = service.HashSet(teachingClass, tokenClaims.Identify, mobile)
		if err != nil {
			fmt.Println("将选课信息存入课程信息失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		tool.Success(ctx, 200, "选课成功")
		return
	}
	tool.Failure(ctx, 400, "你已选择过该课程")
}

//查询个人选课信息
func selection(ctx *gin.Context) {
	//解析token获取id
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		fmt.Println("token解析失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	mobile, err := service.HashGet(tokenClaims.Identify, "mobile")
	if err != nil {
		fmt.Println("查询电话号码失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//查询所选课程编号
	err, courseNumberArr := service.HKeys(mobile)
	if err != nil {
		fmt.Println("查询课程编号失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	teaching := model.Selection{}
	var teachingSum []model.Selection
	for _, v := range courseNumberArr {
		if v != "studentId" {
			teaching.CourseNumber = v
			teaching.CourseCredit, err = service.HashGet(v, "courseCredit")
			if err != nil {
				fmt.Println("查询课程学分失败", err)
				tool.Failure(ctx, 500, "服务器错误")
				return
			}
			teaching.CourseType, err = service.HashGet(v, "courseType")
			if err != nil {
				fmt.Println("查询课程类型失败", err)
				tool.Failure(ctx, 500, "服务器错误")
				return
			}
			//获取教学班编号
			teachingClassNumber, err := service.HashGet(mobile, v)
			if err != nil {
				fmt.Println("查询教学班编号失败", err)
				tool.Failure(ctx, 500, "服务器错误")
				return
			}
			//获取教学班开设时间
			teaching.SetTime, err = service.HashGet(teachingClassNumber, "setTime")
			if err != nil {
				fmt.Println("查询教学班开设时间失败", err)
				tool.Failure(ctx, 500, "服务器错误")
				return
			}
			//查询教师工号
			workNumber, err := service.HashGet(teachingClassNumber, "workNumber")
			if err != nil {
				fmt.Println("查询教师编号失败", err)
				tool.Failure(ctx, 500, "服务器错误")
				return
			}
			//获取教师姓名
			teaching.TeacherName, err = service.HashGet(workNumber, "teacherName")
			if err != nil {
				fmt.Println("查询教师姓名失败", err)
				tool.Failure(ctx, 500, "服务器错误")
				return
			}
			teachingSum = append(teachingSum, teaching)
		}
	}

	tool.Success(ctx, http.StatusOK, teachingSum)
}

//退课
func quit(ctx *gin.Context) {
	//解析token获取id
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		fmt.Println("token解析失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	courseNumber := ctx.PostForm("courseNumber")
	classNumber := ctx.PostForm("classNumber")
	if classNumber == "" || courseNumber == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	//查询电话号码
	mobile, err := service.HashGet(tokenClaims.Identify, "mobile")
	if err != nil {
		fmt.Println("查询电话号码失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//删除学生信息中的选课信息
	err = service.HDelSingle(mobile, courseNumber)
	if err != nil {
		fmt.Println("删除学生信息中的选课信息失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//删除教学班信息中的学生信息
	err = service.HDelSingle(classNumber, tokenClaims.Identify)
	if err != nil {
		fmt.Println("删除教学班信息中的学生信息失败")
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	tool.Success(ctx, 200, "你已经退出该班级")
}
