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
	resultChan := make(chan struct {
		info *ipinfo.Core
		err  error
	}, 1)

	go func() {
		info, err := ipi.client.GetIPInfo(net.ParseIP(ip))
		resultChan <- struct {
			info *ipinfo.Core
			err  error
		}{
			info: info,
			err:  err,
		}
	}()

	select {
	case result := <-resultChan:
		if result.err != nil {
			return nil, result.err
		}
		return &Stats{
			IP:          result.info.IP,
			City:        result.info.City,
			Region:      result.info.Region,
			Country:     result.info.Country,
			CountryName: result.info.CountryName,
			Location:    result.info.Location,
			ISP:         result.info.ASN.Name,
			ASN:         result.info.ASN.ASN,
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
