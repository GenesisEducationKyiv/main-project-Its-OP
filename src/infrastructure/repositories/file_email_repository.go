package repositories

import (
	"btcRate/domain"
	"encoding/json"
	"fmt"
	"os"
)

type FileEmailRepository struct {
	emails      map[string]struct{}
	storageFile string
}

func NewFileEmailRepository(storageFile string) (*FileEmailRepository, error) {
	emails := map[string]struct{}{}

	if fileExists(storageFile) {
		data, err := os.ReadFile(storageFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, &emails)
		if err != nil {
			return nil, err
		}
	}

	r := FileEmailRepository{emails: emails, storageFile: storageFile}
	return &r, nil
}

func (r *FileEmailRepository) AddEmail(email string) error {
	if _, exists := r.emails[email]; exists {
		return &domain.DataConsistencyError{Message: fmt.Sprintf("Email address '%s' is already present in the database", email)}
	}

	r.emails[email] = struct{}{}
	return nil
}

func (r *FileEmailRepository) GetAll() []string {
	if !fileExists(r.storageFile) {
		return []string{}
	}

	var emails []string
	for email := range r.emails {
		emails = append(emails, email)
	}

	return emails
}

func (r *FileEmailRepository) Save() error {
	data, err := json.Marshal(r.emails)
	if err != nil {
		return err
	}

	if !fileExists(r.storageFile) {
		err = createFile(r.storageFile)
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
