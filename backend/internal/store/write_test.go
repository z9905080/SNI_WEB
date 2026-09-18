package store

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"
	"testing"
)

const (
	writeFixture = `INSERT INTO page_group (id, group_name, page_sort) VALUES (1, '首頁', ''), (2, '活動', '[11,10]'), (3, '空', '')`
	writePages   = `INSERT INTO page_content (id, page_group_id, page_name, html_context) VALUES (10, 2, 'A', ''), (11, 2, 'B', ''), (12, 1, 'H', '')`
	writeConfig  = `INSERT INTO web_config (data_key, data_value) VALUES ('page_group_sort', '[2,1,3]')`
)

var ctx = context.Background()

func TestCreateRenameGroup(t *testing.T) {
	s, conn := newTestStore(t, writeFixture, writeConfig)
	g, err := s.CreateGroup(ctx, "新群組")
	if err != nil {
		t.Fatal(err)
	}
	if g.ID == 0 || g.Name != "新群組" || g.Pages == nil {
		t.Fatalf("got %+v", g)
	}
	if got := queryString(t, conn, "SELECT data_value FROM web_config WHERE data_key='page_group_sort'"); got != "[2,1,3,"+itoa(g.ID)+"]" {
		t.Fatalf("排序未附加：%s", got)
	}

	if err := s.RenameGroup(ctx, g.ID, "改名"); err != nil {
		t.Fatal(err)
	}
	if err := s.RenameGroup(ctx, g.ID, "改名"); err != nil {
		t.Fatalf("名稱未變也應成功：%v", err)
	}
	if err := s.RenameGroup(ctx, 999, "x"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
}

func TestCreateGroupWithoutConfigRow(t *testing.T) {
	s, conn := newTestStore(t)
	g, err := s.CreateGroup(ctx, "第一個")
	if err != nil {
		t.Fatal(err)
	}
	if got := queryString(t, conn, "SELECT data_value FROM web_config WHERE data_key='page_group_sort'"); got != "["+itoa(g.ID)+"]" {
		t.Fatalf("應新增設定列：%s", got)
	}
}

func TestDeleteGroupRules(t *testing.T) {
	s, conn := newTestStore(t, writeFixture, writePages, writeConfig)
	if err := s.DeleteGroup(ctx, 1); !errors.Is(err, ErrGroupProtected) {
		t.Fatalf("首頁：%v", err)
	}
	if err := s.DeleteGroup(ctx, 2); !errors.Is(err, ErrGroupNotEmpty) {
		t.Fatalf("有頁面：%v", err)
	}
	if err := s.DeleteGroup(ctx, 999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("不存在：%v", err)
	}
	if err := s.DeleteGroup(ctx, 3); err != nil {
		t.Fatal(err)
	}
	if got := queryString(t, conn, "SELECT data_value FROM web_config WHERE data_key='page_group_sort'"); got != "[2,1]" {
		t.Fatalf("排序未移除：%s", got)
	}
}

func TestSetGroupOrder(t *testing.T) {
	s, conn := newTestStore(t, writeFixture, writeConfig)
	if err := s.SetGroupOrder(ctx, []int64{3, 2}); !errors.Is(err, ErrOrderMismatch) {
		t.Fatalf("少一個：%v", err)
	}
	if err := s.SetGroupOrder(ctx, []int64{3, 2, 1, 9}); !errors.Is(err, ErrOrderMismatch) {
		t.Fatalf("多一個：%v", err)
	}
	if err := s.SetGroupOrder(ctx, []int64{3, 1, 2}); err != nil {
		t.Fatal(err)
	}
	if got := queryString(t, conn, "SELECT data_value FROM web_config WHERE data_key='page_group_sort'"); got != "[3,1,2]" {
		t.Fatalf("got %s", got)
	}
}

func TestSetPageOrder(t *testing.T) {
	s, _ := newTestStore(t, writeFixture, writePages)
	if err := s.SetPageOrder(ctx, 2, []int64{10}); !errors.Is(err, ErrOrderMismatch) {
		t.Fatalf("got %v", err)
	}
	if err := s.SetPageOrder(ctx, 2, []int64{10, 12}); !errors.Is(err, ErrOrderMismatch) {
		t.Fatalf("別的群組的頁面：%v", err)
	}
	if err := s.SetPageOrder(ctx, 999, []int64{}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
	if err := s.SetPageOrder(ctx, 2, []int64{10, 11}); err != nil {
		t.Fatal(err)
	}
	nav, _ := s.Nav(ctx)
	for _, g := range nav {
		if g.ID == 2 && !slices.Equal(pageIDs(g), []int64{10, 11}) {
			t.Fatalf("got %v", pageIDs(g))
		}
	}
}

func TestCreatePage(t *testing.T) {
	s, conn := newTestStore(t, writeFixture, writePages)
	p, err := s.CreatePage(ctx, Page{GroupID: 2, Name: "新頁", HTML: `<p onclick="x()">內容</p>`})
	if err != nil {
		t.Fatal(err)
	}
	if p.ID == 0 || p.HTML != "<p>內容</p>" {
		t.Fatalf("got %+v", p)
	}
	if got := queryString(t, conn, "SELECT page_sort FROM page_group WHERE id=2"); got != "[11,10,"+itoa(p.ID)+"]" {
		t.Fatalf("排序未附加：%s", got)
	}
	if _, err := s.CreatePage(ctx, Page{GroupID: 999, Name: "x"}); !errors.Is(err, ErrInvalidGroup) {
		t.Fatalf("got %v", err)
	}
	if _, err := s.CreatePage(ctx, Page{GroupID: 2, Name: "x", HTML: "<p>" + strings.Repeat("字", 22000) + "</p>"}); !errors.Is(err, ErrContentTooLong) {
		t.Fatalf("got %v", err)
	}
}

func TestUpdatePageMovesGroup(t *testing.T) {
	s, conn := newTestStore(t, writeFixture, writePages)
	p, err := s.UpdatePage(ctx, Page{ID: 10, GroupID: 3, Name: "A2", HTML: "<p>x</p>"})
	if err != nil {
		t.Fatal(err)
	}
	if p != (Page{ID: 10, GroupID: 3, Name: "A2", HTML: "<p>x</p>"}) {
		t.Fatalf("got %+v", p)
	}
	if got := queryString(t, conn, "SELECT page_sort FROM page_group WHERE id=2"); got != "[11]" {
		t.Fatalf("舊群組：%s", got)
	}
	if got := queryString(t, conn, "SELECT page_sort FROM page_group WHERE id=3"); got != "[10]" {
		t.Fatalf("新群組：%s", got)
	}

	// 同群組更新不動排序
	if _, err := s.UpdatePage(ctx, Page{ID: 11, GroupID: 2, Name: "B2"}); err != nil {
		t.Fatal(err)
	}
	if got := queryString(t, conn, "SELECT page_sort FROM page_group WHERE id=2"); got != "[11]" {
		t.Fatalf("同群組：%s", got)
	}

	if _, err := s.UpdatePage(ctx, Page{ID: 999, GroupID: 2, Name: "x"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
	if _, err := s.UpdatePage(ctx, Page{ID: 11, GroupID: 999, Name: "x"}); !errors.Is(err, ErrInvalidGroup) {
		t.Fatalf("got %v", err)
	}
}

func TestDeletePage(t *testing.T) {
	s, conn := newTestStore(t, writeFixture, writePages)
	if err := s.DeletePage(ctx, 10); err != nil {
		t.Fatal(err)
	}
	if got := queryString(t, conn, "SELECT page_sort FROM page_group WHERE id=2"); got != "[11]" {
		t.Fatalf("got %s", got)
	}
	if _, err := s.Page(ctx, 10); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
	if err := s.DeletePage(ctx, 10); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
}

func TestUpdateSettings(t *testing.T) {
	s, _ := newTestStore(t, `INSERT INTO web_config (data_key, data_value) VALUES ('web_title', '舊')`)
	want := Settings{WebTitle: "新", WebSubTitle: "副", FacebookURL: "https://www.facebook.com/x"}
	if err := s.UpdateSettings(ctx, want); err != nil {
		t.Fatal(err)
	}
	got, _ := s.Settings(ctx)
	if got != want {
		t.Fatalf("got %+v", got)
	}
}

func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}
