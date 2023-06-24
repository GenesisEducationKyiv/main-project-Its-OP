package validators

import (
	"btcRate/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate_UnsupportedCurrency(t *testing.T) {
	// Arrange
	validator := NewSupportedCurrencyValidator([]string{"USD"})

	// Act
	err := validator.Validate("GBP")

	// Assert
	assert.NotNil(t, err)
	assert.IsType(t, &domain.ArgumentError{}, err)
	assert.Equal(t, "currency GBP is not supported", err.(*domain.ArgumentError).Message)
}
