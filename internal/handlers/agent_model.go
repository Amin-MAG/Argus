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
