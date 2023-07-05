package application

import "btcRate/application/services"

type ChainedCoinClientFactory struct {
	factories []services.ICoinClientFactory
}

func NewChainedCoinClientFactory(cfs []services.ICoinClientFactory) *ChainedCoinClientFactory {
	return &ChainedCoinClientFactory{factories: cfs}
}

func (f *ChainedCoinClientFactory) CreateClient() services.ICoinClient {
	if len(f.factories) == 0 {
		return nil
	} else if len(f.factories) == 1 {
		return f.factories[0].CreateClient()
	}

	lastChainedClient := newChainedCoinClientProxy(f.factories[len(f.factories)-1].CreateClient(), nil)
	for i := len(f.factories) - 2; i >= 0; i-- {
		chainedClient := newChainedCoinClientProxy(f.factories[i].CreateClient(), lastChainedClient)
		lastChainedClient = chainedClient
	}

	return lastChainedClient
}
