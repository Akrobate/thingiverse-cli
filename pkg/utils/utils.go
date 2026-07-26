package utils

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	_, err := os.Stat(filepath.Clean(path))
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func CalculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()
	hasher := md5.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to compute file hash: %w", err)
	}
	hashBytes := hasher.Sum(nil)
	base64Hash := base64.StdEncoding.EncodeToString(hashBytes)
	return base64Hash, nil
}

func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("Cannot get files stats %s: %w", filePath, err)
	}

	return info.Size(), nil
}
