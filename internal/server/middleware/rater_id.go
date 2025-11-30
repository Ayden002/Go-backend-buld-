package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRaterID 验证X-Rater-Id请求头
func (m *MiddlewareSet) RequireRaterID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取X-Rater-Id头
		raterID := c.GetHeader("X-Rater-Id")
		if raterID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Missing or invalid authentication information",
			})
			c.Abort()
			return
		}

		// 将rater_id存储在context中供handler使用
		c.Set("rater_id", raterID)
		c.Next()
	}
}
