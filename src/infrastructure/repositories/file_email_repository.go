package repositories

import (
	"btcRate/domain"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"os"
)

type FileEmailRepository struct {
	emails      hashset.Set
	storageFile string
}

func NewFileEmailRepository(storageFile string) (*FileEmailRepository, error) {
	emails := *hashset.New()

	if fileExists(storageFile) {
		data, err := os.ReadFile(storageFile)
		if err != nil {
			return nil, err
		}

		err = emails.FromJSON(data)
		if err != nil {
			return nil, err
		}
	}

	r := FileEmailRepository{emails: emails, storageFile: storageFile}
	return &r, nil
}

func (r *FileEmailRepository) AddEmail(email string) error {
	if r.emails.Contains(email) {
		return &domain.DataConsistencyError{Message: fmt.Sprintf("Email address '%s' is already present in the database", email)}
	}

	r.emails.Add(email)
	return nil
}

func (r *FileEmailRepository) GetAll() []string {
	if !fileExists(r.storageFile) {
		return []string{}
	}

	values := r.emails.Values()

	emailsArray := make([]string, len(values))
	for i, value := range values {
		emailsArray[i] = value.(string)
	}

	return emailsArray
}

func (r *FileEmailRepository) Save() error {
	data, err := r.emails.MarshalJSON()
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