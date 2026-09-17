// Package store 封裝網站內容的讀寫與排序 transaction。
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/z9905080/SNI_WEB/backend/internal/db/dbgen"
)

const (
	HomeGroupID    int64 = 1
	KeyWebTitle          = "web_title"
	KeyWebSubTitle       = "web_sub_title"
	KeyFacebookURL       = "facebook_url"
	KeyGroupSort         = "page_group_sort"
	MaxHTMLBytes         = 65535
)

var (
	ErrNotFound       = errors.New("store: not found")
	ErrGroupNotEmpty  = errors.New("store: group not empty")
	ErrGroupProtected = errors.New("store: group protected")
	ErrOrderMismatch  = errors.New("store: order mismatch")
	ErrInvalidGroup   = errors.New("store: invalid group")
	ErrContentTooLong = errors.New("store: content too long")
)

type Settings struct {
	WebTitle    string `json:"web_title"`
	WebSubTitle string `json:"web_sub_title"`
	FacebookURL string `json:"facebook_url"`
}

type PageRef struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Group struct {
	ID    int64     `json:"id"`
	Name  string    `json:"name"`
	Pages []PageRef `json:"pages"`
}

type Page struct {
	ID      int64  `json:"id"`
	GroupID int64  `json:"group_id"`
	Name    string `json:"name"`
	HTML    string `json:"html"`
}

type Carousel struct {
	ID    int64  `json:"id"`
	Image string `json:"image"`
	URL   string `json:"url"`
}

type Marquee struct {
	ID    int64  `json:"id"`
	Text  string `json:"text"`
	Color string `json:"color"`
}

type Store struct {
	db *sql.DB
	q  *dbgen.Queries
}

func New(conn *sql.DB) *Store {
	return &Store{db: conn, q: dbgen.New(conn)}
}

func (s *Store) inTx(ctx context.Context, fn func(q *dbgen.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(s.q.WithTx(tx)); err != nil {
		_ = tx.Rollback() // 回傳原本的錯誤；rollback 失敗時 driver 會關閉連線
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

// notFound 把 sql.ErrNoRows 轉成 ErrNotFound。
func notFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
