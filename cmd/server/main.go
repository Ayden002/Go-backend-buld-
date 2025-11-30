package main

import (
	"fmt"
	"interview/internal/config"
	"interview/internal/db"
	"interview/internal/server"
	"interview/internal/server/handlers"
	"interview/internal/server/middleware"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	//临时export环境变量用的（生产环境会通过docker-compose注入）
	_ = godotenv.Load()
	//从环境变量加载配置
	cfg := config.Load()
	//连接数据库
	conn := db.Connect(cfg.DBURL)

	// 运行数据库迁移
	db.Migrate(conn)

	// 初始化handler和middleware
	h := handlers.NewHandlerSet(conn, cfg)
	m := middleware.NewMiddlewareSet(cfg)

	// 生成Gin router
	r := server.NewRouter(server.RouterDeps{
		Handlers:   h,
		Middleware: m,
	})

	// 4. 启动HTTP server
	addr := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	log.Printf("server listening on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
