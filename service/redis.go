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

func HashSet(hashTableName, fieldName string, value string) error {
	fields := make(map[string]interface{})
	fields[fieldName] = value
	err := dao.HashSet(hashTableName, fields)
	if err != nil {
		return err
	}
	return nil
}

func HashGet(hashTableName string, fieldName string) (string, error) {
	value, err := dao.HashGet(hashTableName, fieldName)
	if err != nil {
		return value, err
	}
	return value, nil
}

func HashGetAll(hashTableName string) (error, map[string]string) {
	err, result := dao.HashGetAll(hashTableName)
	return err, result
}

func SetAdd(key string, member interface{}) error {
	err := dao.SetAdd(key, member)
	return err
}

func SetGet(setName string) (error, []string) {
	err, members := dao.SetGet(setName)
	return err, members
}

func SScan(key string, cursor uint64, match string, count int64) []string {
	val, _ := dao.SScan(key, cursor, match, count)
	return val
}

func SIsMember(key string, member interface{}) (error, bool) {
	err, flag := dao.SIsMember(key, member)
	return err, flag
}
