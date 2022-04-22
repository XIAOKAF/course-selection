package service

import (
	"course-selection/dao"
	"course-selection/model"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"math/rand"
	"time"
)

func CreateCode() string {
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	return code
}

func SendSms(mobile string, code string, sms model.Message) error {
	//连接
	credential := common.NewCredential(sms.SignId, sms.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, err := tencentsms.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		return err
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
		return err
	}
	//将验证码存入redis之中并且设置过期时间(mobile:verifiedCode)
	err = dao.Set(mobile, code, 2)
	if err != nil {
		return err
	}
	return nil
}

// CheckSms 查询验证码以及过期时间
func CheckSms(mobile string) (string, time.Duration, error) {
	code, err := dao.Get(mobile)
	if err != nil {
		return code, 0, err
	}
	duration, err := dao.TTL(mobile)
	if err != nil {
		return code, duration, err
	}
	return code, duration, nil
}

// IsMobileExist 查询电话号码是否存在
func IsMobileExist(mobile string) (bool, error) {
	_, err := dao.Get(mobile)
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return true, err
	}
	return true, nil
}
