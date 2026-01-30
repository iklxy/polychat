package model

import (
	"gorm.io/gorm"
)

type User struct {
	/*
		type Model struct {
			ID        uint `gorm:"primarykey"`
			CreatedAt time.Time
			UpdatedAt time.Time
			DeletedAt DeletedAt `gorm:"index"`
		}*/
	gorm.Model
	//uniqueIndex : 为数据库建立唯一索引，防止用户名重复
	Username string `gorm:"uniqueIndex;type:varchar(20);not null"`
	//密码允许重复
	Password string `gorm:"type:varchar(100);not null"` // 存加密之后的哈希值
	Email    string `gorm:"type:varchar(100)"`
	Avatar   string `gorm:"type:varchar(255)"` // 头像URL
}
