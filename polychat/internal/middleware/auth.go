package middleware

import (
	"net/http"
	"polychat/pkg/util"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		// 从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		// 如果有Authorization字段，检查Bearer格式
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{"code": "402", "msg": "格式错误"})
				c.Abort()
				return
			}
			token = parts[1]
		} else {
			// 如果没有Authorization字段，尝试从查询参数中获取token
			token = c.Query("token")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "403", "msg": "未提供token"})
			c.Abort()
			return
		}

		// 校验token是否有效
		claims, err := util.ParseToken(token)
		if err != nil {
			// 添加错误日志输出
			// fmt.Println("Token解析失败:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"code": "403", "msg": "token无效: " + err.Error()})
			c.Abort()
			return
		}
		//将claims中的用户信息设置到上下文
		c.Set("userID", claims.UserID) // 便利后续使用
		c.Set("claims", claims)
		c.Next()
	}
}
