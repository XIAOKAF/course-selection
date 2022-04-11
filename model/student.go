package model

type Student struct {
	StudentId   int `gorm:"primary_key"`
	UnifiedCode string
	StudentName string
	Gender      int
	Grade       string
	Class       string
	Password    string
	Department  string
	Major       string
	Image       string
}
