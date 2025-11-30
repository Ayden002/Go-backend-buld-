package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /healthz
func (h *HandlerSet) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
