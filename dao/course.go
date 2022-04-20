package dao

import "course-selection/model"

func InsertCourse(course model.Course) error {
	result := DB.Select("courseNumber", "courseName", "courseDepartment", "courseCredit", "courseType", "teacher", "teachingClass", "courseGrade", "duration").Create(&course)
	return result.Error
}
