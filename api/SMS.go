package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

//发送短信
func sendSms(ctx *gin.Context) {
	mobile := ctx.PostForm("mobile")

	//查询电话号码是否存在
	_, err := service.HashGet(mobile, "studentId")

	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "电话号码错误")
			return
		}
		fmt.Println("查询电话号码是否存在失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	//生成随机验证码
	code := service.CreateCode()

	//解析短信配置文件
	var sms model.Message
	sms, err = service.ParseSmsConfig()
	if err != nil {
		fmt.Println("解析短信配置文件错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	//连接
	credential := common.NewCredential(sms.SecretId, sms.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, err := tencentsms.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		fmt.Println("短信相关配置错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	request := tencentsms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(sms.AppId)
	request.SignName = common.StringPtr(sms.Sign)
	request.SenderId = common.StringPtr("")
	request.ExtendCode = common.StringPtr("")
	request.TemplateParamSet = common.StringPtrs([]string{code, "5"})
	request.TemplateId = common.StringPtr(sms.TemplateId)
	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + mobile})

	//发送短信
	_, err = client.SendSms(request)
	if err != nil {
		fmt.Println(err)
		fmt.Println("发送短信错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	//将验证码存入redis之中并且设置过期时间(mobile+"code"-verifiedCode)
	err = service.Set(mobile+"code", code, 5)
	if err != nil {
		fmt.Println("储存验证码错误", err)
		tool.Failure(ctx, 200, "服务器错误")
		return
	}
	tool.Success(ctx, 200, "短信发送成功(p≧w≦q)")
}

//校验验证码
func checkSms(ctx *gin.Context) {
	mobile := ctx.PostForm("mobile")
	code := ctx.PostForm("code")
	result, duration, err := service.CheckSms(mobile + "code")
	if err != nil {
		if err == redis.Nil {
			tool.Failure(ctx, 400, "验证码已过期")
			return
		}
		fmt.Println("查询验证码失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}

	if duration == -1 {
		fmt.Println("验证码没有设置过期时间")
		//删除电话号码-验证码键值对
		err = service.Del(mobile + "code")
		if err != nil {
			fmt.Println("删除验证码错误且验证码未设置过期时间", err)
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
	tool.Failure(ctx, 200, "验证码正确")
}
