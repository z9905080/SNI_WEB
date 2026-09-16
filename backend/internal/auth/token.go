package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// NewToken 產生放在 cookie 的原始 token 與存入資料庫的 SHA-256 雜湊。
func NewToken() (raw, hash string) {
	b := make([]byte, 32)
	rand.Read(b)
	raw = base64.RawURLEncoding.EncodeToString(b)
	return raw, HashToken(raw)
}

func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
