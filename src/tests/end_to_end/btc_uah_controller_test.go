package end_to_end

import (
	"btcRate/coin/web"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func setup(t *testing.T) (web.ServerManager, func()) {
	server := web.NewServerManager()
	stop, err := server.RunServer("./logs/coin-logs.csv")
	if err != nil {
		t.Fatal("unable to start the server")
	}
	time.Sleep(2 * time.Second)

	return server, func() {
		if err := stop(); err != nil {
			t.Fatal("unable to stop the server")
		}
	}
}

func TestRateApi(t *testing.T) {
	// Arrange
	server, stop := setup(t)
	defer stop()

	// Act
	resp, err := server.GetRate("http://localhost:8080")

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.True(t, resp.Successful)
	assert.True(t, resp.Body.Amount > 0)
	assert.True(t, resp.Body.Currency == "UAH")
}
