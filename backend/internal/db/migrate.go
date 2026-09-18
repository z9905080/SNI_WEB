package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
)

// newProvider 用 goose 的 Provider API：沒有套件層級全域狀態，並且能加鎖。
func newProvider(conn *sql.DB) (*goose.Provider, error) {
	// Provider 預期遷移檔位於 fsys 根目錄
	fsys, err := fs.Sub(Migrations, "migrations")
	if err != nil {
		return nil, err
	}
	// serve 啟動時會套用遷移，多副本同時啟動需要鎖來避免重複套用
	locker, err := lock.NewMySQLTableLocker()
	if err != nil {
		return nil, err
	}
	return goose.NewProvider(goose.DialectMySQL, conn, fsys, goose.WithLocker(locker))
}

// Migrate 套用所有尚未執行的遷移，回傳實際套用的版本號（已是最新版時為空）。
// 進度記在 goose_db_version 表，因此可重複執行。
//
// 只呼叫一次 Up：每次 Provider 操作都會重新取鎖，分開呼叫 GetDBVersion
// 會與 Up 互相卡住，直到鎖的租約逾時。
func Migrate(ctx context.Context, conn *sql.DB) ([]int64, error) {
	p, err := newProvider(conn)
	if err != nil {
		return nil, fmt.Errorf("建立 goose provider 失敗：%w", err)
	}
	results, err := p.Up(ctx)
	if err != nil {
		return nil, fmt.Errorf("套用遷移失敗：%w", err)
	}
	applied := make([]int64, 0, len(results))
	for _, r := range results {
		applied = append(applied, r.Source.Version)
	}
	return applied, nil
}
