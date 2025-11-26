package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	n, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	if n != 32 {
		return "", fmt.Errorf("did not read 32 bytes")
	}

	return hex.EncodeToString(randomBytes), nil
}
