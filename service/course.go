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

func RDetailsCourse(teaching model.Teaching) error {
	teachingMap := make(map[string]interface{})
	teachingMap["courseNumber"] = teaching.CourseNumber
	teachingMap["setTime"] = teaching.SetTime
	teachingMap["teacherNumber"] = teaching.TeacherNumber
	err := dao.HashSet(teaching.TeachingClass, teachingMap)
	return err
}

// IsRepeated 判断时间是否有重叠的部分
func IsRepeated(selectArr, choice []string) bool {
	//有时间冲突则返回true，反之false
	for _, value := range selectArr {
		for _, v := range choice {
			if value == v {
				return true
			}
		}
	}
	return false
}

func ChooseCourse(choice model.Choice) error {
	err := dao.ChooseCourse(choice)
	return err
}
