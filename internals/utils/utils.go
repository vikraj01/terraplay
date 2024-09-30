package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateRandomID(length int) string {
	byteLength := (length + 1) / 2
	randomBytes := make([]byte, byteLength)

	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(randomBytes)[:length]
}

func GenerateTimestampID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func GenerateCryptoID(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func GenerateUUID() string {
	id := uuid.New()
	return id.String()
}


