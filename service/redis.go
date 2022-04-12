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

func Get(key string) (error, bool, string) {
	err, result := dao.Get(key)
	if err != nil {
		if err == redis.Nil {
			//false表示未查询到该键，反之查询到
			return nil, false, ""
		}
		return err, true, ""
	}
	return nil, true, result
}
