package infrastructure

import (
	"btcRate/domain"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"os"
	"path/filepath"
)

type FileEmailRepository struct {
	Emails      hashset.Set
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

	r := FileEmailRepository{Emails: emails, storageFile: storageFile}
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
	if !fileExists(r.storageFile) {
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

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
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
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()
	return nil
}
