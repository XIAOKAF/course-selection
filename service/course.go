package service

import (
	"course-selection/dao"
	"course-selection/model"
)

// InsertCourse mysql插入新的课程信息
func InsertCourse(course model.Course) error {
	err := dao.InsertCourse(course)
	return err
}

// RInsertCourse redis插入新的课程信息
func RInsertCourse(course model.Course) error {
	courseMap := make(map[string]interface{})
	courseMap["courseName"] = course.CourseName
	courseMap["courseDepartment"] = course.CourseDepartment
	courseMap["courseCredit"] = course.CourseCredit
	courseMap["courseType"] = course.CourseType
	courseMap["setTime"] = course.SetTime
	courseMap["duration"] = course.Duration
	err := dao.HashSet(course.CourseNumber, courseMap)
	return err
}
