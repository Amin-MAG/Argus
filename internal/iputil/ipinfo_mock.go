package iputil

import (
	"errors"
	"github.com/ipinfo/go/v2/ipinfo"
	"net"
	"time"
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

// MockIPInfoClientWithError is a mock implementation of IPInfoClient that returns error.
type MockIPInfoClientWithError struct{}

func (m *MockIPInfoClientWithError) GetIPInfo(ip net.IP) (*ipinfo.Core, error) {
	return nil, errors.New("cannot get ip info")
}

// MockIPInfoClientWithTimeout is a mock implementation of IPInfoClient that has timeout.
type MockIPInfoClientWithTimeout struct{}

func (m *MockIPInfoClientWithTimeout) GetIPInfo(ip net.IP) (*ipinfo.Core, error) {
	time.Sleep(time.Minute)
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
