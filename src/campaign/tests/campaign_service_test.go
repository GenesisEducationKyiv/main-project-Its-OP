package tests

import (
	"btcRate/campaign/application"
	"btcRate/campaign/application/validators"
	"btcRate/campaign/domain"
	"btcRate/campaign/infrastructure/repositories"
	"btcRate/common/application/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"sync"
	"testing"
)

const storageFile = "artifacts/emails.json"

func setup(t *testing.T) *application.CampaignService {
	mutex := &sync.RWMutex{}
	emailRepo, _ := repositories.NewFileEmailRepository(storageFile, mutex)
	emailValidator := &validators.EmailValidator{}

	logger := mocks.NewILogger(t)
	logger.EXPECT().Info(mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return()

	service := application.NewCampaignService(emailRepo, nil, nil, emailValidator, logger)

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
	service := setup(t)

	// Act
	err := service.Subscribe("test@example.com")

	// Assert
	assert.Nil(t, err)
}

func TestSubscribe_Duplicate(t *testing.T) {
	// Arrange
	defer teardown(t)
	service := setup(t)
	err := service.Subscribe("test@example.com")
	assert.Nil(t, err)

	// Act
	err = service.Subscribe("test@example.com")

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, "Email address 'test@example.com' is already present in the database", err.(*domain.DataConsistencyError).Message)
}
