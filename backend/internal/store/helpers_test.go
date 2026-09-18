package store

import (
	"database/sql"
	"testing"

	"github.com/z9905080/SNI_WEB/backend/internal/testutil"
)

func newTestStore(t *testing.T, stmts ...string) (*Store, *sql.DB) {
	t.Helper()
	conn, _ := testutil.MySQL(t)
	testutil.Exec(t, conn, stmts...)
	return New(conn), conn
}

func pageIDs(g Group) []int64 {
	ids := []int64{}
	for _, p := range g.Pages {
		ids = append(ids, p.ID)
	}
	return ids
}

func groupIDs(gs []Group) []int64 {
	ids := []int64{}
	for _, g := range gs {
		ids = append(ids, g.ID)
	}
	return ids
}

func queryString(t *testing.T, conn *sql.DB, q string, args ...any) string {
	t.Helper()
	var s string
	if err := conn.QueryRow(q, args...).Scan(&s); err != nil {
		t.Fatalf("%s: %v", q, err)
	}
	return s
}
