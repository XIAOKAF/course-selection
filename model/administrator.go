package model

type Administrator struct {
	AdministratorId string `gorm:"primary_key"`
	Password        string
	RuleLevel       string
}
