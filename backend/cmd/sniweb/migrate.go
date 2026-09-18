package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/z9905080/SNI_WEB/backend/internal/db"
)

// migrateCmd 只需要 DATABASE_URL，不必湊齊 serve 用的其他設定。
func migrateCmd(ctx context.Context, stdout io.Writer, getenv func(string) string) error {
	dsn := getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("缺少環境變數：DATABASE_URL")
	}
	conn, err := db.Open(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close()

	applied, err := db.Migrate(ctx, conn)
	if err != nil {
		return err
	}
	if len(applied) == 0 {
		fmt.Fprintln(stdout, "資料庫已是最新版本")
		return nil
	}
	fmt.Fprintf(stdout, "遷移完成，套用版本：%v\n", applied)
	return nil
}
