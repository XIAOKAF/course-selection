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

func Del(key string) error {
	err := dao.Del(key)
	if err != nil {
		return err
	}
	return nil
}

func HashSet(userId string, fieldName string, value string) error {
	fields := make(map[string]interface{})
	fields[fieldName] = value
	err := dao.HashSet(userId, fields)
	if err != nil {
		return err
	}
	return nil
}

func HashGet(userId string, fieldName string) (string, error) {
	value, err := dao.HashGet(userId, fieldName)
	if err != nil {
		return value, err
	}
	return value, nil
}
