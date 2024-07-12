package handlers

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

// Agent represents information about an agent.
type Agent struct {
	ID        uint      `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	IPAddress string    `json:"ip_address,omitempty"`
	ASN       string    `json:"asn,omitempty"`
	ISP       string    `json:"isp,omitempty"`
	City      string    `json:"city,omitempty"`
	Country   string    `json:"country,omitempty"`
	Location  string    `json:"location,omitempty"`
}

// CreateAgentRequest represents the request format for creating a new agent.
type CreateAgentRequest struct {
	IPAddress string `json:"ip_address"`
}

// CreateAgentResponse represents the response format for creating a new agent.
func (req CreateAgentRequest) validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.IPAddress,
			validation.Required.Error("IP address cannot be empty"),
			is.IPv4.Error("IP should be in format of IPv4"),
		),
	)
}

// CreateAgentResponse represents the response format for creating a new agent.
type CreateAgentResponse struct {
	Message string `json:"message"`
	Agent   Agent  `json:"agent"`
}

// GetAgentsQueryParams represents the query parameters for fetching agents.
type GetAgentsQueryParams struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	IPAddress string `form:"ip_address"`
	SortBy    string `form:"sort_by"`
	Order     string `form:"order"`
}

// AgentPagination represents pagination details for a list of agents.
type AgentPagination struct {
	TotalAgents int64 `json:"total_agents"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
}

// AgentsData represents data containing a list of agents and pagination details.
type AgentsData struct {
	Agents     []Agent         `json:"agents"`
	Pagination AgentPagination `json:"pagination"`
}

// GetAgentsResponse represents the response format for fetching agents.
type GetAgentsResponse struct {
	Message string     `json:"message"`
	Data    AgentsData `json:"data"`
}

// AgentDetailedResponse represents the response format for fetching detailed agent information.
type AgentDetailedResponse struct {
	Message string `json:"message,omitempty"`
	Agent   Agent  `json:"agent,omitempty"`
}
