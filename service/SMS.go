package service

import (
	"course-selection/dao"
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"time"
)

func CreateCode() string {
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	return code
}

func SelectMobile(mobile string) (bool, error) {
	_, err := dao.Get(mobile)
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

func CheckSms(mobile string) (string, time.Duration, bool, error) {
	result, err := dao.Get(mobile)
	if err != nil {
		if err == redis.Nil {
			return result, -2, false, nil
		}
		return result, -2, true, err
	}
	duration, err := dao.TTL(mobile)
	if err != nil {
		return result, duration, true, err
	}
	return result, duration, true, nil
}
