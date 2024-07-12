package handlers

import (
	"argus/internal/db"
	"argus/pkg/logger"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"slices"
	"strconv"
	"time"
)

// Constants for default values and sorting orders.
const (
	AgentsDefaultPage     = 1
	AgentsDefaultPageSize = 10
)

const (
	Asc  = "asc"
	Desc = "desc"
)

var (
	ValidAgentSorts  = []string{"id"}      // ValidAgentSorts defines valid fields for sorting agents.
	ValidAgentOrders = []string{Asc, Desc} // ValidAgentOrders defines valid sorting orders.
)

// HandleCreateAgent handles requests to create a new agent
// @Summary Create a new agent
// @Description Create a new agent with the provided IP address and retrieve its details
// @Tags agents
// @Accept json
// @Produce json
// @Param request body CreateAgentRequest true "Request body for creating a new agent"
// @Success 201 {object} CreateAgentResponse "Successfully created agent"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /agents [post]
func (gh *GinHandler) HandleCreateAgent(ctx *gin.Context) {
	var createRequest CreateAgentRequest
	if err := ctx.ShouldBindJSON(&createRequest); err != nil {
		logger.WithError(err).Debug("cannot parse create new agent request")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "cannot parse request body"})
		return
	}

	// Validate the input data
	err := createRequest.validate()
	if err != nil {
		logger.WithError(err).Debug("cannot validate create new agent request")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Create a context timeout to circuit break in case of long API call
	getIPInfoCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	// Make a call to IPInfo to get the stats about IP
	stats, err := gh.ipStatsGatherer.GetInfo(getIPInfoCtx, createRequest.IPAddress)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.WithField("ip", createRequest.IPAddress).WithError(err).Warn("timeout exceeded for IPInfo API")
			ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{Error: "timeout exceeded"})
			return
		}

		logger.WithField("ip", createRequest.IPAddress).WithError(err).Warn("cannot gather statistics for this IP address")
		ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{Error: "cannot gather statistics for this IP address"})
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
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "cannot create agent"})
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

// HandleGetAgents handles retrieving agents
// @Summary Get a list of agents
// @Description Retrieve a list of agents based on optional query parameters
// @Tags agents
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination (default is 1)"
// @Param page_size query int false "Number of agents per page (default is 10)"
// @Param ip_address query string false "Filter agents by IP address"
// @Param sort_by query string false "Field to sort agents by (e.g., 'id')"
// @Param order query string false "Sorting order ('asc' or 'desc')"
// @Success 200 {object} GetAgentsResponse "Successfully retrieved agents"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 404 {object} GetAgentsResponse "No agents found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /agents [get]
func (gh *GinHandler) HandleGetAgents(ctx *gin.Context) {
	// Handle query params
	var queryParams GetAgentsQueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		logger.WithError(err).Debug("cannot bind query params")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "bad query params"})
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
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "the valid fields are: id",
			})
			return
		}
		agentSort.SortBy = &queryParams.SortBy
	}
	if queryParams.Order != "" {
		if !slices.Contains(ValidAgentOrders, queryParams.Order) {
			logger.WithField("order", queryParams.Order).Debug("cannot parse the order parameter")
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "the valid fields are: asc, desc",
			})
			return
		}
		agentSort.OrderBy = &queryParams.Order
	}

	// Retrieve agents from database
	agentsResult, err := gh.db.GetAllAgents(ctx, agentsFilter, queryParams.Page, queryParams.PageSize, agentSort)
	if err != nil {
		logger.WithError(err).Warn("cannot retrieve the agents from the database")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "cannot retrieve the agents from the database"})
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

// HandleGetAgentDetail handles getting details about each agent
// @Summary Get details of a specific agent
// @Description Retrieve detailed information of a specific agent by ID
// @Tags agents
// @Accept json
// @Produce json
// @Param agent_id path int true "ID of the agent to retrieve"
// @Success 200 {object} AgentDetailedResponse "Successfully retrieved agent details"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 404 {object} ErrorResponse "Agent not found"
// @Router /agents/{agent_id} [get]
func (gh *GinHandler) HandleGetAgentDetail(ctx *gin.Context) {
	// Parsing the agent_id
	agentIDParam, err := strconv.Atoi(ctx.Param("agent_id"))
	if err != nil {
		logger.WithError(err).Warn("cannot parse agent id")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "agent_id is not provided or is not valid",
		})
		return
	}
	agentID := uint(agentIDParam)

	// Retrieve the agent from the database
	agent, err := gh.db.GetAgentByID(ctx, agentID)
	if err != nil {
		logger.WithError(err).Warn("cannot retrieve agent by id")
		ctx.JSON(http.StatusNotFound, ErrorResponse{Error: "cannot find such agent by id"})
		return
	}

	ctx.JSON(http.StatusOK, AgentDetailedResponse{
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
