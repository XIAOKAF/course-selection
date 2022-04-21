package service

import (
	"course-selection/dao"
	"course-selection/model"
	"gorm.io/gorm"
)

func SelectUnifiedCode(unifiedCode string) (bool, error, string) {
	err, pwd := dao.SelectUnifiedCode(unifiedCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil, pwd
		}
		return true, err, pwd
	}
	return true, nil, pwd
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
