package service

import (
	"course-selection/dao"
	"fmt"
	"math/rand"
	"time"
)

func CreateCode() string {
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	return code
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
