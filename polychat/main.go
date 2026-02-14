package main

import (
	"fmt"
	"polychat/internal/api"
	"polychat/internal/middleware"
	"polychat/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化数据库连接
	database.InitDB()

	gin.SetMode(gin.ReleaseMode)
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
	RelationHandle := api.RelationHandler{}
	//公开接口，不需要Token验证
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", userHandle.Register)
		v1.POST("/login", userHandle.Login)
	}

	//受保护的接口，需要验证Token
	authorized := v1.Group("/")
	authorized.Use(middleware.JWTAuthMiddleware())
	{
		authorized.GET("/chat", api.ConnectWS)

		// 好友关系模块
		relationGroup := authorized.Group("/relation")
		{
			relationGroup.POST("/add", RelationHandle.AddFriend)
			relationGroup.POST("/delete", RelationHandle.DeleteFriend)
			relationGroup.GET("/list", RelationHandle.GetFriend)
			relationGroup.POST("/update_note", RelationHandle.UpdateFriendNote)
		}
	}
	fmt.Println("服务器运行在8080端口")
	if err := r.Run(":8080"); err != nil {
		panic("运行失败" + err.Error())
	}
}
