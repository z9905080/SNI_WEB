package store

import (
	"context"
	"errors"
	"slices"
	"testing"
)

const navFixture = `INSERT INTO page_group (id, group_name, page_sort) VALUES
 (1, '首頁', ''), (2, '活動', '[5,99,3]'), (3, '簡介', 'bad json'), (4, '孤兒排序', '[]')`

const pagesFixture = `INSERT INTO page_content (id, page_group_id, page_name, html_context) VALUES
 (3, 2, 'C', '<p>c</p>'), (4, 2, 'D', ''), (5, 2, 'E', ''), (6, 3, 'F', ''), (7, 1, 'H2', '<p>第二</p>'), (8, 1, 'H1', '<p>第一</p>')`

func TestNavOrdering(t *testing.T) {
	s, _ := newTestStore(t, navFixture, pagesFixture,
		`INSERT INTO web_config (data_key, data_value) VALUES ('page_group_sort', '[3,1]')`)
	nav, err := s.Nav(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := groupIDs(nav); !slices.Equal(got, []int64{3, 1, 2, 4}) {
		t.Fatalf("群組順序 %v", got)
	}
	if got := pageIDs(nav[2]); !slices.Equal(got, []int64{5, 3, 4}) {
		t.Fatalf("活動頁面順序 %v", got)
	}
	if got := pageIDs(nav[0]); !slices.Equal(got, []int64{6}) {
		t.Fatalf("無效 JSON 應視為空陣列：%v", got)
	}
	if nav[3].Pages == nil || len(nav[3].Pages) != 0 {
		t.Fatalf("空群組 Pages 應為空切片：%#v", nav[3].Pages)
	}
	if nav[2].Name != "活動" || nav[2].Pages[0].Name != "E" {
		t.Fatalf("名稱錯誤：%+v", nav[2])
	}
}

func TestNavWithoutGroupSortConfig(t *testing.T) {
	s, _ := newTestStore(t, navFixture)
	nav, err := s.Nav(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := groupIDs(nav); !slices.Equal(got, []int64{1, 2, 3, 4}) {
		t.Fatalf("沒有設定時依 id：%v", got)
	}
}

func TestSettings(t *testing.T) {
	s, _ := newTestStore(t, `INSERT INTO web_config (data_key, data_value) VALUES
	 ('web_title', '生長之家'), ('web_sub_title', '感謝'), ('page_group_sort', '[1]')`)
	got, err := s.Settings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got != (Settings{WebTitle: "生長之家", WebSubTitle: "感謝"}) {
		t.Fatalf("got %+v", got)
	}
}

func TestHomePage(t *testing.T) {
	ctx := context.Background()

	s, _ := newTestStore(t, navFixture)
	if _, ok, err := s.HomePage(ctx); err != nil || ok {
		t.Fatalf("沒有首頁頁面時 ok 應為 false：%v %v", ok, err)
	}

	s, _ = newTestStore(t, navFixture, pagesFixture)
	p, ok, err := s.HomePage(ctx)
	if err != nil || !ok {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	if p.ID != 7 || p.HTML != "<p>第二</p>" {
		t.Fatalf("page_sort 為空時取 id 最小者：%+v", p)
	}

	s, _ = newTestStore(t,
		`INSERT INTO page_group (id, group_name, page_sort) VALUES (1, '首頁', '[8,7]')`, pagesFixture)
	p, _, _ = s.HomePage(ctx)
	if p.ID != 8 || p.GroupID != 1 || p.Name != "H1" {
		t.Fatalf("應依排序取第一個：%+v", p)
	}
}

func TestPage(t *testing.T) {
	s, _ := newTestStore(t, pagesFixture)
	p, err := s.Page(context.Background(), 3)
	if err != nil {
		t.Fatal(err)
	}
	if p != (Page{ID: 3, GroupID: 2, Name: "C", HTML: "<p>c</p>"}) {
		t.Fatalf("got %+v", p)
	}
	if _, err := s.Page(context.Background(), 404); !errors.Is(err, ErrNotFound) {
		t.Fatalf("應回 ErrNotFound：%v", err)
	}
}

func TestCarouselsAndMarqueesReadEmpty(t *testing.T) {
	s, _ := newTestStore(t)
	cs, err := s.Carousels(context.Background())
	if err != nil || cs == nil || len(cs) != 0 {
		t.Fatalf("%#v %v", cs, err)
	}
	ms, err := s.Marquees(context.Background())
	if err != nil || ms == nil || len(ms) != 0 {
		t.Fatalf("%#v %v", ms, err)
	}
}
