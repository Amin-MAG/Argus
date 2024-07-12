package handlers

// PingResponse represents the response format for the Ping API.
type PingResponse struct {
	DatabaseStatus string `json:"database_status"`
}
