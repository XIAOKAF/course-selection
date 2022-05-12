package dao

import "course-selection/model"

func AdministratorLogin(administrator model.Administrator) (error, string) {
	result := DB.Where("administrator_id = ?", administrator.AdministratorId).First(&administrator)
	if result.Error != nil {
		return result.Error, administrator.Password
	}
	return nil, administrator.Password
}

// Cancel 注销学生账号
func Cancel(studentId string) error {
	var student model.Student
	result := DB.Where("student_id = ?", studentId).Delete(&student)
	return result.Error
}
