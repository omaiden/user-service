package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

const Bearer = "bearer "

func ExtractAuthToken(auth string, prefix string) string {
	tk := strings.TrimSpace(auth)
	if len(tk) <= len(prefix) {
		return ""
	}
	if !strings.EqualFold(tk[:len(prefix)], prefix) {
		return ""
	}
	tk = tk[len(prefix):]
	return tk
}

func HashToken(token string) string {
	hashedToken := sha256.Sum256([]byte(token))
	return base64.RawStdEncoding.EncodeToString(hashedToken[:])
}

func GenerateToken(size int) string {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return base64.RawStdEncoding.EncodeToString(b)
}
