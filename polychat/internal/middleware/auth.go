package middleware

import (
	"net/http"
	"polychat/pkg/util"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized,
				gin.H{"code": "401", "msg": "未授权"})
		}
		// 检查Authorization字段是否符合Bearer格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized,
				gin.H{"code": "402", "msg": "格式错误"})
		}
		//校验token是否有效
		token := parts[1]
		claims, err := util.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized,
				gin.H{"code": "403", "msg": "token无效"})
		}
		//将claims中的用户信息设置到上下文
		c.Set("claims", claims)
		c.Next()
	}
}
