package service

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateVerifyToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
