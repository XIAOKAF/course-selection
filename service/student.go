package service

import (
	"course-selection/dao"
	"course-selection/model"
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

// UpdatePassword 更改用户密码
func UpdatePassword(student model.Student) error {
	err := dao.UpdatePassword(student)
	if err != nil {
		return err
	}
	return nil
}

// UpdateMobile 更新用户电话号码
func UpdateMobile(student model.Student) error {
	err := dao.UpdateMobile(student)
	if err != nil {
		return err
	}
	return nil
}
