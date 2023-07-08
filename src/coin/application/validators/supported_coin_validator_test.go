package validators

import (
	"coin/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate_UnsupportedCoin(t *testing.T) {
	// Arrange
	validator := NewSupportedCoinValidator([]string{"BTC"})

	// Act
	err := validator.Validate("DOGE")

	// Assert
	assert.NotNil(t, err)
	assert.IsType(t, &domain.ArgumentError{}, err)
	assert.Equal(t, "coin DOGE is not supported", err.(*domain.ArgumentError).Message)
}
