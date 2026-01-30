package dao

import (
	"polychat/internal/model"
	"polychat/pkg/database"
)

// 调用数据库中的DB来创建新用户
func CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

// 根据用户名来查询用户
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
