package db

import (
	"context"
	"fmt"
	"time"
)

// Agent contains data for each agent request
type Agent struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	IPAddress string `gorm:"index,not null"`
	ASN       string
	ISP       string
	City      string
	Country   string
	Location  string
}

type AgentFilter struct {
	IPAddress *string
}

type AgentSort struct {
	SortBy  *string
	OrderBy *string
}

type AgentsResult struct {
	Agents      []Agent
	TotalAgents int64
}

// CreateNewAgent creates a new agent in the database
func (gdb *GormDB) CreateNewAgent(ctx context.Context, a *Agent) (*Agent, error) {
	a.CreatedAt = time.Now()

	return a, gdb.db.Create(a).Error
}

func (gdb *GormDB) GetAllAgents(ctx context.Context, filter *AgentFilter, page int, pageSize int, sort *AgentSort) (*AgentsResult, error) {
	var agents []Agent
	var count int64

	query := gdb.db.Model(&Agent{})

	// Apply filters
	if filter != nil {
		if filter.IPAddress != nil {
			query.Where("ip_address = ?", filter.IPAddress)
		}
	}

	// Calculate the number of campaigns
	err := query.Count(&count).Error
	if err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Handle sort
	if sort != nil && sort.SortBy != nil {
		orderBy := "desc"
		if sort.OrderBy != nil {
			orderBy = *sort.OrderBy
		}
		query = query.Order(fmt.Sprintf("%s %s", *sort.SortBy, orderBy))
	}
	query.Order("agents.id")

	// Execute query
	err = query.Find(&agents).Error
	if err != nil {
		return nil, err
	}

	return &AgentsResult{
		Agents:      agents,
		TotalAgents: count,
	}, nil
}

func (gdb *GormDB) GetAgentByID(ctx context.Context, agentID uint) (*Agent, error) {
	var agent Agent
	return &agent, gdb.db.Where("id = ?", agentID).First(&agent).Error
}
