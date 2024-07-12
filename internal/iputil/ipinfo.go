package iputil

import (
	"context"
	"fmt"
	"github.com/ipinfo/go/v2/ipinfo"
	"net"
)

type IPInfoClient struct {
	client *ipinfo.Client
}

func NewIPInfoClient(token string) (IPStatsGatherer, error) {
	client := ipinfo.NewClient(nil, nil, token)
	return &IPInfoClient{client: client}, nil
}

func (ipi *IPInfoClient) GetInfo(ctx context.Context, ip string) (*Stats, error) {
	info, err := ipi.client.GetIPInfo(net.ParseIP(ip))
	if err != nil {
		return nil, err
	}

	fmt.Println(info.Country)
	fmt.Println(info.City)
	fmt.Println(info.Location)

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
