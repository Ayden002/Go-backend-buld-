package server

import (
	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
	"interview/internal/server/handlers"
	"interview/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	Handlers   *handlers.HandlerSet
	Middleware *middleware.MiddlewareSet
}

func NewRouter(deps RouterDeps) *gin.Engine {
	// 使用 gin.New() 而不是 gin.Default()，这样可以避免默认的错误处理
	r := gin.New()

	// 手动添加必要的中间件
	r.Use(gin.Logger())
	// 全局错误处理
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.AbortWithStatusJSON(500, gin.H{"code": "INTERNAL_ERROR", "message": "Internal server error"})
	}))

	// 健康检查
	r.GET("/healthz", deps.Handlers.Healthz)

	// 鉴权中间件
	authRequired := deps.Middleware.AuthBearer()
	raterIDRequired := deps.Middleware.RequireRaterID()

	// 电影路由
	// GET /movies - 列表和搜索（公开）
	r.GET("/movies", deps.Handlers.ListMovies)

	// POST /movies - 创建电影（需要鉴权）
	r.POST("/movies", authRequired, deps.Handlers.CreateMovie)

	// GET /movies/:title/rating - 获取评分聚合（公开）
	r.GET("/movies/:title/rating", deps.Handlers.GetRating)

	// POST /movies/:title/ratings - 提交评分（需要X-Rater-Id）
	r.POST("/movies/:title/ratings", raterIDRequired, deps.Handlers.SubmitRating)
	// Swagger 文档
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r

}
