package proxies

import "coin/application/services"

type chainedCoinClientProxy struct {
	client services.ICoinClient
	next   services.ICoinClient
}

func newChainedCoinClientProxy(client services.ICoinClient, next services.ICoinClient) *chainedCoinClientProxy {
	return &chainedCoinClientProxy{client: client, next: next}
}

func (c *chainedCoinClientProxy) GetRate(currency string, coin string) (*services.SpotPrice, error) {
	price, err := c.client.GetRate(currency, coin)

	if err != nil && c.next != nil {
		return c.next.GetRate(currency, coin)
	}

	return price, err
}
