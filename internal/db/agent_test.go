package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateNewAgent(t *testing.T) {
	ctx := context.Background()
	tdb := getTestDatabase(ctx, t)

	// Create a new agent
	newAgent := &Agent{
		IPAddress: "192.168.1.100",
		ASN:       "AS12345",
		ISP:       "Test ISP",
	}

	createdAgent, err := tdb.(*GormDB).CreateNewAgent(ctx, newAgent)
	assert.NoError(t, err, "error creating new agent")
	assert.NotNil(t, createdAgent, "created agent should not be nil")
	assert.NotZero(t, createdAgent.ID, "agent ID should not be zero")
}

func TestGetAllAgents(t *testing.T) {
	ctx := context.Background()
	tdb := getTestDatabase(ctx, t)

	// Prepare test data
	filter := &AgentFilter{
		IPAddress: nil,
	}
	page := 1
	pageSize := 10
	sort := &AgentSort{
		SortBy:  nil,
		OrderBy: nil,
	}

	// Call GetAllAgents
	result, err := tdb.(*GormDB).GetAllAgents(ctx, filter, page, pageSize, sort)
	assert.NoError(t, err, "error fetching agents")
	assert.NotNil(t, result, "result should not be nil")
	assert.NotZero(t, result.TotalAgents, "total agents count should not be zero")
}

func TestGetAgentByID(t *testing.T) {
	ctx := context.Background()
	tdb := getTestDatabase(ctx, t)

	// Create a new agent for testing
	newAgent := &Agent{
		IPAddress: "192.168.1.101",
		ASN:       "AS54321",
		ISP:       "Test ISP 2",
	}
	createdAgent, err := tdb.(*GormDB).CreateNewAgent(ctx, newAgent)
	assert.NoError(t, err, "error creating new agent")

	// Fetch the agent by ID
	fetchedAgent, err := tdb.(*GormDB).GetAgentByID(ctx, createdAgent.ID)
	assert.NoError(t, err, "error fetching agent by ID")
	assert.NotNil(t, fetchedAgent, "fetched agent should not be nil")
	assert.Equal(t, createdAgent.ID, fetchedAgent.ID, "fetched agent ID should match created agent ID")
	assert.Equal(t, newAgent.IPAddress, fetchedAgent.IPAddress, "fetched agent IP address should match")
}
