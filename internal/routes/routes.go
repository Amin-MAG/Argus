package routes

import (
	"argus/config"
	"argus/internal/db"
	"argus/internal/handlers"
	"argus/internal/iputil"
	"argus/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zsais/go-gin-prometheus"
	"net/http"
)

const ApiV1 = "/api/v1"

// NewGinServer creates a new Server instance.
func NewGinServer(cfg config.Config, db db.DB, ipStatsGatherer iputil.IPStatsGatherer) (*http.Server, error) {
	// Gin Configuration
	gin.SetMode(cfg.Argus.GinMode)

	// Create new engine for the server
	engine := gin.Default()

	// Set up the middlewares
	// TODO: Add middlewares here
	p := ginprometheus.NewPrometheus("gin")
	p.Use(engine)

	// Create the HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Argus.Port),
		Handler: engine,
	}
	logger.Info("new gin server has been created")

	// Create new Gin handler
	ginHandler := handlers.NewGinHandler(cfg, db, ipStatsGatherer)

	// Register routes of modules
	v1 := engine.Group(ApiV1)
	// Health APIs
	v1.GET("/health/ping", ginHandler.Ping)
	// AgentDetailedResponse Monitoring APIs
	v1.POST("/agents", ginHandler.HandleCreateAgent)
	v1.GET("/agents", ginHandler.HandleGetAgents)
	v1.GET("/agents/:agent_id", ginHandler.HandleGetAgentDetail)

	return server, nil
}
