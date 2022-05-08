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

// HashGetAll 获取哈希表中所有的值
func HashGetAll(hashTableName string) (error, map[string]string) {
	result := RDB.HGetAll(hashTableName)
	return result.Err(), result.Val()
}

// HExists 检索哈希表中是否存在指定字段
func HExists(key string, filed string) (error, bool) {
	result := RDB.HExists(key, filed)
	return result.Err(), result.Val()
}

// HKeys 获取哈希表中所有字段
func HKeys(key string) (error, []string) {
	result := RDB.HKeys(key)
	return result.Err(), result.Val()
}

// HVals 获取哈希表中所有值
func HVals(key string) (error, []string) {
	result := RDB.HVals(key)
	return result.Err(), result.Val()
}

// HDel 删除哈希表中的值
func HDel(key string, fields []string) error {
	for _, i := range fields {
		result := RDB.HDel(key, i)
		if result.Err() != nil {
			return result.Err()
		}
	}
	return nil
}

// HDelSingle 删除哈希表中的一个值
func HDelSingle(key string, filed string) error {
	result := RDB.HDel(key, filed)
	return result.Err()
}

// SetAdd 像列表中插入数据
func SetAdd(key string, member interface{}) error {
	result := RDB.SAdd(key, member)
	return result.Err()
}

// SetGet 从列表中获取数据
func SetGet(setName string) (error, []string) {
	result := RDB.SMembers(setName)
	return result.Err(), result.Val()
}

// SScan 模糊搜索
func SScan(key string, cursor uint64, match string, count int64) ([]string, uint64) {
	result := RDB.SScan(key, cursor, match, count)
	return result.Val()
}

// SIsMember 判断给定member是否在集合里面
func SIsMember(key string, members interface{}) (error, bool) {
	result := RDB.SIsMember(key, members)
	return result.Err(), result.Val()
}
