package repositories

import (
	"fmt"
	"os"
)

type FileLogRepository struct {
	storageFile string
}

func NewFileLogRepository(storageFile string) *FileLogRepository {
	return &FileLogRepository{storageFile: storageFile}
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
