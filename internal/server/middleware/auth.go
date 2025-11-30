package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 验证Bearer Token
func (m *MiddlewareSet) AuthBearer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Missing or invalid authentication information",
			})
			c.Abort()
			return
		}

		// 检查格式 "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Missing or invalid authentication information",
			})
			c.Abort()
			return
		}

		// 验证token
		token := parts[1]
		if token != m.Config.AuthToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Missing or invalid authentication information",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
