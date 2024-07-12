package iputil

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert" // Using testify for assertions
)

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
