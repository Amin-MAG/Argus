package handlers

import (
	"argus/config"
	"argus/internal/db"
	"argus/internal/iputil"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCreateAgent_Success(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)

	testData := struct {
		requestBody      CreateAgentRequest
		expectedCode     int
		expectedResponse CreateAgentResponse
	}{
		requestBody: CreateAgentRequest{
			IPAddress: "192.168.1.1",
		},
		expectedCode: http.StatusCreated,
		expectedResponse: CreateAgentResponse{
			Message: "agent has been created successfully",
			Agent: Agent{
				IPAddress: "8.8.8.8",
				ASN:       "AS15169",
				City:      "Mountain View",
				Country:   "US",
				Location:  "37.386,-122.0838",
				ISP:       "Google LLC",
			},
		},
	}

	argusIpClient, err := iputil.NewArgusIPClient(&iputil.MockIPInfoClient{})
	assert.NoError(t, err)
	assert.NotNil(t, argusIpClient)

	gh := NewGinHandler(config.Config{}, getTestDatabase(ctx, t), argusIpClient)

	router := gin.Default()
	router.POST("/agents", gh.HandleCreateAgent)

	body, _ := json.Marshal(testData.requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/agents", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, testData.expectedCode, w.Code)

	response, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	var createAgentResponse CreateAgentResponse
	err = json.Unmarshal(response, &createAgentResponse)
	assert.NoError(t, err)
	assert.Equal(t, testData.expectedResponse.Agent.ASN, createAgentResponse.Agent.ASN)
	assert.Equal(t, testData.expectedResponse.Agent.ISP, createAgentResponse.Agent.ISP)
	assert.Equal(t, testData.expectedResponse.Agent.IPAddress, createAgentResponse.Agent.IPAddress)
	assert.Equal(t, testData.expectedResponse.Agent.Location, createAgentResponse.Agent.Location)
	assert.Equal(t, testData.expectedResponse.Agent.Country, createAgentResponse.Agent.Country)
	assert.Equal(t, testData.expectedResponse.Agent.City, createAgentResponse.Agent.City)
}

func TestHandleCreateAgent_EmptyRequestBody(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)

	testData := struct {
		requestBody      CreateAgentRequest
		expectedCode     int
		expectedResponse ErrorResponse
	}{
		requestBody:  CreateAgentRequest{},
		expectedCode: http.StatusBadRequest,
		expectedResponse: ErrorResponse{
			Error: "ip_address: IP address cannot be empty.",
		},
	}

	argusIpClient, err := iputil.NewArgusIPClient(&iputil.MockIPInfoClient{})
	assert.NoError(t, err)
	assert.NotNil(t, argusIpClient)

	gh := NewGinHandler(config.Config{}, getTestDatabase(ctx, t), argusIpClient)

	router := gin.Default()
	router.POST("/agents", gh.HandleCreateAgent)

	body, _ := json.Marshal(testData.requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/agents", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, testData.expectedCode, w.Code)

	response, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	var errorResponse ErrorResponse
	err = json.Unmarshal(response, &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, testData.expectedResponse.Error, errorResponse.Error)
}

func TestHandleCreateAgent_InvalidRequestBody(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)

	testData := struct {
		requestBody      CreateAgentRequest
		expectedCode     int
		expectedResponse ErrorResponse
	}{
		requestBody: CreateAgentRequest{
			IPAddress: ".1.1",
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: ErrorResponse{
			Error: "ip_address: IP should be in format of IPv4.",
		},
	}

	argusIpClient, err := iputil.NewArgusIPClient(&iputil.MockIPInfoClient{})
	assert.NoError(t, err)
	assert.NotNil(t, argusIpClient)

	gh := NewGinHandler(config.Config{}, getTestDatabase(ctx, t), argusIpClient)

	router := gin.Default()
	router.POST("/agents", gh.HandleCreateAgent)

	body, _ := json.Marshal(testData.requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/agents", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, testData.expectedCode, w.Code)

	response, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	var errorResponse ErrorResponse
	err = json.Unmarshal(response, &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, testData.expectedResponse.Error, errorResponse.Error)
}

func TestHandleCreateAgent_IPInfoError(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)

	testData := struct {
		requestBody      CreateAgentRequest
		expectedCode     int
		expectedResponse ErrorResponse
	}{
		requestBody: CreateAgentRequest{
			IPAddress: "1.1.1.1",
		},
		expectedCode: http.StatusServiceUnavailable,
		expectedResponse: ErrorResponse{
			Error: "cannot gather statistics for this IP address",
		},
	}

	argusIpClient, err := iputil.NewArgusIPClient(&iputil.MockIPInfoClientWithError{})
	assert.NoError(t, err)
	assert.NotNil(t, argusIpClient)

	gh := NewGinHandler(config.Config{}, getTestDatabase(ctx, t), argusIpClient)

	router := gin.Default()
	router.POST("/agents", gh.HandleCreateAgent)

	body, _ := json.Marshal(testData.requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/agents", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, testData.expectedCode, w.Code)

	response, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	var errorResponse ErrorResponse
	err = json.Unmarshal(response, &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, testData.expectedResponse.Error, errorResponse.Error)
}

func TestHandleCreateAgent_IPInfoTimeout(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)

	testData := struct {
		requestBody      CreateAgentRequest
		expectedCode     int
		expectedResponse ErrorResponse
	}{
		requestBody: CreateAgentRequest{
			IPAddress: "1.1.1.1",
		},
		expectedCode: http.StatusServiceUnavailable,
		expectedResponse: ErrorResponse{
			Error: "timeout exceeded",
		},
	}

	argusIpClient, err := iputil.NewArgusIPClient(&iputil.MockIPInfoClientWithTimeout{})
	assert.NoError(t, err)
	assert.NotNil(t, argusIpClient)

	gh := NewGinHandler(config.Config{}, getTestDatabase(ctx, t), argusIpClient)

	router := gin.Default()
	router.POST("/agents", gh.HandleCreateAgent)

	body, _ := json.Marshal(testData.requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/agents", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, testData.expectedCode, w.Code)

	response, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	var errorResponse ErrorResponse
	err = json.Unmarshal(response, &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, testData.expectedResponse.Error, errorResponse.Error)
}

func TestHandleGetAgents_Success(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)

	testData := struct {
		queryParams      string
		expectedCode     int
		expectedResponse GetAgentsResponse
	}{
		queryParams:  "?page=1&page_size=1",
		expectedCode: http.StatusOK,
		expectedResponse: GetAgentsResponse{
			Message: "retrieved agents successfully",
			Data: AgentsData{
				Agents: []Agent{
					{
						IPAddress: "8.8.8.8",
						ASN:       "AS15169",
						City:      "Mountain View",
						Country:   "US",
						Location:  "37.386,-122.0838",
					},
				},
				Pagination: AgentPagination{
					TotalAgents: 1,
					TotalPages:  1,
					CurrentPage: 1,
					PerPage:     1,
				},
			},
		},
	}

	argusIpClient, err := iputil.NewArgusIPClient(&iputil.MockIPInfoClient{})
	assert.NoError(t, err)
	assert.NotNil(t, argusIpClient)

	gh := NewGinHandler(config.Config{}, getTestDatabase(ctx, t), argusIpClient)

	router := gin.Default()
	router.GET("/agents", gh.HandleGetAgents)

	req, _ := http.NewRequest(http.MethodGet, "/agents"+testData.queryParams, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, testData.expectedCode, w.Code)

	response, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	var getAgentsResponse GetAgentsResponse
	err = json.Unmarshal(response, &getAgentsResponse)
	assert.NoError(t, err)
	assert.Equal(t, testData.expectedResponse.Message, getAgentsResponse.Message)
	assert.Equal(t, int(getAgentsResponse.Data.Pagination.TotalAgents), len(getAgentsResponse.Data.Agents))
	assert.Equal(t, len(testData.expectedResponse.Data.Agents), len(getAgentsResponse.Data.Agents))
	firstExpectedAgent := testData.expectedResponse.Data.Agents[0]
	firstActualAgent := getAgentsResponse.Data.Agents[0]
	assert.Equal(t, firstExpectedAgent.ASN, firstActualAgent.ASN)
	assert.Equal(t, firstExpectedAgent.IPAddress, firstActualAgent.IPAddress)
}

func TestHandleGetAgents_NotFound(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)

	testData := struct {
		queryParams      string
		expectedCode     int
		expectedResponse GetAgentsResponse
	}{
		queryParams:  "?page=10&page_size=1",
		expectedCode: http.StatusNotFound,
		expectedResponse: GetAgentsResponse{
			Message: "there is no agents for this page",
			Data: AgentsData{
				Agents: []Agent{},
				Pagination: AgentPagination{
					TotalAgents: 1,
					TotalPages:  1,
					CurrentPage: 10,
					PerPage:     1,
				},
			},
		},
	}

	argusIpClient, err := iputil.NewArgusIPClient(&iputil.MockIPInfoClient{})
	assert.NoError(t, err)
	assert.NotNil(t, argusIpClient)

	gh := NewGinHandler(config.Config{}, getTestDatabase(ctx, t), argusIpClient)

	router := gin.Default()
	router.GET("/agents", gh.HandleGetAgents)

	req, _ := http.NewRequest(http.MethodGet, "/agents"+testData.queryParams, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, testData.expectedCode, w.Code)

	response, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	var getAgentsResponse GetAgentsResponse
	err = json.Unmarshal(response, &getAgentsResponse)
	assert.NoError(t, err)
	assert.Equal(t, testData.expectedResponse.Message, getAgentsResponse.Message)
	assert.Equal(t, len(testData.expectedResponse.Data.Agents), len(getAgentsResponse.Data.Agents))
}

func TestHandleGetAgentDetails_Success(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)

	// Create a test data
	testDB := getTestDatabase(ctx, t)
	createdAgent, err := testDB.CreateNewAgent(ctx, &db.Agent{
		IPAddress: "8.8.8.8",
		ASN:       "AS15169",
		City:      "Mountain View",
		Country:   "US",
		Location:  "37.386,-122.0838",
		ISP:       "Google LLC",
	})
	assert.NoError(t, err)
	assert.NotNil(t, createdAgent)

	testData := struct {
		id               uint
		expectedCode     int
		expectedResponse AgentDetailedResponse
	}{
		id:           createdAgent.ID,
		expectedCode: http.StatusOK,
		expectedResponse: AgentDetailedResponse{
			Message: "agent has been retrieved successfully",
			Agent: Agent{
				ID:        createdAgent.ID,
				IPAddress: "8.8.8.8",
				ASN:       "AS15169",
				City:      "Mountain View",
				Country:   "US",
				Location:  "37.386,-122.0838",
				ISP:       "Google LLC",
			},
		},
	}

	argusIpClient, err := iputil.NewArgusIPClient(&iputil.MockIPInfoClient{})
	assert.NoError(t, err)
	assert.NotNil(t, argusIpClient)

	gh := NewGinHandler(config.Config{}, testDB, argusIpClient)

	router := gin.Default()
	router.GET("/agents/:agent_id", gh.HandleGetAgentDetail)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/agents/%d", testData.id), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, testData.expectedCode, w.Code)

	response, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	var agentDetailedResponse AgentDetailedResponse
	err = json.Unmarshal(response, &agentDetailedResponse)
	assert.NoError(t, err)
	assert.Equal(t, testData.expectedResponse.Message, agentDetailedResponse.Message)
	assert.Equal(t, testData.expectedResponse.Agent.ID, agentDetailedResponse.Agent.ID)
	assert.Equal(t, testData.expectedResponse.Agent.ASN, agentDetailedResponse.Agent.ASN)
	assert.Equal(t, testData.expectedResponse.Agent.ISP, agentDetailedResponse.Agent.ISP)
	assert.Equal(t, testData.expectedResponse.Agent.IPAddress, agentDetailedResponse.Agent.IPAddress)
	assert.Equal(t, testData.expectedResponse.Agent.Location, agentDetailedResponse.Agent.Location)
	assert.Equal(t, testData.expectedResponse.Agent.Country, agentDetailedResponse.Agent.Country)
	assert.Equal(t, testData.expectedResponse.Agent.City, agentDetailedResponse.Agent.City)
}

func TestHandleGetAgentDetails_NotFound(t *testing.T) {
	ctx := context.Background()
	gin.SetMode(gin.TestMode)

	testData := struct {
		id               uint
		expectedCode     int
		expectedResponse ErrorResponse
	}{
		id:           10000000,
		expectedCode: http.StatusNotFound,
		expectedResponse: ErrorResponse{
			Error: "cannot find such agent by id",
		},
	}

	argusIpClient, err := iputil.NewArgusIPClient(&iputil.MockIPInfoClient{})
	assert.NoError(t, err)
	assert.NotNil(t, argusIpClient)

	gh := NewGinHandler(config.Config{}, getTestDatabase(ctx, t), argusIpClient)

	router := gin.Default()
	router.GET("/agents/:agent_id", gh.HandleGetAgentDetail)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/agents/%d", testData.id), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, testData.expectedCode, w.Code)

	response, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	var errorResposne ErrorResponse
	err = json.Unmarshal(response, &errorResposne)
	assert.NoError(t, err)
	assert.Equal(t, testData.expectedResponse.Error, errorResposne.Error)
}
