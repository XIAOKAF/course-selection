package service

import (
	"course-selection/dao"
	"course-selection/model"
)

func SpiderMan(student model.Student) error {
	return dao.SpiderMan(student)
}
