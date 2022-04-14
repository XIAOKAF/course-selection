package service

import "course-selection/dao"

// SelectRoleLevel 查询用户权限等级
func SelectRoleLevel(userNumber string) (error, string) {
	err, roleLevel := dao.Get(userNumber)
	if err != nil {
		return err, roleLevel
	}
	return nil, roleLevel
}
