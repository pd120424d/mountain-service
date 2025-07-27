package server

import (
	"os"
	"strings"
)

// ReadSecret reads a secret from a file path and returns it trimmed of whitespace
func ReadSecret(filePath string) (string, error) {
	secret, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(secret)), nil
}
