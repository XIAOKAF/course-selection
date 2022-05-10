package dao

import "course-selection/model"

func CreateCourse(course model.Course) error {
	result := DB.Select("course_number", "course_name", "course_department", "course_grade", "course_credit", "course_type", "duration").Create(&course)
	return result.Error
}
