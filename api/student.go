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
	unifiedCode := ctx.PostForm("unifiedCode")
	password := ctx.PostForm("password")
	if unifiedCode == "" || password == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	//查询该学生是否是本校学生，是则返回true，不是则返回false
	flag, err, pwd := service.SelectUnifiedCode(unifiedCode)
	tool.DealWithErr(ctx, err, "查询统一验证码错误")
	if !flag {
		tool.Failure(ctx, 400, "是本校学生(⊙o⊙)吗？")
		return
	}
	if pwd != password {
		tool.Failure(ctx, 400, "密码错误（提示一下哦，初始密码是姓名拼音")
		return
	}
	tool.Success(ctx, 200, "亲爱的"+unifiedCode+"，你已经成功激活账户啦！o(*￣▽￣*)ブ")
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
	tool.DealWithErr(ctx, err, "查询统一验证码错误")
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
		tool.DealWithErr(ctx, err, "创建token错误")
		err = service.HashSet("token", unifiedCode, token)
		tool.DealWithErr(ctx, err, "存储token错误")
		tool.Success(ctx, 200, token)
		return
	}
	err, token := service.CreateToken(unifiedCode, 2)
	tool.DealWithErr(ctx, err, "创建token错误")
	err, refreshToken := service.RememberStatus(unifiedCode, 5)
	tool.DealWithErr(ctx, err, "创建refreshToken错误")
	err = service.HashSet("token", unifiedCode, token)
	tool.DealWithErr(ctx, err, "储存token失败")
	err = service.HashSet("refreshToken", unifiedCode, refreshToken)
	tool.DealWithErr(ctx, err, "存储refreshToken错误")
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
		tool.DealWithErr(ctx, err, "查询旧密码错误")
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
	tool.DealWithErr(ctx, err, "MySQL更新密码错误")
	//redis更新
	err = service.HashSet(tokenClaims.UserId, "password", newPwd)
	tool.DealWithErr(ctx, err, "redis更新密码错误")
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
		tool.DealWithErr(ctx, err, "查询统一验证码错误")
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
	tool.DealWithErr(ctx, err, "redis储存验证码错误")
	tool.Success(ctx, 200, "成功发送短信")
}

//更新电话号码时校验验证码
func checkCodeForUpdate(ctx *gin.Context) {
	//确认登录状态
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	tool.DealWithErr(ctx, err, "解析token出错")
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
	tool.DealWithErr(ctx, err, "查询验证码错误")
	if duration == -1 {
		fmt.Println(ctx, "验证码没有设置过期时间")
		//删除电话号码-验证码键值对
		err = service.Del(newMobile)
		tool.DealWithErr(ctx, err, "删除电话号码验证码键值对出错")
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
	tool.DealWithErr(ctx, err, "MySQL更新出错")
	//更新redis
	err = service.HashSet(tokenClaims.UserId, "mobile", newMobile)
	if err != nil {
		fmt.Println("redis更新电话号码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//删除新电话号码-验证码键值对
	err = service.Del(newMobile)
	tool.DealWithErr(ctx, err, "redis删除验证码出错")
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
	tool.DealWithErr(ctx, err, "获取学生姓名错误")
	gender, err := service.HashGet(tokenClaims.UserId, "gender")
	tool.DealWithErr(ctx, err, "获取学生性别错误")
	g, err := strconv.Atoi(gender)
	tool.DealWithErr(ctx, err, "string转int错误")
	grade, err := service.HashGet(tokenClaims.UserId, "grade")
	tool.DealWithErr(ctx, err, "获取学生年级错误")
	class, err := service.HashGet(tokenClaims.UserId, "class")
	tool.DealWithErr(ctx, err, "获取学生班级错误")
	department, err := service.HashGet(tokenClaims.UserId, "department")
	tool.DealWithErr(ctx, err, "获取学生院系错误")
	major, err := service.HashGet(tokenClaims.UserId, "major")
	tool.DealWithErr(ctx, err, "获取学生专业错误")

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
	tool.DealWithErr(ctx, err, "从腾讯云下载图片出错")

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
	tool.DealWithErr(ctx, err, "token解析出错")
	courseNumber := ctx.PostForm("courseNumber")
	teachingClass := ctx.PostForm("teachingClass")
	if courseNumber == "" || teachingClass == "" {
		tool.Failure(ctx, 400, "必要字段不能为空")
		return
	}
	//在redis中查询该课程编号哈希表
	_, err = service.HashGet(courseNumber, "courseName")
	tool.DealWithErr(ctx, err, "查询课程编号出错")
	//查询学生是否已经选择过该课程
	err, flag := service.HExists(tokenClaims.UserId, courseNumber)
	tool.DealWithErr(ctx, err, "查询学生是否已经选择过该课程出错")
	if flag {
		tool.Failure(ctx, 400, "你已经选择过该课程")
		return
	}
	//判断选课时间是否冲突
	//查询学生已选课程时间
	err, selectCurriculumArr := service.HKeys(tokenClaims.UserId + "teaching")
	tool.DealWithErr(ctx, err, "查询学生已选择课程出错")
	err, selectCourseArr := service.HVals(tokenClaims.UserId + "teaching")
	tool.DealWithErr(ctx, err, "查询学生已加入教学班出错")
	selectString := ""
	var build strings.Builder
	for i, _ := range selectCourseArr {
		selectTime, err := service.HashGet(selectCurriculumArr[i]+"teaching", selectCourseArr[i])
		tool.DealWithErr(ctx, err, "查询课程开设时间出错")
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
	tool.DealWithErr(ctx, err, "查询当前所选课程开设时间出错")
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
	tool.DealWithErr(ctx, err, "将选课信息存入redis出错")
	//存入以教学班编号为名的哈希表
	err = service.HashSet(choice.TeachingClass, choice.UnifiedCode, choice.StudentName)
	tool.DealWithErr(ctx, err, "将选课信息存入redis出错")
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
	tool.DealWithErr(ctx, err, "查询课程编号失败")
	err, classNumberArr := service.HVals(tokenClaims.UserId + "teaching")
	tool.DealWithErr(ctx, err, "查询教学班失败")
	teaching := model.Selection{}
	var teachingSum []model.Selection
	for i, v := range courseNumberArr {
		teaching.CourseNumber = v
		teaching.TeachingClass = classNumberArr[i]
		teaching.CourseCredit, err = service.HashGet(tokenClaims.UserId, "courseCredit")
		tool.DealWithErr(ctx, err, "查询课程学分错误")
		teaching.CourseType, err = service.HashGet(tokenClaims.UserId, "courseType")
		tool.DealWithErr(ctx, err, "查询课程类型错误")
		teaching.SetTime, err = service.HashGet(v+"teaching", classNumberArr[i])
		tool.DealWithErr(ctx, err, "查询教学班开设时间错误")
		teacherNumber, err := service.HashGet(v+"teacher", classNumberArr[i])
		tool.DealWithErr(ctx, err, "查询教师编号错误")
		teaching.TeacherName, err = service.HashGet(teacherNumber, "teacherName")
		tool.DealWithErr(ctx, err, "查询教师姓名错误")
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
	tool.DealWithErr(ctx, err, "删除学生信息中的选课记录失败")
	idArr = append(idArr, tokenClaims.UserId)
	err = service.HDel(classNumber, idArr)
	tool.DealWithErr(ctx, err, "删除教学班内的学生信息失败")
	tool.Success(ctx, 200, "你已经退出该班级")
}
