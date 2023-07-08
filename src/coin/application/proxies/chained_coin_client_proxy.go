package proxies

import (
	"coin/application"
)

type chainedCoinClientProxy struct {
	client application.ICoinClient
	next   application.ICoinClient
}

func newChainedCoinClientProxy(client application.ICoinClient, next application.ICoinClient) *chainedCoinClientProxy {
	return &chainedCoinClientProxy{client: client, next: next}
}

func (c *chainedCoinClientProxy) GetRate(currency string, coin string) (*application.SpotPrice, error) {
	price, err := c.client.GetRate(currency, coin)

	if err != nil && c.next != nil {
		return c.next.GetRate(currency, coin)
	}

	return price, err
}
