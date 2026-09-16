package auth

import (
	"encoding/base64"
	"regexp"
	"testing"
)

func TestNewToken(t *testing.T) {
	raw, hash := NewToken()
	b, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil || len(b) != 32 {
		t.Fatalf("raw 應為 32 bytes base64url：%q (%v)", raw, err)
	}
	if !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(hash) {
		t.Fatalf("hash 應為 sha256 hex：%q", hash)
	}
	if HashToken(raw) != hash {
		t.Fatal("HashToken 應與 NewToken 回傳的 hash 一致")
	}
	raw2, _ := NewToken()
	if raw == raw2 {
		t.Fatal("token 不應重複")
	}
}
