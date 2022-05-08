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
	"os"
	"strconv"
	"strings"
)

func studentRegister(ctx *gin.Context) {
	studentId := ctx.PostForm("userId")
	password := ctx.PostForm("password")
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

func studentLogin(ctx *gin.Context) {
	unifiedCode := ctx.PostForm("unifiedCode")
	password := ctx.PostForm("pwd")
	auth := ctx.PostForm("auth")
	if unifiedCode == "" || password == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	flag, err, pwd := service.SelectUnifiedCode(unifiedCode)
	if err != nil {
		fmt.Println("查询统一验证码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if !flag {
		tool.Failure(ctx, 400, "是本校学生(⊙o⊙)吗？")
		return
	}
	if pwd != password {
		tool.Failure(ctx, 400, "密码错误（提示一下哦，初始密码是姓名拼音")
		return
	}
	if auth == "" {
		err, token := service.CreateToken(unifiedCode, 2)
		if err != nil {
			fmt.Println("创建token失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		err = service.HashSet("token", unifiedCode, token)
		if err != nil {
			fmt.Println("储存token失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		tool.Success(ctx, 200, token)
		return
	}
	err, token := service.CreateToken(unifiedCode, 2)
	if err != nil {
		fmt.Println("创建token失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err, refreshToken := service.RememberStatus(unifiedCode, 5)
	if err != nil {
		fmt.Println("创建refreshToken失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet("token", unifiedCode, token)
	if err != nil {
		fmt.Println("储存token失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.HashSet("refreshToken", unifiedCode, refreshToken)
	if err != nil {
		fmt.Println("储存refreshToken失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	tool.Success(ctx, 200, token)
}

//通过原密码修改密码
func changePwdByOldPwd(ctx *gin.Context) {
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		fmt.Println("token解析失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	oldPwd := ctx.PostForm("oldPwd")
	newPwd := ctx.PostForm("newPwd")
	if oldPwd == "" {
		tool.Failure(ctx, 400, "悄悄提醒你，初始验证码为姓名拼音哦")
		return
	}
	if newPwd == "" {
		tool.Failure(ctx, 400, "(・∀・新密码是？")
		return
	}
	result, err := service.HashGet(tokenClaims.UserId, "password")
	if err != nil {
		fmt.Println("查询旧密码失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if result != oldPwd {
		tool.Failure(ctx, 400, "原来的密码不正确哦")
		return
	}
	student := model.Student{
		UnifiedCode: tokenClaims.UserId,
		Password:    newPwd,
	}
	//MySQL更新
	err = service.UpdatePassword(student)
	if err != nil {
		fmt.Println("MySQL更新密码失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//redis更新
	err = service.HashSet(tokenClaims.UserId, "password", newPwd)
	if err != nil {
		fmt.Println("redis更新密码失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	tool.Success(ctx, 200, "成功♪(^∇^*)")
}

func updateMobile(ctx *gin.Context) {
	//新电话号码
	newMobile := ctx.PostForm("newMobile")
	if newMobile == "" {
		tool.Failure(ctx, 400, "电话号码不能为空哦")
		return
	}
	//确认登录状态
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		fmt.Println("token解析失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	_, err = service.HashGet(tokenClaims.UserId, "studentName")
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "统一验证码不存在")
			return
		}
		if err != nil {
			fmt.Println("查询统一验证码失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
	}
	//发送校验短信
	code := service.CreateCode()
	var sms model.Message
	sms, err = service.ParseSmsConfig(sms)
	if err != nil {
		fmt.Println("解析短信配置文件错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err = service.SendSms(newMobile, code, sms)
	if err != nil {
		fmt.Println("短信发送错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//redis存储
	err = service.Set(tokenClaims.UserId, code, 2)
	if err != nil {
		fmt.Println("redis储存验证码失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	tool.Success(ctx, 200, "成功发送短信")
}

//更新电话号码时校验验证码
func checkCodeForUpdate(ctx *gin.Context) {
	//确认登录状态
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		fmt.Println("解析token失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	newMobile := ctx.PostForm("newMobile")
	code := ctx.PostForm("code")
	//验证码是否正确且在保质期内
	result, duration, err := service.CheckSms(newMobile)
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "验证码已过期或电话号码错误")
			return
		}
	}
	if err != nil {
		fmt.Println("查询验证码失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if duration == -1 {
		fmt.Println(ctx, "验证码没有设置过期时间")
		//删除电话号码-验证码键值对
		err = service.Del(newMobile)
		if err != nil {
			fmt.Println("删除电话号码-验证码键值对失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if code != result {
		tool.Failure(ctx, 400, "验证码错误")
		return
	}
	//跟新MySQL
	student := model.Student{
		UnifiedCode: tokenClaims.UserId,
		Mobile:      newMobile,
	}
	err = service.UpdateMobile(student)
	if err != nil {
		fmt.Println("MySQL更新失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//更新redis
	err = service.HashSet(tokenClaims.UserId, "mobile", newMobile)
	if err != nil {
		fmt.Println("redis更新电话号码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//删除新电话号码-验证码键值对
	err = service.Del(newMobile)
	if err != nil {
		fmt.Println("redis删除验证码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	tool.Success(ctx, 200, "电话号码更新成功")
}

//更新头像
func updateAvatar(ctx *gin.Context) {
	//确认登录状态
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		fmt.Println("token解析失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//获取头像文件
	avatar, err := os.Open("avatar")
	if err != nil {
		fmt.Println("获取头像错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	u, err := url.Parse("https://examplebucket-1250000000.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BatchURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  " ",
			SecretKey: " ",
		},
	})
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "ipj",
		},
	}
	_, err = client.Object.Delete(ctx, tokenClaims.UserId)
	if err != nil {
		fmt.Println("删除原有头像错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	_, err = client.Object.Put(ctx, tokenClaims.UserId, avatar, opt)
	if err != nil {
		fmt.Println("储存头像错误", err)
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
	studentName, err := service.HashGet(tokenClaims.UserId, "studentName")
	if err != nil {
		fmt.Println("获取学生姓名失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	gender, err := service.HashGet(tokenClaims.UserId, "gender")
	if err != nil {
		fmt.Println("获取学生性别失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	g, err := strconv.Atoi(gender)
	if err != nil {
		fmt.Println("string转int失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	grade, err := service.HashGet(tokenClaims.UserId, "grade")
	if err != nil {
		fmt.Println("获取学生年级失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	class, err := service.HashGet(tokenClaims.UserId, "class")
	if err != nil {
		fmt.Println("获取学生班级失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	department, err := service.HashGet(tokenClaims.UserId, "department")
	if err != nil {
		fmt.Println("获取学生院系失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	major, err := service.HashGet(tokenClaims.UserId, "major")
	if err != nil {
		fmt.Println("获取学生专业失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	//从腾讯云对象储存内获取头像
	u, _ := url.Parse("https://examplebucket-1250000000.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  " ",
			SecretKey: " ",
		},
	})
	file := "localfile"
	opt := &cos.MultiDownloadOptions{
		ThreadPoolSize: 5,
	}
	_, err = client.Object.Download(ctx, tokenClaims.UserId, file, opt)
	if err != nil {
		fmt.Println("获取图片错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	student := model.Student{
		StudentName: studentName,
		Gender:      g,
		Grade:       grade,
		Class:       class,
		Department:  department,
		Major:       major,
	}
	tool.Success(ctx, 200, student)
}

//选课
func chooseCourse(ctx *gin.Context) {
	//中间件验证请求头是否携带token且token存在并合格
	//从token中获取id方便后续插入数据的操作
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
	//在redis中查询该课程编号哈希表
	_, err = service.HashGet(courseNumber, "courseName")
	if err != nil {
		fmt.Println("查询课程编号失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//查询学生是否已经选择过该课程
	err, flag := service.HExists(tokenClaims.UserId, courseNumber)
	if err != nil {
		fmt.Println("查询学生是否已经选择过该课程失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if flag {
		tool.Failure(ctx, 400, "你已经选择过该课程")
		return
	}
	//判断选课时间是否冲突
	//查询学生已选课程时间
	err, selectCurriculumArr := service.HKeys(tokenClaims.UserId + "teaching")
	if err != nil {
		fmt.Println("查询学生已选择课程失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err, selectCourseArr := service.HVals(tokenClaims.UserId + "teaching")
	if err != nil {
		fmt.Println("获取学生已加入教学班失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	selectString := ""
	var build strings.Builder
	for i, _ := range selectCourseArr {
		selectTime, err := service.HashGet(selectCurriculumArr[i]+"teaching", selectCourseArr[i])
		if err != nil {
			fmt.Println("查询课程开设时间失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		build.WriteString(selectString)
		if i == 0 {
			build.WriteString(selectTime)
		} else {
			build.WriteString(" " + selectTime)
		}
		selectString = build.String()
	}
	selectArr := strings.Fields(selectString)
	//查询当前所选课程时间
	setTime, err := service.HashGet(courseNumber+"teaching", teachingClass)
	if err != nil {
		fmt.Println("查询当前所选课程开设时间出错", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	setTimeArr := strings.Fields(setTime)
	//判断二者是否有交集
	flag = service.IsRepeated(selectArr, setTimeArr)
	if !flag {
		tool.Failure(ctx, 400, "课程时间出现冲突")
		return
	}
	//根据token中提供的学生统一验证码检索学生姓名
	name, err := service.HashGet(tokenClaims.UserId, "studentName")
	choice := model.Choice{
		TeachingClass: teachingClass,
		UnifiedCode:   tokenClaims.UserId,
		StudentName:   name,
	}
	//将学生选课信息存入redis
	//存入统一验证码为名的哈希表
	err = service.HashSet(tokenClaims.UserId, courseNumber, teachingClass)
	if err != nil {
		fmt.Println("将选课信息存入redis失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//存入以教学班编号为名的哈希表
	err = service.HashSet(choice.TeachingClass, choice.UnifiedCode, choice.StudentName)
	if err != nil {
		fmt.Println("将选课信息存入redis失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	tool.Success(ctx, 200, "选课成功")
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
	err, courseNumberArr := service.HKeys(tokenClaims.UserId + "teaching")
	if err != nil {
		fmt.Println("查询课程编号失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	err, classNumberArr := service.HVals(tokenClaims.UserId + "teaching")
	if err != nil {
		fmt.Println("查询教学班失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	teaching := model.Selection{}
	var teachingSum []model.Selection
	for i, v := range courseNumberArr {
		teaching.CourseNumber = v
		teaching.TeachingClass = classNumberArr[i]
		teaching.CourseCredit, err = service.HashGet(tokenClaims.UserId, "courseCredit")
		if err != nil {
			fmt.Println("查询课程学分失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		teaching.CourseType, err = service.HashGet(tokenClaims.UserId, "courseType")
		if err != nil {
			fmt.Println("查询课程类型失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		teaching.SetTime, err = service.HashGet(v+"teaching", classNumberArr[i])
		if err != nil {
			fmt.Println("查询教学班开设时间失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		teacherNumber, err := service.HashGet(v+"teacher", classNumberArr[i])
		if err != nil {
			fmt.Println("查询教师编号失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		teaching.TeacherName, err = service.HashGet(teacherNumber, "teacherName")
		if err != nil {
			fmt.Println("查询教师姓名失败", err)
			tool.Failure(ctx, 500, "服务器错误")
			return
		}
		teachingSum[i] = teaching
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
	if classNumber == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	var idArr []string
	var courseNumberArr []string
	courseNumberArr = append(courseNumberArr, courseNumber)
	err = service.HDel(tokenClaims.UserId, courseNumberArr)
	if err != nil {
		fmt.Println("删除学生信息中的选课失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	idArr = append(idArr, tokenClaims.UserId)
	err = service.HDel(classNumber, idArr)
	if err != nil {
		fmt.Println("删除教学班内的学生信息失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	tool.Success(ctx, 200, "你已经退出该班级")
}
