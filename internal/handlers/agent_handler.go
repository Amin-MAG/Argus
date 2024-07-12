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

	// Make a call to IPInfo to get the stats about IP
	stats, err := gh.ipStatsGatherer.GetInfo(ctx, createRequest.IPAddress)
	if err != nil {
		logger.WithField("ip", createRequest.IPAddress).WithError(err).Warn("cannot gather statistics for this IP address")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "cannot gather statistics for this IP address"})
		return
	}

	// Create the row in database
	agent, err := gh.db.CreateNewAgent(ctx, &db.Agent{
		IPAddress: stats.IP.String(),
		ASN:       stats.ASN,
		ISP:       stats.ISP,
		City:      stats.City,
		Country:   stats.Country,
		Location:  stats.Location,
	})
	if err != nil {
		logger.WithError(err).Warn("cannot create agent")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create agent"})
		return
	}

	ctx.JSON(http.StatusCreated, CreateAgentResponse{
		Message: "agent has been created successfully",
		Agent: Agent{
			ID:        agent.ID,
			CreatedAt: agent.CreatedAt,
			IPAddress: agent.IPAddress,
			ASN:       agent.ASN,
			ISP:       agent.ISP,
			City:      agent.City,
			Country:   agent.Country,
			Location:  agent.Location,
		},
	})
}

func (gh *GinHandler) HandleGetAgents(ctx *gin.Context) {
	// Handle query params
	var queryParams GetAgentsQueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		logger.WithError(err).Debug("cannot bind query params")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad query params"})
		return
	}
	if queryParams.Page == 0 {
		queryParams.Page = AgentsDefaultPage
	}
	if queryParams.PageSize == 0 {
		queryParams.PageSize = AgentsDefaultPageSize
	}
	agentsFilter := &db.AgentFilter{}
	if queryParams.IPAddress != "" {
		agentsFilter.IPAddress = &queryParams.IPAddress
	}
	agentSort := &db.AgentSort{}
	if queryParams.SortBy != "" {
		if !slices.Contains(ValidAgentSorts, queryParams.SortBy) {
			logger.WithField("sort_by", queryParams.SortBy).Debug("cannot parse the sort by parameter")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "the valid fields are: id",
			})
			return
		}
		agentSort.SortBy = &queryParams.SortBy
	}
	if queryParams.Order != "" {
		if !slices.Contains(ValidAgentOrders, queryParams.Order) {
			logger.WithField("order", queryParams.Order).Debug("cannot parse the order parameter")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "the valid fields are: asc, desc",
			})
			return
		}
		agentSort.OrderBy = &queryParams.Order
	}

	// Retrieve agents from database
	agentsResult, err := gh.db.GetAllAgents(ctx, agentsFilter, queryParams.Page, queryParams.PageSize, agentSort)
	if err != nil {
		logger.WithError(err).Warn("cannot retrieve the agents from the database")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot retrieve the agents from the database"})
		return
	}
	agentStats := AgentPagination{
		TotalAgents: agentsResult.TotalAgents,
		TotalPages:  int(math.Ceil(float64(agentsResult.TotalAgents) / float64(queryParams.PageSize))),
		CurrentPage: queryParams.Page,
		PerPage:     queryParams.PageSize,
	}

	// Handle no agent found
	if len(agentsResult.Agents) == 0 {
		ctx.JSON(http.StatusNotFound, GetAgentsResponse{
			Message: "there is no agents for this page",
			Data: AgentsData{
				Agents:     []Agent{},
				Pagination: agentStats,
			},
		})
		return
	}

	// Converting the agents to response model
	var agents []Agent
	for _, a := range agentsResult.Agents {
		agents = append(agents, Agent{
			ID:        a.ID,
			IPAddress: a.IPAddress,
			CreatedAt: a.CreatedAt,
			ASN:       a.ASN,
		})
	}

	ctx.JSON(http.StatusOK, GetAgentsResponse{
		Message: "retrieved agents successfully",
		Data: AgentsData{
			Agents:     agents,
			Pagination: agentStats,
		},
	})
}

func (gh *GinHandler) HandleGetAgentDetail(ctx *gin.Context) {
	// Parsing the agent_id
	agentIDParam, err := strconv.Atoi(ctx.Param("agent_id"))
	if err != nil {
		logger.WithError(err).Warn("cannot parse agent id")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "agent_id is not provided or is not valid",
		})
		return
	}
	agentID := uint(agentIDParam)

	// Retrieve the agent from the database
	agent, err := gh.db.GetAgentByID(ctx, agentID)
	if err != nil {
		logger.WithError(err).Warn("cannot retrieve agent by id")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "cannot find such agent by id"})
		return
	}

	ctx.JSON(http.StatusCreated, AgentDetailedResponse{
		Message: "agent has been retrieved successfully",
		Agent: Agent{
			ID:        agent.ID,
			CreatedAt: agent.CreatedAt,
			IPAddress: agent.IPAddress,
			ASN:       agent.ASN,
			ISP:       agent.ISP,
			City:      agent.City,
			Country:   agent.Country,
			Location:  agent.Location,
		},
	})
}
