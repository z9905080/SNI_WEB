// Package testutil 提供需要 Docker 的整合測試工具。
package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/z9905080/SNI_WEB/backend/internal/db"
)

var (
	once     sync.Once
	baseDSN  string
	startErr error
	counter  atomic.Int64
)

// MySQL 回傳一個已匯入 schema 的全新資料庫。整個測試程序共用一個容器（由 testcontainers 的 reaper 回收）。
func MySQL(t *testing.T) (*sql.DB, string) {
	t.Helper()
	if testing.Short() {
		t.Skip("需要 Docker，-short 時略過")
	}
	once.Do(func() {
		ctx := context.Background()
		c, err := mysql.Run(ctx, "mysql:8.4",
			mysql.WithUsername("root"),
			mysql.WithPassword("test"),
			mysql.WithDatabase("bootstrap"),
		)
		if err != nil {
			startErr = err
			return
		}
		host, err := c.Host(ctx)
		if err != nil {
			startErr = err
			return
		}
		port, err := c.MappedPort(ctx, "3306/tcp")
		if err != nil {
			startErr = err
			return
		}
		baseDSN = fmt.Sprintf("root:test@tcp(%s:%s)/", host, port.Port())
	})
	if startErr != nil {
		t.Fatalf("啟動 MySQL 容器失敗：%v", startErr)
	}

	name := fmt.Sprintf("t%d_%d", os.Getpid(), counter.Add(1))
	admin, err := sql.Open("mysql", baseDSN+"?multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	Exec(t, admin, "CREATE DATABASE `"+name+"` CHARACTER SET utf8mb4")

	dsn := baseDSN + name
	conn, err := db.Open(context.Background(), dsn)
	if err != nil {
		t.Fatal(err)
	}
	// 走正式環境同一條路徑：同樣用 db.Open 的連線設定（goose 的表格鎖靠
	// RowsAffected 判斷取鎖成敗，換成裸的 sql.Open 會一直取不到而卡住），
	// 遷移檔壞掉時測試就會抓到
	if _, err := db.Migrate(context.Background(), conn); err != nil {
		t.Fatalf("套用遷移失敗：%v", err)
	}
	t.Cleanup(func() {
		conn.Close()
		if a, err := sql.Open("mysql", baseDSN); err == nil {
			_, _ = a.Exec("DROP DATABASE `" + name + "`")
			a.Close()
		}
	})
	return conn, dsn
}

func Exec(t *testing.T, conn *sql.DB, stmts ...string) {
	t.Helper()
	for _, s := range stmts {
		if strings.TrimSpace(s) == "" {
			continue
		}
		if _, err := conn.Exec(s); err != nil {
			t.Fatalf("執行 SQL 失敗：%v\n%s", err, s)
		}
	}
}
