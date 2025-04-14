package utils

import (
    "crypto/rand"
    "encoding/hex"
    "log"
)

// GenerateRandomID generates a 24-character hexadecimal string similar to MongoDB ObjectID
func GenerateRandomID() string {
    bytes := make([]byte, 12)
    _, err := rand.Read(bytes)
    if err != nil {
        log.Fatalf("failed to generate random ID: %v", err)
    }

    return hex.EncodeToString(bytes)
}