// Package auth 處理密碼、登入 session、限流與來源檢查。
package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const MinPasswordLen = 8

var (
	ErrPasswordTooShort = errors.New("密碼至少需要 8 個字元")
	ErrPasswordTooLong  = errors.New("密碼不可超過 72 bytes")
)

func HashPassword(pw string) (string, error) {
	if len(pw) < MinPasswordLen {
		return "", ErrPasswordTooShort
	}
	if len(pw) > 72 {
		return "", ErrPasswordTooLong
	}
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(h), err
}

func CheckPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}
