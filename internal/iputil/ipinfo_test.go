package iputil

import (
	"context"
	"github.com/ipinfo/go/v2/ipinfo"
	"net"
	"testing"

	"github.com/stretchr/testify/assert" // Using testify for assertions
)

// MockIPInfoClient is a mock implementation of IPInfoClient.
type MockIPInfoClient struct{}

func (m *MockIPInfoClient) GetIPInfo(ip net.IP) (*ipinfo.Core, error) {
	core := &ipinfo.Core{
		IP:          net.ParseIP("8.8.8.8"),
		City:        "Mountain View",
		Region:      "CA",
		Country:     "US",
		CountryName: "United States",
		Location:    "37.386,-122.0838",
		ASN: &ipinfo.CoreASN{
			ASN:  "AS15169",
			Name: "Google LLC",
		},
	}
	return core, nil
}

func TestGetInfo(t *testing.T) {
	// Setup
	mockClient := &MockIPInfoClient{}
	ip := "8.8.8.8" // Example IP address

	argusClient, err := NewArgusIPClient(mockClient)
	assert.NoError(t, err)
	assert.NotNil(t, argusClient)

	// Execute
	stats, err := argusClient.GetInfo(context.Background(), ip)

	// Assert
	assert.NoError(t, err, "GetInfo should not return an error")
	assert.NotNil(t, stats, "Stats should not be nil")
	assert.Equal(t, net.ParseIP("8.8.8.8"), stats.IP, "IP should match")
	assert.Equal(t, "Mountain View", stats.City, "City should match")
}
