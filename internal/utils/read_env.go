package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var envLoadErr error

func init() {
	envLoadErr = godotenv.Load(".env")
	if envLoadErr != nil {
		fmt.Println("Warning: Failed to load .env file:", envLoadErr)
	}

}

func ReadEnv(lookupStr string) (string, error) {
	if envLoadErr != nil {
		return "", fmt.Errorf("failed to load .env file: %w", envLoadErr)
	}

	found := os.Getenv(lookupStr)
	if found == "" {
		return "", errors.New(lookupStr + "not found in environment variables")
	}
	return found, nil
}