package handlers

import (
	"argus/config"
	"argus/internal/iputil"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
