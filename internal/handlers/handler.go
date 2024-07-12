package handlers

import (
	"argus/config"
	"argus/internal/db"
	"argus/internal/iputil"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GinHandler struct {
	cfg             config.Config
	db              db.DB
	ipStatsGatherer iputil.IPStatsGatherer
}

func NewGinHandler(cfg config.Config, db db.DB, ipStatsGatherer iputil.IPStatsGatherer) *GinHandler {
	return &GinHandler{cfg: cfg, db: db, ipStatsGatherer: ipStatsGatherer}
}

func (gh *GinHandler) BillingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.Header.Get("API-Key")
		// TODO: should be defined for each agent in a separate table (Out of scope of this project)
		if key != gh.cfg.Argus.APIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
		}
		c.Next()
	}
}
