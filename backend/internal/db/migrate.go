package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/pressly/goose/v3"
)

// goose 的 dialect 與 base FS 是套件層級的全域狀態，並行測試下重複設定會觸發 race
var (
	gooseOnce sync.Once
	gooseErr  error
)

func setupGoose() error {
	gooseOnce.Do(func() {
		goose.SetBaseFS(Migrations)
		goose.SetLogger(goose.NopLogger()) // 由呼叫端決定要印什麼
		gooseErr = goose.SetDialect("mysql")
	})
	if gooseErr != nil {
		return fmt.Errorf("設定 goose 失敗：%w", gooseErr)
	}
	return nil
}

// Migrate 套用所有尚未執行的遷移，回傳套用前後的版本號。
// 進度記在 goose_db_version 表，因此可重複執行。
func Migrate(ctx context.Context, conn *sql.DB) (before, after int64, err error) {
	if err := setupGoose(); err != nil {
		return 0, 0, err
	}
	// 全新的資料庫還沒有 goose_db_version 表，查不到版本就當作 0
	before, _ = goose.GetDBVersionContext(ctx, conn)
	if err := goose.UpContext(ctx, conn, "migrations"); err != nil {
		return before, 0, fmt.Errorf("套用遷移失敗：%w", err)
	}
	after, err = goose.GetDBVersionContext(ctx, conn)
	if err != nil {
		return before, 0, fmt.Errorf("讀取遷移版本失敗：%w", err)
	}
	return before, after, nil
}
