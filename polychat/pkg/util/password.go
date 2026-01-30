package util

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 对密码进行加密
// 返回加密后的哈希字符串
func HashPassword(password string) (string, error) {
	// DefaultCost 默认为 10
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 比对明文密码和数据库中的哈希值是否匹配
// 返回 true 表示匹配成功
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
