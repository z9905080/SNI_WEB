// Package storage 抽象圖片檔的存放位置（本機磁碟或 S3 相容儲存）。
package storage

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"
)

var (
	ErrNotFound    = errors.New("storage: not found")
	ErrExists      = errors.New("storage: already exists")
	ErrInvalidName = errors.New("storage: invalid name")
)

type Object struct {
	Name    string
	Size    int64
	ModTime time.Time
}

type Storage interface {
	Put(ctx context.Context, name string, r io.ReadSeeker, size int64, contentType string) error
	Delete(ctx context.Context, name string) error
	List(ctx context.Context) ([]Object, error)
	Serve(w http.ResponseWriter, r *http.Request, name string)
}

const immutableCache = "public, max-age=31536000, immutable"
