package tests

import (
	"btcRate/application"
	"btcRate/application/validators"
	"btcRate/domain"
	"btcRate/infrastructure"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const storageFile = "artifacts/emails.json"

func setup() *application.CampaignService {
	emailRepo, _ := infrastructure.NewFileEmailRepository(storageFile)
	emailValidator := &validators.EmailValidator{}
	service := application.NewCampaignService(emailRepo, nil, emailValidator)

	return service
}

func teardown(t *testing.T) {
	err := os.Remove(storageFile)
	if err != nil {
		t.Fatal("failed to delete file")
	}
}

func TestSubscribe_Success(t *testing.T) {
	// Arrange
	defer teardown(t)
	service := setup()

	// Act
	err := service.Subscribe("test@example.com")

	// Assert
	assert.Nil(t, err)
}

func TestSubscribe_Duplicate(t *testing.T) {
	// Arrange
	defer teardown(t)
	service := setup()
	err := service.Subscribe("test@example.com")
	assert.Nil(t, err)

	// Act
	err = service.Subscribe("test@example.com")

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, "Email address 'test@example.com' is already present in the database", err.(*domain.DataConsistencyError).Message)
}
