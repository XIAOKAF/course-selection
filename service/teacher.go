package service

import (
	"course-selection/dao"
	"gorm.io/gorm"
)

func SelectTeacher(teacherNumber string) (bool, error) {
	err := dao.SelectTeacher(teacherNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return true, err
	}
	return true, nil
}
