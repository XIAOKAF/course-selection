package dao

import "course-selection/model"

func SelectUnifiedCode(unifiedCode string) (error, string) {
	var pwd string
	result := DB.Select("student").Where("unifiedCode = ?", unifiedCode).Take(&pwd)
	return result.Error, pwd
}

// UpdatePassword 修改密码
func UpdatePassword(student model.Student) error {
	result := DB.Model(&student).Where("unifiedCode = ?", student.UnifiedCode).Update("password", student.Password)
	return result.Error
}

func UpdateMobile(student model.Student) error {
	result := DB.Model(&student).Where("unifiedCode = ?", student.UnifiedCode).Update("mobile", student.Mobile)
	return result.Error
}
