package handlers

import (
	"argus/internal/db"
	"argus/pkg/logger"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"slices"
	"strconv"
)

const (
	AgentsDefaultPage     = 1
	AgentsDefaultPageSize = 10
)

const (
	Asc  = "asc"
	Desc = "desc"
)

var (
	ValidAgentSorts  = []string{"id"}
	ValidAgentOrders = []string{Asc, Desc}
)

func (gh *GinHandler) HandleCreateAgent(ctx *gin.Context) {
	var createRequest CreateAgentRequest
	if err := ctx.ShouldBindJSON(&createRequest); err != nil {
		logger.WithError(err).Debug("cannot parse create new agent request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot parse request body"})
		return
	}

	// Validate the input data
	err := createRequest.validate()
	if err != nil {
		logger.WithError(err).Debug("cannot validate create new agent request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the row in database
	agent, err := gh.db.CreateNewAgent(ctx, &db.Agent{
		IPAddress: createRequest.IPAddress,
		ASN:       "dum",
		ISP:       "dum",
	})
	if err != nil {
		logger.WithError(err).Warn("cannot create agent")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create agent"})
		return
	}

	ctx.JSON(http.StatusCreated, CreateAgentResponse{
		Message: "agent has been created successfully",
		Agent: AgentDetailedResponse{
			ID:        agent.ID,
			CreatedAt: agent.CreatedAt,
			IPAddress: agent.IPAddress,
			ASN:       agent.ASN,
			ISP:       agent.ISP,
		},
	})
}

func (gh *GinHandler) HandleGetAgents(ctx *gin.Context) {
}

func (gh *GinHandler) HandleGetAgentDetail(ctx *gin.Context) {
}
