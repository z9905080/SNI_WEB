// Package db 提供 schema 與資料庫連線；查詢程式碼位於 dbgen（sqlc 產生）。
package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Migrations 是 goose 的遷移檔；sqlc 也從同一個目錄推導 schema，兩者不會分歧。
//
//go:embed migrations/*.sql
var Migrations embed.FS

// NormalizeDSN 接受 go-sql-driver DSN 或 mysql:// URL，並強制 parseTime、UTC、clientFoundRows 與 utf8mb4。
func NormalizeDSN(dsn string) (string, error) {
	if strings.HasPrefix(dsn, "mysql://") {
		u, err := url.Parse(dsn)
		if err != nil {
			return "", fmt.Errorf("DATABASE_URL 格式錯誤：%w", err)
		}
		pw, _ := u.User.Password()
		dsn = fmt.Sprintf("%s:%s@tcp(%s)%s", u.User.Username(), pw, u.Host, u.Path)
		if u.RawQuery != "" {
			dsn += "?" + u.RawQuery
		}
	}
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return "", fmt.Errorf("DATABASE_URL 格式錯誤：%w", err)
	}
	cfg.ParseTime = true
	cfg.Loc = time.UTC
	cfg.ClientFoundRows = true
	// 沒有明確設定時給預設逾時，避免半死的連線讓請求無限期卡住
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 30 * time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 30 * time.Second
	}
	out := cfg.FormatDSN()
	if !strings.Contains(out, "charset=") {
		sep := "?"
		if strings.Contains(out, "?") {
			sep = "&"
		}
		out += sep + "charset=utf8mb4"
	}
	return out, nil
}

func Open(ctx context.Context, dsn string) (*sql.DB, error) {
	normalized, err := NormalizeDSN(dsn)
	if err != nil {
		return nil, err
	}
	conn, err := sql.Open("mysql", normalized)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)
	conn.SetConnMaxIdleTime(time.Minute)
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := conn.PingContext(pingCtx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("無法連線資料庫：%w", err)
	}
	return conn, nil
}
