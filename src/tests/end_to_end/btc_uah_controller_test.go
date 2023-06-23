package end_to_end

import (
	"btcRate/web"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestRateApi(t *testing.T) {
	// Arrange
	stop, err := web.RunBtcUahController("./data/emails.json")
	if err != nil {
		t.Fatal("unable to start the server")
	}
	time.Sleep(2 * time.Second)

	defer func() {
		if err := stop(); err != nil {
			t.Fatal("unable to stop the server")
		}
	}()

	// Act
	resp, err := http.Get("http://localhost:8080/api/v1/rate")
	if err != nil {
		t.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error("failed to close response body")
		}
	}(resp.Body)

	var price int
	err = json.NewDecoder(resp.Body).Decode(&price)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, price)
}
