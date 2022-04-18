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

func TTL(key string) (time.Duration, error) {
	result := RDB.TTL(key)
	return result.Result()
}
