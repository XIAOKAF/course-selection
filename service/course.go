package service

import (
	"course-selection/dao"
	"course-selection/model"
	"gorm.io/gorm"
)

// SelectCourse 查询课程是否已经存在
func SelectCourse(courseNumber string) (bool, error) {
	err := dao.SelectCourse(courseNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

// CreateCourse mysql插入新的课程信息
func CreateCourse(course model.Course) error {
	err := dao.CreateCourse(course)
	return err
}

// RCreateCourse redis插入新的课程信息
func RCreateCourse(course model.Course) error {
	courseMap := make(map[string]interface{})
	courseMap["courseName"] = course.CourseName
	courseMap["courseDepartment"] = course.CourseDepartment
	courseMap["courseCredit"] = course.CourseCredit
	courseMap["courseType"] = course.CourseType
	courseMap["duration"] = course.Duration
	err := dao.HashSet(course.CourseNumber, courseMap)
	return err
}

func ChooseCourse(choice model.Choice) error {
	err := dao.ChooseCourse(choice)
	return err
}
