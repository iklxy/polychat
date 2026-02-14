package database

import (
	"fmt"
	"polychat/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 操作数据库的主要对象
var DB *gorm.DB

func InitDB() {
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "admin:YY010303@tcp(127.0.0.1:3306)/polychat_db?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s&allowNativePasswords=true&tls=false"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		//在err不为空的时候，说明连接失败，抛出异常且终止流程
		panic("数据库连接失败" + err.Error())
	}

	//通过model包中的User结构体，自动创建数据库中的user表
	err = DB.AutoMigrate(&model.User{})
	err = DB.AutoMigrate(&model.Relation{})
	if err != nil {
		//在err不为空的时候，说明创建表失败，抛出异常且终止流程
		panic("数据库创建表失败" + err.Error())
	}
	fmt.Println("数据库连接成功")
}
