package db

import "context"

type DB interface {
	Ping(ctx context.Context) error

	CreateNewAgent(ctx context.Context, agent *Agent) (*Agent, error)
	GetAllAgents(ctx context.Context, filter *AgentFilter, page int, pageSize int, sort *AgentSort) (*AgentsResult, error)
	GetAgentByID(ctx context.Context, agentID uint) (*Agent, error)
}
