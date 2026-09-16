package testutil

import (
	"testing"
)

func TestMySQLIsolatedDatabases(t *testing.T) {
	a, _ := MySQL(t)
	b, dsn := MySQL(t)
	if dsn == "" {
		t.Fatal("應回傳 DSN")
	}
	Exec(t, a, "INSERT INTO page_group (group_name, page_sort) VALUES ('甲', '')")

	var n int
	if err := b.QueryRow("SELECT COUNT(*) FROM page_group").Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("資料庫應彼此隔離，b 有 %d 筆", n)
	}
	var name string
	if err := a.QueryRow("SELECT group_name FROM page_group").Scan(&name); err != nil || name != "甲" {
		t.Fatalf("utf8mb4 讀寫失敗：%q %v", name, err)
	}
}
