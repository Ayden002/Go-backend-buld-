package server

import (
	"interview/internal/server/handlers"
	"interview/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	Handlers   *handlers.HandlerSet
	Middleware *middleware.MiddlewareSet
}

func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/healthz", deps.Handlers.Healthz)

	//all路由TODO
	//

	return r
}
