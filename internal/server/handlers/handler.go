package handlers

import (
	"github.com/jmoiron/sqlx"
	"interview/internal/boxoffice"
	"interview/internal/config"
)

type HandlerSet struct {
	DB        *sqlx.DB
	BoxOffice *boxoffice.Client
}

func NewHandlerSet(db *sqlx.DB, cfg config.Config) *HandlerSet {
	// 初始化BoxOffice客户端
	boxOfficeClient := boxoffice.New(cfg.BoxofficeURL, cfg.BoxofficeAPIKey)

	return &HandlerSet{
		DB:        db,
		BoxOffice: boxOfficeClient,
	}
}
