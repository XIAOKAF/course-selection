package dao

import "course-selection/model"

func AdministratorLogin(administrator model.Administrator) (error, string) {
	result := DB.Where("administratorId = ?", administrator.AdministratorId).First(&administrator)
	if result.Error != nil {
		return result.Error, administrator.Password
	}
	return nil, administrator.Password
}
