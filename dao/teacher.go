package dao

import "course-selection/model"

func SelectTeacher(teacherNumber string) error {
	var teaching model.Teaching
	result := DB.Select("teacherName").Where("teacherNumber = ?", teacherNumber).Take(&teaching)
	return result.Error
}
