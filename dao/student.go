package dao

import "course-selection/model"

func SpiderMan(student model.Student) error {
	result := DB.Select("student_id", "student_name", "gender", "grade", "class", "department", "major", "rule_id").Create(&student)
	return result.Error
}
