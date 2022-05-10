package service

import (
	"course-selection/dao"
	"course-selection/model"
)

func InsertTeacher(teacher model.Teacher) error {
	return dao.InsertTeachers(teacher)
}
