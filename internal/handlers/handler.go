package handlers

import (
	"argus/config"
	"argus/internal/db"
)

type GinHandler struct {
	cfg config.Config
	db  db.DB
}

func NewGinHandler(cfg config.Config, db db.DB) *GinHandler {
	return &GinHandler{cfg: cfg, db: db}
}
