package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	}
	tool.Success(ctx, 200, "亲爱的"+unifiedCode+"，你已经成功激活账户啦！o(*￣▽￣*)ブ")
}

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
		tool.Failure(ctx, 400, "密码不能为空哦，悄悄提醒你，初始验证码为姓名拼音哦")
		return
	}
	if newPwd == "" {
		tool.Failure(ctx, 400, "你还要不要改密码了")
		return
	}
	result, flag, err := service.Get(tokenClaims.UserId)
	if err != nil {
		fmt.Println("查询统一验证码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if !flag {
		tool.Failure(ctx, 400, "该统一验证码不存在")
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
	err = service.UpdatePassword(student)
	if err != nil {
		fmt.Println("跟新密码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	tool.Success(ctx, 200, "成功♪(^∇^*)")
}

func updateMobile(ctx *gin.Context) {
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
	//新电话号码
	newMobile := ctx.PostForm("newMobile")
	if newMobile == "" {
		tool.Failure(ctx, 400, "电话号码不能为空哦")
		return
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
	//将新电话号码存入redis之中
	err = service.Set(newMobile, code, 2)
	if err != nil {
		fmt.Println("储存新电话号码错误", err)
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
	newMobile := ctx.PostForm("newMobile")
	code := ctx.PostForm("code")
	//查询新电话号码是否正确
	flag, err = service.IsMobileExist(newMobile)
	if err != nil {
		fmt.Println("查询新电话号码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if !flag {
		tool.Failure(ctx, 400, "电话号码错误")
		return
	}
	//验证码是否正确且在保质期内
	result, time, err := service.CheckSms(newMobile)
	if err != nil {
		fmt.Println("查询验证码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if time < 0 {
		tool.Failure(ctx, 400, "验证码已过期")
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
		fmt.Println("MySQL更新电话号码错误", err)
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
		fmt.Println("删除新电话号码错误", err)
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

func SelectPersonalCourse(ctx *gin.Context) {
	//获取token以便后续查询
	tokenString := ctx.Request.Header.Get("token")
	tokenClaims, err := service.ParseToken(tokenString)
	tool.DealWithErr(ctx, err, "解析token失败")
	_, flag, err := service.Get(tokenClaims.UserId)
	tool.DealWithErr(ctx, err, "从redis中查询统一验证码失败")
	if !flag {
		tool.Failure(ctx, 400, "token不存在")
		return
	}

}
