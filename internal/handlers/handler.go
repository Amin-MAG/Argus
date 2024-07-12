package handlers

import (
	"argus/config"
	"argus/internal/db"
	"argus/internal/iputil"
)

type GinHandler struct {
	cfg             config.Config
	db              db.DB
	ipStatsGatherer iputil.IPStatsGatherer
}

func NewGinHandler(cfg config.Config, db db.DB, ipStatsGatherer iputil.IPStatsGatherer) *GinHandler {
	return &GinHandler{cfg: cfg, db: db, ipStatsGatherer: ipStatsGatherer}
}
