package handlers

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type AgentResponse struct {
	ID        uint   `json:"id"`
	IPAddress string `json:"ip_address"`
}

type AgentDetailedResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	IPAddress string    `json:"ip_address"`
	ASN       string    `json:"asn"`
	ISP       string    `json:"isp"`
}

type CreateAgentRequest struct {
	IPAddress string `json:"ip_address"`
}

func (req CreateAgentRequest) validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.IPAddress,
			validation.Required.Error("IP address cannot be empty"),
			is.IPv4.Error("IP should be in format of IPv4"),
		),
	)
}

type CreateAgentResponse struct {
	Message string                `json:"message"`
	Agent   AgentDetailedResponse `json:"data"`
}

type GetAgentsQueryParams struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	IPAddress string `form:"ip_address"`
	SortBy    string `form:"sort_by"`
	Order     string `form:"order"`
}

type AgentPagination struct {
	TotalAgents int64 `json:"total_agents"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
}

type AgentsData struct {
	Agents     []AgentResponse `json:"agents"`
	Pagination AgentPagination `json:"pagination"`
}

type GetAgentsResponse struct {
	Message string     `json:"message"`
	Data    AgentsData `json:"data"`
}
