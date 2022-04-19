package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
)

func studentRegister(ctx *gin.Context) {
	unifiedCode := ctx.PostForm("unifiedCode")
	if unifiedCode == "" {
		tool.Failure(ctx, 400, "还没有输入统一验证码哦")
		return
	}
	//查询该学生是否是本校学生，是则返回true，不是则返回false
	flag, err := service.SelectStudentByUnifiedCode(unifiedCode)
	if err != nil {
		fmt.Println("查询统一验证码错误", err)
		tool.Failure(ctx, 400, "服务器错误")
		return
	}
	if !flag {
		tool.Failure(ctx, 400, "是本校学生(⊙o⊙)吗？")
		return
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
