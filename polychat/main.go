package main

import (
	"fmt"
	"polychat/internal/api"
	"polychat/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化数据库连接
	database.InitDB()

	// 2.初始化gin引擎
	r := gin.Default()

	// 静态文件服务
	r.Static("/static", "./static")
	// 将根路径 / 直接映射到 index.html
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	//3.注册路由
	userHandle := api.UserHandle{}

	//路由组
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", userHandle.Register)
		v1.POST("/login", userHandle.Login)
	}
	fmt.Println("服务器运行在8080端口")
	if err := r.Run(":8080"); err != nil {
		panic("运行失败" + err.Error())
	}
}
