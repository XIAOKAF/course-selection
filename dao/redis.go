package dao

import (
	"time"
)

func Set(key string, val string, expiration time.Duration) error {
	err := RDB.Set(key, val, expiration*time.Minute)
	if err != nil {
		return err.Err()
	}
	return nil
}

func Get(key string) (string, error) {
	result, err := RDB.Get(key).Result()
	if err != nil {
		return result, err
	}
	return result, nil
}

func Del(key string) error {
	result := RDB.Del(key)
	return result.Err()
}

func TTL(key string) (time.Duration, error) {
	result := RDB.TTL(key)
	return result.Result()
}

// HashSet 将用户信息以键值对的信息存入哈希表
func HashSet(hashTableName string, fields map[string]interface{}) error {
	result := RDB.HMSet(hashTableName, fields)
	return result.Err()
}

// HashGet 获取哈希表中指定的值
func HashGet(hashTableName, fieldsName string) (string, error) {
	result := RDB.HGet(hashTableName, fieldsName)
	return result.Val(), result.Err()
}
