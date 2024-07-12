package iputil

import (
	"context"
	"github.com/ipinfo/go/v2/ipinfo"
	"net"
)

type IPInfoClient interface {
	GetIPInfo(ip net.IP) (*ipinfo.Core, error)
}

type ArgusIPInfoClient struct {
	client IPInfoClient
}

func NewArgusIPClient(client IPInfoClient) (IPStatsGatherer, error) {
	return &ArgusIPInfoClient{client: client}, nil
}

func (ipi *ArgusIPInfoClient) GetInfo(ctx context.Context, ip string) (*Stats, error) {
	info, err := ipi.client.GetIPInfo(net.ParseIP(ip))
	if err != nil {
		return nil, err
	}

	return &Stats{
		IP:          info.IP,
		City:        info.City,
		Region:      info.Region,
		Country:     info.Country,
		CountryName: info.CountryName,
		Location:    info.Location,
		ISP:         info.ASN.Name,
		ASN:         info.ASN.ASN,
	}, nil
}
