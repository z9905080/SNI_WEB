package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWrapRecoversAndSetsHeaders(t *testing.T) {
	h := Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if RequestID(r.Context()) == "" {
			t.Error("缺少 request id")
		}
		panic("boom")
	}), "default-src 'self'")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/api/v1/x", nil))
	if w.Code != 500 {
		t.Fatalf("status %d", w.Code)
	}
	id := w.Header().Get("X-Request-ID")
	if id == "" || !strings.Contains(w.Body.String(), id) || !strings.Contains(w.Body.String(), `"code":"internal"`) {
		t.Fatalf("回應應含 request id：%s / %s", id, w.Body.String())
	}
	for k, v := range map[string]string{
		"X-Content-Type-Options":  "nosniff",
		"Referrer-Policy":         "strict-origin-when-cross-origin",
		"Content-Security-Policy": "default-src 'self'",
	} {
		if got := w.Header().Get(k); got != v {
			t.Errorf("%s = %q", k, got)
		}
	}
}

func TestBuildCSP(t *testing.T) {
	csp := BuildCSP("http://toolkit.url.com.tw/counter/setcounter.php?sid=24073")
	for _, want := range []string{
		"default-src 'self'",
		"script-src 'self' https://www.googletagmanager.com http://toolkit.url.com.tw",
		"frame-src 'self' https://www.youtube.com https://youtube.com https://www.youtube-nocookie.com https://players.brightcove.net https://drive.google.com https://www.facebook.com",
		"img-src 'self' data: blob: https: http:",
		"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
		"font-src 'self' data: https://fonts.gstatic.com",
		"object-src 'none'",
		"base-uri 'self'",
		"frame-ancestors 'self'",
	} {
		if !strings.Contains(csp, want) {
			t.Errorf("CSP 缺少 %q\n%s", want, csp)
		}
	}
	if strings.Contains(BuildCSP(""), "toolkit") {
		t.Error("沒有計數器網址時不應加入")
	}
}
