package repositories

import (
	"os"
	"path/filepath"
)

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
