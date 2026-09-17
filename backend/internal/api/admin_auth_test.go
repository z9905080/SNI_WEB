package api

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/z9905080/SNI_WEB/backend/internal/auth"
)

type userBody struct {
	User auth.User `json:"user"`
}

func TestLoginMeLogout(t *testing.T) {
	h := newHarness(t)
	expectError(t, h.do("GET", "/api/v1/admin/auth/me", nil), 401, "unauthorized")

	h.login()
	cookies := h.client.Jar.Cookies(mustURL(h.srv.URL))
	if len(cookies) != 1 || cookies[0].Name != auth.CookieName {
		t.Fatalf("cookie %+v", cookies)
	}

	me := decode[userBody](t, h.do("GET", "/api/v1/admin/auth/me", nil))
	if me.User.Account != "admin" || me.User.Name != "管理者" {
		t.Fatalf("%+v", me)
	}

	expectStatus(t, h.do("POST", "/api/v1/admin/auth/logout", nil), 204)
	expectError(t, h.do("GET", "/api/v1/admin/auth/me", nil), 401, "unauthorized")
}

func TestLoginValidationAndFailure(t *testing.T) {
	h := newHarness(t)
	h.login()
	expectError(t, h.do("POST", "/api/v1/admin/auth/login", map[string]string{"account": " ", "password": "x"}), 400, "validation")
	e := expectError(t, h.do("POST", "/api/v1/admin/auth/login", map[string]string{"account": "admin", "password": "nope"}), 401, "invalid_credentials")
	e2 := expectError(t, h.do("POST", "/api/v1/admin/auth/login", map[string]string{"account": "ghost", "password": "nope"}), 401, "invalid_credentials")
	if e.Error.Message != e2.Error.Message {
		t.Fatal("帳號不存在與密碼錯誤的訊息應相同")
	}
}

func TestSessionExpiryAndSliding(t *testing.T) {
	h := newHarness(t)
	h.login()

	h.clock.Advance(45 * time.Minute)
	resp := h.do("GET", "/api/v1/admin/auth/me", nil)
	expectStatus(t, resp, 200)
	if len(resp.Cookies()) != 1 {
		t.Fatal("剩餘少於 30 分鐘時應重設 cookie")
	}

	h.clock.Advance(59 * time.Minute) // 延長後仍在 1 小時內
	expectStatus(t, h.do("GET", "/api/v1/admin/auth/me", nil), 200)

	h.clock.Advance(2 * time.Hour)
	resp = h.do("GET", "/api/v1/admin/auth/me", nil)
	expectError(t, resp, 401, "unauthorized")
}

func TestAccountRateLimit(t *testing.T) {
	h := newHarness(t)
	h.login()
	for i := 0; i < 10; i++ {
		expectStatus(t, h.do("POST", "/api/v1/admin/auth/login", map[string]string{"account": "admin", "password": "bad"}), 401)
	}
	expectError(t, h.do("POST", "/api/v1/admin/auth/login", map[string]string{"account": "admin", "password": "password1"}), 429, "too_many_requests")

	h.clock.Advance(16 * time.Minute)
	expectStatus(t, h.do("POST", "/api/v1/admin/auth/login", map[string]string{"account": "admin", "password": "password1"}), 200)
}

func TestAccountRateLimitCaseInsensitive(t *testing.T) {
	h := newHarness(t)
	h.login()
	variants := []string{"admin", "Admin", "ADMIN"}
	for i := 0; i < 10; i++ {
		account := variants[i%len(variants)]
		expectStatus(t, h.do("POST", "/api/v1/admin/auth/login", map[string]string{"account": account, "password": "bad"}), 401)
	}
	expectError(t, h.do("POST", "/api/v1/admin/auth/login", map[string]string{"account": "admin", "password": "password1"}), 429, "too_many_requests")
}

func TestIPRateLimit(t *testing.T) {
	h := newHarness(t)
	for i := 0; i < 20; i++ {
		expectStatus(t, h.do("POST", "/api/v1/admin/auth/login", map[string]string{"account": fmt.Sprintf("u%d", i), "password": "bad"}), 401)
	}
	expectError(t, h.do("POST", "/api/v1/admin/auth/login", map[string]string{"account": "other", "password": "bad"}), 429, "too_many_requests")
}

func TestCSRF(t *testing.T) {
	h := newHarness(t)
	h.login()

	req, _ := http.NewRequest("POST", h.srv.URL+"/api/v1/admin/auth/logout", nil)
	req.Header.Set("Origin", "https://evil.example")
	expectError(t, h.send(req), 403, "forbidden_origin")

	req, _ = http.NewRequest("POST", h.srv.URL+"/api/v1/admin/auth/login", strings.NewReader(`{"account":"admin","password":"password1"}`))
	expectError(t, h.send(req), 403, "forbidden_origin")

	req, _ = http.NewRequest("GET", h.srv.URL+"/api/v1/admin/auth/me", nil)
	expectStatus(t, h.send(req), 200) // GET 不檢查來源
}
