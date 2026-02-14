package database

import (
	"fmt"
	"os"
	"polychat/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 操作数据库的主要对象
var DB *gorm.DB

func InitDB() {
	// 获取数据库Host，默认为远程IP（用于本地开发），部署到服务器时可通过环境变量 DB_HOST=127.0.0.1 指定
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "47.110.94.115"
	}

	dsn := fmt.Sprintf("admin:YY010303@tcp(%s:3306)/polychat_db?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s&allowNativePasswords=true&tls=false", dbHost)

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
