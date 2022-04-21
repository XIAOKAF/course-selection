package dao

import "course-selection/model"

func InsertCourse(course model.Course) error {
	result := DB.Select("courseNumber", "courseName", "courseDepartment", "courseCredit", "courseType", "courseGrade", "duration").Create(&course)
	return result.Error
}

func ChooseCourse(choice model.Choice) error {
	result := DB.Select("teachingClass", "unifiedCode").Create(&choice)
	return result.Error
}
