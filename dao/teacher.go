package dao

import "course-selection/model"

func InsertTeachers(teacher model.Teacher) error {
	result := DB.Select("teacher_id", "teacher_number", "teacher_name", "rule_level").Create(&teacher)
	return result.Error
}
