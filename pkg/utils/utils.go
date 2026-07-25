package utils

import (
	"crypto/md5"
	"encoding/hex"
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
		return "", fmt.Errorf("Cannot open file %s: %w", filePath, err)
	}
	defer file.Close()
	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("Hash calculation error: %w", err)
	}
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}

func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("Cannot get files stats %s: %w", filePath, err)
	}

	return info.Size(), nil
}
