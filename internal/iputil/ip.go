package iputil

import (
	"context"
	"net"
)

type IPStatsGatherer interface {
	GetInfo(ctx context.Context, ip string) (*Stats, error)
}

type Stats struct {
	IP          net.IP `json:"ip" csv:"ip"`
	City        string `json:"city,omitempty" yaml:"city,omitempty"`
	Region      string `json:"region,omitempty" yaml:"region,omitempty"`
	Country     string `json:"country,omitempty" yaml:"country,omitempty"`
	CountryName string `json:"country_name,omitempty" yaml:"countryName,omitempty"`
	Location    string `json:"loc,omitempty" yaml:"location,omitempty"`
	ISP         string `json:"isp,omitempty" yaml:"isp,omitempty"`
	ASN         string `json:"asn,omitempty" yaml:"asn,omitempty"`
}
