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
	return value, err
}

func HashGetAll(hashTableName string) (error, map[string]string) {
	err, result := dao.HashGetAll(hashTableName)
	return err, result
}

func HExists(hashTableName string, filedName string) (error, bool) {
	err, flag := dao.HExists(hashTableName, filedName)
	return err, flag
}

func HKeys(hashTableName string) (error, []string) {
	err, keys := dao.HKeys(hashTableName)
	return err, keys
}

func HVals(hashTableName string) (error, []string) {
	err, vals := dao.HVals(hashTableName)
	return err, vals
}

func HDel(hashTableName string, filedNameArr []string) error {
	err := dao.HDel(hashTableName, filedNameArr)
	return err
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
