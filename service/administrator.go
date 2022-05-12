package service

import (
	"course-selection/dao"
	"course-selection/model"
)

func AdministratorLogin(administrator model.Administrator) (error, string) {
	err, password := dao.AdministratorLogin(administrator)
	if err != nil {
		return err, password
	}
	return nil, password
}

func Cancel(studentId string) error {
	err := dao.Cancel(studentId)
	return err
}
