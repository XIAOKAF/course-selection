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
