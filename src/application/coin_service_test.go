package application

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type TestCoinClient struct {
	mock.Mock
}

func (m *TestCoinClient) GetRate(currency string, coin string) (float64, time.Time, error) {
	args := m.Called(currency, coin)
	return args.Get(0).(float64), args.Get(1).(time.Time), args.Error(2)
}

func setup() (*CoinService, *TestCoinClient) {
	coinClient := &TestCoinClient{}

	service := NewCoinService(coinClient, nil, nil, nil)

	return service, coinClient
}

func TestGetCurrentRate_Success(t *testing.T) {
	// Arrange
	service, client := setup()
	client.On("GetRate", "USD", "BTC").Return(31000.0, time.Now(), nil)

	// Act
	_, err := service.GetCurrentRate("USD", "BTC")

	// Assert
	assert.Nil(t, err)
	client.AssertCalled(t, "GetRate", "USD", "BTC")
}
