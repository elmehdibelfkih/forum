package auth

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// uuid replacer for the session id

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("failed to generate token %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes) // a-z A-Z /+
}
