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

// Ping checks for health of the system
// @Summary Check health status
// @Description Check if the health of system is ok or not
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {object} PingResponse
// @Router /ping [get]
func (gh *GinHandler) Ping(ctx *gin.Context) {
	dbStatus := ServiceUp
	err := gh.db.Ping(ctx)
	if err != nil {
		logger.WithError(err).Warn("cannot access underlying database")
		dbStatus = ServiceDown
	}

	ctx.JSON(http.StatusOK, PingResponse{
		DatabaseStatus: dbStatus,
	})
}
