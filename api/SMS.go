package api

import (
	"course-selection/dao"
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

//发送短信（短信登录，通过短信找回密码
func sendSms(ctx *gin.Context) {
	mobile := ctx.PostForm("mobile")
	//查询电话号码是否存在
	flag, err := service.SelectMobile(mobile)
	if err != nil {
		fmt.Println("查询电话号码错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	if flag {
		tool.Failure(ctx, 400, "电话号码不存在")
		return
	}
	//生成随机验证码
	code := service.CreateCode()
	//解析短信配置文件
	var sms model.Message
	sms, err = service.ParseSmsConfig(sms)
	if err != nil {
		fmt.Println("解析短信配置文件错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//连接
	credential := common.NewCredential(sms.SignId, sms.SecretKey)
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
	request.TemplateParamSet = common.StringPtrs([]string{code})
	request.TemplateId = common.StringPtr(sms.TemplateId)
	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + mobile})
	//发送短信
	_, err = client.SendSms(request)
	if err != nil {
		fmt.Println("发送短信错误", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//将验证码存入redis之中并且设置过期时间(mobile:verifiedCode)
	dao.Set(mobile, code, 2)
	tool.Success(ctx, 200, "短信发送成功(p≧w≦q)")
}
