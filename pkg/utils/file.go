package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func ReadFile(filepath string) []byte {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil
	}
	return data
}

func WriteFile(filePath string, content []byte) bool {
	err := os.WriteFile(filePath, content, 0644)
	return err == nil
}

func FindFile(rootDir, fileName string) (filePath string) {
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil // Skip if directory
		}
		if strings.ToLower(info.Name()) == strings.ToLower(fileName) {
			filePath = path
		}
		return nil
	})
	if err != nil {
		return ""
	}

	return filePath
}

func WriteOutput(filePath string, content []byte) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		panic(err)
	}
}
