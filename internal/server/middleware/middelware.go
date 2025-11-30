package middleware

import "interview/internal/config"

type MiddlewareSet struct {
	Config config.Config
}

func NewMiddlewareSet(cfg config.Config) *MiddlewareSet {
	return &MiddlewareSet{Config: cfg}
}
