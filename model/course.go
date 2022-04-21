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

type Choice struct {
	ChoiceId      int
	TeachingClass string
	UnifiedCode   string
}
