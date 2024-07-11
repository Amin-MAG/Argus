package handlers

import (
	"argus/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ServiceUp   = "up"
	ServiceDown = "down"
)

func (gh *GinHandler) Ping(ctx *gin.Context) {
	dbStatus := ServiceUp
	err := gh.db.Ping(ctx)
	if err != nil {
		logger.WithError(err).Warn("cannot access underlying database")
		dbStatus = ServiceDown
	}

	ctx.JSON(http.StatusOK, gin.H{
		"database_status": dbStatus,
	})
}
