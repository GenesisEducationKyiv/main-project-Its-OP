package infrastructure

import (
	"btcRate/domain"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"os"
	"path/filepath"
)

const storageFile = "./data/emails.json"

type FileEmailRepository struct {
	Emails hashset.Set
}

func NewFileEmailRepository() (*FileEmailRepository, error) {
	emails := *hashset.New()

	if fileExists() {
		data, err := os.ReadFile(storageFile)
		if err != nil {
			dbError := &domain.InternalError{NestedError: err}
			return nil, dbError
		}

		err = emails.FromJSON(data)
		if err != nil {
			dbError := &domain.InternalError{NestedError: err}
			return nil, dbError
		}
	}

	r := FileEmailRepository{Emails: emails}
	return &r, nil
}

func (r *FileEmailRepository) AddEmail(email string) error {
	if r.Emails.Contains(email) {
		return &domain.DataConsistencyError{Message: fmt.Sprintf("Email address '%s' is already present in the database", email)}
	}

	r.Emails.Add(email)
	return nil
}

func (r *FileEmailRepository) GetAll() []string {
	if !fileExists() {
		return []string{}
	}

	values := r.Emails.Values()

	emailsArray := make([]string, len(values))
	for i, value := range values {
		emailsArray[i] = value.(string)
	}

	return emailsArray
}

func (r *FileEmailRepository) Save() error {
	data, err := r.Emails.MarshalJSON()
	if err != nil {
		dbError := &domain.InternalError{NestedError: err}
		return dbError
	}

	if !fileExists() {
		err = createFile(storageFile)
		if err != nil {
			dbError := &domain.InternalError{NestedError: err}
			return dbError
		}
	}

	permissionCode := 0644

	err = os.WriteFile(storageFile, data, os.FileMode(permissionCode))
	if err != nil {
		dbError := &domain.InternalError{NestedError: err}
		return dbError
	}

	return nil
}

func fileExists() bool {
	info, err := os.Stat(storageFile)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func createFile(filePath string) error {
	dirPath := filepath.Dir(filePath)

	permissionCode := 0755

	err := os.MkdirAll(dirPath, os.FileMode(permissionCode))

	if err != nil {
		dbError := &domain.InternalError{NestedError: err}
		return dbError
	}

	file, err := os.Create(filePath)
	if err != nil {
		dbError := &domain.InternalError{NestedError: err}
		return dbError
	}

	defer file.Close()
	return nil
}
