package service

import (
	"course-selection/dao"
	"github.com/go-redis/redis"
	"time"
)

func Set(key string, val string, expiration time.Duration) error {
	err := dao.Set(key, val, expiration)
	if err != nil {
		return err
	}
	return nil
}

func Get(key string) (string, bool, error) {
	result, err := dao.Get(key)
	if err != nil {
		if err == redis.Nil {
			//false表示未查询到该键，反之查询到
			return "", false, nil
		}
		return "", true, err
	}
	return result, true, err
}
