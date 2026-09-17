package api

import (
	"testing"

	"github.com/z9905080/SNI_WEB/backend/internal/store"
)

func TestCarouselEndpoints(t *testing.T) {
	h := newHarness(t)
	h.login()

	for _, bad := range []map[string]string{
		{"image": "", "url": ""},
		{"image": "javascript:alert(1)", "url": ""},
		{"image": "/php/picture/a.jpg", "url": "javascript:alert(1)"},
		{"image": "/php/picture/a.jpg", "url": "//evil.example"},
	} {
		expectError(t, h.do("POST", "/api/v1/admin/carousels", bad), 400, "validation")
	}

	resp := h.do("POST", "/api/v1/admin/carousels", map[string]string{"image": "/php/picture/a.jpg", "url": "/page/3"})
	expectStatus(t, resp, 201)
	c := decode[struct {
		Carousel store.Carousel `json:"carousel"`
	}](t, resp).Carousel

	upd := decode[struct {
		Carousel store.Carousel `json:"carousel"`
	}](t, h.do("PATCH", "/api/v1/admin/carousels/"+itoa(c.ID), map[string]string{"url": "https://www.facebook.com/x"})).Carousel
	if upd.Image != "/php/picture/a.jpg" || upd.URL != "https://www.facebook.com/x" {
		t.Fatalf("%+v", upd)
	}

	list := decode[struct {
		Carousels []store.Carousel `json:"carousels"`
	}](t, h.do("GET", "/api/v1/admin/carousels", nil))
	if len(list.Carousels) != 1 || list.Carousels[0] != upd {
		t.Fatalf("%+v", list)
	}

	expectError(t, h.do("PATCH", "/api/v1/admin/carousels/999", map[string]string{"url": ""}), 404, "not_found")
	expectStatus(t, h.do("DELETE", "/api/v1/admin/carousels/"+itoa(c.ID), nil), 204)
	expectError(t, h.do("DELETE", "/api/v1/admin/carousels/"+itoa(c.ID), nil), 404, "not_found")
}

func TestMarqueeEndpoints(t *testing.T) {
	h := newHarness(t)
	h.login()

	expectError(t, h.do("POST", "/api/v1/admin/marquees", map[string]string{"text": "", "color": "#fff"}), 400, "validation")
	expectError(t, h.do("POST", "/api/v1/admin/marquees", map[string]string{"text": "x", "color": "red"}), 400, "validation")

	resp := h.do("POST", "/api/v1/admin/marquees", map[string]string{"text": " 合掌感謝！ ", "color": "#1EFF00"})
	expectStatus(t, resp, 201)
	m := decode[struct {
		Marquee store.Marquee `json:"marquee"`
	}](t, resp).Marquee
	if m.Text != "合掌感謝！" {
		t.Fatalf("%+v", m)
	}

	upd := decode[struct {
		Marquee store.Marquee `json:"marquee"`
	}](t, h.do("PATCH", "/api/v1/admin/marquees/"+itoa(m.ID), map[string]string{"color": "#abc"})).Marquee
	if upd.Color != "#abc" || upd.Text != "合掌感謝！" {
		t.Fatalf("%+v", upd)
	}

	list := decode[struct {
		Marquees []store.Marquee `json:"marquees"`
	}](t, h.do("GET", "/api/v1/admin/marquees", nil))
	if len(list.Marquees) != 1 {
		t.Fatalf("%+v", list)
	}
	expectStatus(t, h.do("DELETE", "/api/v1/admin/marquees/"+itoa(m.ID), nil), 204)
}

func TestSettingsEndpoints(t *testing.T) {
	h := newHarness(t, `INSERT INTO web_config (data_key, data_value) VALUES ('web_title', '生長之家')`)
	h.login()

	got := decode[struct {
		Settings store.Settings `json:"settings"`
	}](t, h.do("GET", "/api/v1/admin/settings", nil)).Settings
	if got.WebTitle != "生長之家" {
		t.Fatalf("%+v", got)
	}

	expectError(t, h.do("PUT", "/api/v1/admin/settings", store.Settings{WebTitle: ""}), 400, "validation")
	expectError(t, h.do("PUT", "/api/v1/admin/settings", store.Settings{WebTitle: "a", FacebookURL: "ftp://x"}), 400, "validation")

	want := store.Settings{WebTitle: "新標題", WebSubTitle: "副標", FacebookURL: "https://www.facebook.com/seichonoie.tw"}
	saved := decode[struct {
		Settings store.Settings `json:"settings"`
	}](t, h.do("PUT", "/api/v1/admin/settings", want)).Settings
	if saved != want {
		t.Fatalf("%+v", saved)
	}
	site := decode[siteBody](t, h.do("GET", "/api/v1/site", nil))
	if site.Settings != want {
		t.Fatalf("公開 API 應反映新設定：%+v", site.Settings)
	}
}

func TestSettingsUnsupportedCharacters(t *testing.T) {
	h := newHarness(t)
	h.login()
	expectError(t, h.do("PUT", "/api/v1/admin/settings", store.Settings{WebTitle: "感謝🙏"}), 400, "unsupported_characters")
}
