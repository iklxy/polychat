package service

import (
	"errors"
	"polychat/internal/dao"
	"polychat/internal/model"
	"polychat/pkg/util"

	"gorm.io/gorm"
)

type UserService struct{}

// 注册
func (s *UserService) Register(username, password string) error {
	//检查用户名是否存在
	_, err := dao.GetUserByUsername(username)
	if err == nil {
		return errors.New("用户名已存在")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err // 其他数据库错误
	}

	//加密密码
	hashPassword, err := util.HashPassword(password)
	if err != nil {
		return err //密码加密的时候出错
	}

	//创建新用户
	user := &model.User{
		Username: username,
		Password: hashPassword,
	}

	return dao.CreateUser(user)
}

// 登录
func (s *UserService) Login(username, password string) (string, error) {
	//先查询用户是否存在
	user, err := dao.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("用户名不存在") //用户名不存在
	}

	//用户存在的前提下，校验密码
	if !util.CheckPassword(password, user.Password) {
		return "", errors.New("密码错误")
	}

	//密码校验通过，返回token
	token, err := util.GenerateToken(user.ID)
	if err != nil {
		return "", errors.New("token生成失败")
	}
	return token, nil
}
