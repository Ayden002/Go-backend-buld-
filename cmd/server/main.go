package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"interview/internal/config"
	"interview/internal/server"
	"interview/internal/server/handlers"
	"interview/internal/server/middleware"
	"log"
)

func main() {

	// 临时export环境变量用的（生产环境会通过docker-compose注入）
	_ = godotenv.Load()
	// 1.从环境变量加载配置
	cfg := config.Load()

	// 2. 初始化handler和middleware
	h := handlers.NewHandlerSet()
	m := middleware.NewMiddlewareSet(cfg)

	// 3. 生成Gin router
	r := server.NewRouter(server.RouterDeps{
		Handlers:   h,
		Middleware: m,
	})

	// 4. 启动HTTP server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("server listening on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
