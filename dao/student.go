package dao

import "course-selection/model"

func SelectStudentByUnifiedCode(unifiedCode string) error {
	_, err := Get(unifiedCode)
	if err != nil {
		return err
	}
	return nil
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
