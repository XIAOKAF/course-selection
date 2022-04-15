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
