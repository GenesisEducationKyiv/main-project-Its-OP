package repositories

import (
	"btcRate/common/domain"
	"btcRate/common/infrastructure/repositories"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type FileEmailRepository struct {
	emails      map[string]struct{}
	storageFile string
	mutex       *sync.RWMutex
}

func NewFileEmailRepository(storageFile string, mutex *sync.RWMutex) (*FileEmailRepository, error) {
	emails := map[string]struct{}{}

	if repositories.FileExists(storageFile) {
		data, err := os.ReadFile(storageFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, &emails)
		if err != nil {
			return nil, err
		}
	}

	return &FileEmailRepository{emails: emails, storageFile: storageFile, mutex: mutex}, nil
}

func (r *FileEmailRepository) AddEmail(email string) error {
	if r.emailExists(email) {
		return &domain.DataConsistencyError{Message: fmt.Sprintf("Email address '%s' is already present in the database", email)}
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.emails[email] = struct{}{}
	return r.save()
}

func (r *FileEmailRepository) GetAll() []string {
	if !repositories.FileExists(r.storageFile) {
		return []string{}
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var emails []string
	for email := range r.emails {
		emails = append(emails, email)
	}

	return emails
}

func (r *FileEmailRepository) emailExists(email string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.emails[email]
	return exists
}

func (r *FileEmailRepository) save() error {
	data, err := json.Marshal(r.emails)
	if err != nil {
		return err
	}

	if !repositories.FileExists(r.storageFile) {
		err = repositories.CreateFile(r.storageFile)
		if err != nil {
			return err
		}
	}

	permissionCode := 0644

	err = os.WriteFile(r.storageFile, data, os.FileMode(permissionCode))
	if err != nil {
		return err
	}

	return nil
}
