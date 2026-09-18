package auth

import (
	"errors"
	"strings"
	"testing"
)

func TestPassword(t *testing.T) {
	h, err := HashPassword("correct horse")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(h, "$2") {
		t.Fatalf("不是 bcrypt：%s", h)
	}
	if !CheckPassword(h, "correct horse") {
		t.Fatal("正確密碼應通過")
	}
	if CheckPassword(h, "wrong") {
		t.Fatal("錯誤密碼不應通過")
	}
	if CheckPassword("21c2e59531c8710156d34a3c30ac81d5", "x") {
		t.Fatal("非 bcrypt 的舊資料不應通過")
	}
}

func TestPasswordLength(t *testing.T) {
	if _, err := HashPassword("1234567"); !errors.Is(err, ErrPasswordTooShort) {
		t.Fatalf("got %v", err)
	}
	if _, err := HashPassword(strings.Repeat("a", 73)); !errors.Is(err, ErrPasswordTooLong) {
		t.Fatalf("got %v", err)
	}
}
