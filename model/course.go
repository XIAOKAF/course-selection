package model

type Course struct {
	CourseNumber     string
	CourseName       string
	CourseDepartment string
	CourseCredit     float64
	CourseType       int
	CourseGrade      string
	Duration         string
}

type Selection struct {
	CourseNumber  string
	TeachingClass string
	TeacherName   string
	SetTime       string
	CourseCredit  string
	CourseType    string
}

type TeachingClassInfo struct {
	*Course
	TeachingClassNumber string
	SetTime             string
	StudentSum          int
}

type ClassDetails struct {
	*TeachingClassInfo
	TeacherName string
}
