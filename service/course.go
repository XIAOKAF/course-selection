package service

import (
	"course-selection/dao"
	"course-selection/model"
	"sort"
)

// SelectCourse 查询课程是否已经存在
func SelectCourse(courseNumber string) (error, bool) {
	return dao.SIsMember("course", courseNumber)
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
	courseMap["courseGrade"] = course.CourseGrade
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

// JudgeTimeConflict 判断所选课程是否存在时间冲突
func JudgeTimeConflict(selectedTimeArr []string, timeArr []string) bool {
	sort.Strings(selectedTimeArr)
	for _, v := range timeArr {
		index := sort.SearchStrings(selectedTimeArr, v)
		if index < len(selectedTimeArr) && selectedTimeArr[index] == v {
			return true
		}
	}
	return false
}

// IsClassExist 判断教学班是否已经被创建
func IsClassExist(classArr []string, class string) bool {
	sort.Strings(classArr)
	index := sort.SearchStrings(classArr, class)
	if index < len(classArr) && classArr[index] == class {
		return true
	}
	return false
}
