package service

import (
	"course-selection/dao"
	"github.com/go-redis/redis"
)

func SelectStudentByUnifiedCode(unifiedCode string) (bool, error) {
	err := dao.SelectStudentByUnifiedCode(unifiedCode)
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return true, err
	}
	return true, nil
}
