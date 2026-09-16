package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSetAndClearCookie(t *testing.T) {
	w := httptest.NewRecorder()
	exp := time.Date(2026, 1, 1, 1, 0, 0, 0, time.UTC)
	SetCookie(w, "tok", exp, true)
	c := w.Result().Cookies()[0]
	if c.Name != CookieName || c.Value != "tok" || !c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteLaxMode || c.Path != "/" || !c.Expires.Equal(exp) {
		t.Fatalf("cookie 屬性錯誤：%+v", c)
	}

	w = httptest.NewRecorder()
	ClearCookie(w, false)
	c = w.Result().Cookies()[0]
	if c.Name != CookieName || c.MaxAge >= 0 || c.Secure {
		t.Fatalf("清除 cookie 屬性錯誤：%+v", c)
	}
}

func TestOriginAllowed(t *testing.T) {
	allowed := []string{"https://www.seicho-no-ie.org.tw", "http://localhost:5173"}
	cases := []struct {
		origin, referer string
		want            bool
	}{
		{"https://www.seicho-no-ie.org.tw", "", true},
		{"HTTPS://WWW.SEICHO-NO-IE.ORG.TW", "", true},
		{"http://localhost:5173", "", true},
		{"https://evil.example", "", false},
		{"null", "", false},
		{"", "https://www.seicho-no-ie.org.tw/admin/groups", true},
		{"", "https://evil.example/admin", false},
		{"", "not a url", false},
		{"", "", false},
	}
	for _, c := range cases {
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		if c.origin != "" {
			r.Header.Set("Origin", c.origin)
		}
		if c.referer != "" {
			r.Header.Set("Referer", c.referer)
		}
		if got := OriginAllowed(r, allowed); got != c.want {
			t.Errorf("origin=%q referer=%q got %v", c.origin, c.referer, got)
		}
	}
}

func TestClientIP(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.RemoteAddr = "10.0.0.1:5555"
	if got := ClientIP(r); got != "10.0.0.1" {
		t.Errorf("RemoteAddr → %q", got)
	}
	r.Header.Set("X-Forwarded-For", "1.1.1.1, 2.2.2.2 ")
	if got := ClientIP(r); got != "2.2.2.2" {
		t.Errorf("應取最右側：%q", got)
	}
}
