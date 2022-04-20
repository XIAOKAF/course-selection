package service

import (
	"course-selection/dao"
	"course-selection/model"
)

// InsertCourse 插入新的课程信息
func InsertCourse(course model.Course) error {
	err := dao.InsertCourse(course)
	return err
}
