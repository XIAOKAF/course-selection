package dao

import "course-selection/model"

func SelectCourse(courseNumber string) error {
	var courseName string
	result := DB.Table("course").Select("courseName").Scan(&courseName)
	return result.Error
}

func CreateCourse(course model.Course) error {
	result := DB.Select("courseNumber", "courseName", "courseDepartment", "courseCredit", "courseType", "courseGrade", "duration").Create(&course)
	return result.Error
}

func ChooseCourse(choice model.Choice) error {
	result := DB.Select("teachingClass", "unifiedCode").Create(&choice)
	return result.Error
}
