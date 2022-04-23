package model

type Course struct {
	CourseId         string
	CourseNumber     string
	CourseName       string
	CourseDepartment string
	CourseCredit     float64
	CourseType       int
	TeachingClass    string
	CourseGrade      string
	Duration         string
}

type Teaching struct {
	CourseNumber  string
	TeachingClass string
	TeacherNumber string
	SetTime       string
}

type Choice struct {
	ChoiceId      int
	TeachingClass string
	UnifiedCode   string
	StudentName   string
}

type Selection struct {
	CourseNumber  string
	TeachingClass string
	TeacherName   string
	SetTime       string
	CourseCredit  string
	CourseType    string
}
