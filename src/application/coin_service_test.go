package application

import (
	"btcRate/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setup(t *testing.T) (*CoinService, *MockICoinClient) {
	coinClient := NewMockICoinClient(t)
	supportedCurrencies := []string{"UAH", "USD"}
	supportedCoins := []string{"BTC", "ETH"}

	service := NewCoinService(supportedCurrencies, supportedCoins, coinClient, nil, nil)

	return service, coinClient
}

func TestGetCurrentRate_UnsupportedCurrency(t *testing.T) {
	// Arrange
	service, _ := setup(t)

	// Act
	price, err := service.GetCurrentRate("GBP", "BTC")

	// Assert
	assert.Nil(t, price)
	assert.NotNil(t, err)
	assert.Equal(t, "Currency GBP is not supported", err.(domain.ArgumentError).Message)
}

func TestGetCurrentRate_UnsupportedCoin(t *testing.T) {
	// Arrange
	service, _ := setup(t)

	// Act
	price, err := service.GetCurrentRate("UAH", "DOGE")

	// Assert
	assert.Nil(t, price)
	assert.NotNil(t, err)
	assert.Equal(t, "Coin DOGE is not supported", err.(domain.ArgumentError).Message)
}

func TestGetCurrentRate_Success(t *testing.T) {
	// Arrange
	service, client := setup(t)
	currency := "USD"
	coin := "BTC"
	timestamp := time.Now()
	expectedPrice := domain.Price{Amount: 31000, Currency: currency, Timestamp: timestamp}
	client.EXPECT().GetRate(currency, coin).Return(expectedPrice.Amount, timestamp, nil)

	// Act
	price, err := service.GetCurrentRate(currency, coin)

	// Assert
	assert.Equal(t, &expectedPrice, price)
	assert.Nil(t, err)
	client.AssertCalled(t, "GetRate", currency, coin)
}
