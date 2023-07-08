package repositories

import (
	"fmt"
	"os"
)

type FileLogRepository struct {
	storageFile string
}

func NewFileLogRepository(storageFile string) (*FileLogRepository, error) {
	if !FileExists(storageFile) {
		err := CreateFile(storageFile)
		if err != nil {
			return nil, err
		}
	}
	return &FileLogRepository{storageFile: storageFile}, nil
}

func (r *FileLogRepository) Log(data string) error {
	f, err := os.OpenFile(r.storageFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("error: failed to close file %s, %f", r.storageFile, err)
		}
	}()

	if _, err := f.WriteString(fmt.Sprintf("%s\n", data)); err != nil {
		return err
	}

	return nil
}
