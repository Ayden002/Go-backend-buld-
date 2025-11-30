package middleware

import "github.com/gin-gonic/gin"

func (m *MiddlewareSet) RequireRaterID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO
		c.Next()
	}
}
