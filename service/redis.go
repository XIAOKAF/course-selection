package service

import (
	"course-selection/dao"
	"time"
)

func Set(key string, val string, expiration time.Duration) error {
	err := dao.Set(key, val, expiration*time.Minute)
	if err != nil {
		return err
	}
	return nil
}

func Get(key string) (string, error) {
	result, err := dao.Get(key)
	return result, err
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

func HDelSingle(key string, filed string) error {
	return dao.HDelSingle(key, filed)
}

func SScan(key string, cursor uint64, match string, count int64) []string {
	val, _ := dao.SScan(key, cursor, match, count)
	return val
}
