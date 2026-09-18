package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/z9905080/SNI_WEB/backend/internal/store"
)

const (
	publicFixture = `INSERT INTO page_group (id, group_name, page_sort) VALUES (1, '首頁', '[2,1]'), (2, '活動', '')`
	publicPages   = `INSERT INTO page_content (id, page_group_id, page_name, html_context) VALUES (1, 1, '舊消息', '<p>1</p>'), (2, 1, '最新消息', '<p>2</p>'), (3, 2, '練成會', '<p>3</p>')`
	publicConfig  = `INSERT INTO web_config (data_key, data_value) VALUES ('web_title', '生長之家'), ('web_sub_title', '感謝'), ('facebook_url', 'https://fb'), ('page_group_sort', '[2,1]')`
	publicMedia   = `INSERT INTO carousel (id, image, url) VALUES (1, '/php/picture/a.jpg', '')`
)

type siteBody struct {
	store.Settings
	Menu []store.Group `json:"menu"`
}

type pageBody struct {
	Page      *store.Page      `json:"page"`
	Carousels []store.Carousel `json:"carousels"`
	Marquees  []store.Marquee  `json:"marquees"`
}

func TestGetSite(t *testing.T) {
	h := newHarness(t, publicFixture, publicPages, publicConfig)
	resp := h.do("GET", "/api/v1/site", nil)
	expectStatus(t, resp, 200)
	if cc := resp.Header.Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("Cache-Control %q", cc)
	}
	b := decode[siteBody](t, resp)
	if b.WebTitle != "生長之家" || b.FacebookURL != "https://fb" || len(b.Menu) != 2 || b.Menu[0].ID != 2 || b.Menu[1].Pages[0].ID != 2 {
		t.Fatalf("%+v", b)
	}
}

func TestGetHome(t *testing.T) {
	h := newHarness(t, publicFixture, publicPages, publicMedia)
	b := decode[pageBody](t, h.do("GET", "/api/v1/home", nil))
	if b.Page == nil || b.Page.ID != 2 || len(b.Carousels) != 1 || b.Marquees == nil {
		t.Fatalf("%+v", b)
	}

	h = newHarness(t)
	resp := h.do("GET", "/api/v1/home", nil)
	expectStatus(t, resp, 200)
	b = decode[pageBody](t, resp)
	if b.Page != nil || b.Carousels == nil {
		t.Fatalf("沒有首頁內容時 page 為 null：%+v", b)
	}
}

func TestGetPage(t *testing.T) {
	h := newHarness(t, publicFixture, publicPages)
	b := decode[pageBody](t, h.do("GET", "/api/v1/pages/3", nil))
	if b.Page == nil || *b.Page != (store.Page{ID: 3, GroupID: 2, Name: "練成會", HTML: "<p>3</p>"}) {
		t.Fatalf("%+v", b.Page)
	}
	expectError(t, h.do("GET", "/api/v1/pages/99", nil), 404, "not_found")
	expectError(t, h.do("GET", "/api/v1/pages/abc", nil), 400, "invalid_id")
}

func TestHealthAndUnknownAPI(t *testing.T) {
	h := newHarness(t)
	expectStatus(t, h.do("GET", "/healthz", nil), 200)
	expectStatus(t, h.do("GET", "/readyz", nil), 200)
	expectError(t, h.do("GET", "/api/v1/nope", nil), 404, "not_found")
	expectError(t, h.do("DELETE", "/api/v1/site", nil), 404, "not_found")
}

func TestReadyzTimesOutOnSlowPing(t *testing.T) {
	slow := func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}
	s := New(Deps{Ping: slow})
	srv := httptest.NewServer(Wrap(s.Routes(), "default-src 'self'"))
	t.Cleanup(srv.Close)

	start := time.Now()
	resp, err := http.Get(srv.URL + "/readyz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Fatalf("/readyz 應在約 2 秒逾時，實際花了 %v", elapsed)
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("status = %d，want 503", resp.StatusCode)
	}
}
