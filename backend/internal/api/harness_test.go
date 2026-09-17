package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/z9905080/SNI_WEB/backend/internal/auth"
	"github.com/z9905080/SNI_WEB/backend/internal/db/dbgen"
	"github.com/z9905080/SNI_WEB/backend/internal/storage"
	"github.com/z9905080/SNI_WEB/backend/internal/store"
	"github.com/z9905080/SNI_WEB/backend/internal/testutil"
)

const testOrigin = "http://example.test"

type clock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *clock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

type harness struct {
	t      *testing.T
	srv    *httptest.Server
	client *http.Client
	conn   *sql.DB
	q      *dbgen.Queries
	clock  *clock
	disk   *storage.Disk
}

func newHarness(t *testing.T, fixtures ...string) *harness {
	t.Helper()
	conn, _ := testutil.MySQL(t)
	testutil.Exec(t, conn, fixtures...)
	// cookiejar 以真實時間判斷 cookie 是否過期，假時鐘必須從現在開始
	c := &clock{now: time.Now().UTC().Truncate(time.Second)}
	disk, err := storage.NewDisk(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	q := dbgen.New(conn)
	s := New(Deps{
		Store:          store.New(conn),
		Sessions:       auth.NewSessions(q, c.Now),
		Storage:        disk,
		CookieSecure:   false,
		AllowedOrigins: []string{testOrigin},
		Now:            c.Now,
		Ping:           conn.PingContext,
	})
	srv := httptest.NewServer(Wrap(s.Routes(), "default-src 'self'"))
	t.Cleanup(srv.Close)
	jar, _ := cookiejar.New(nil)
	return &harness{t: t, srv: srv, client: &http.Client{Jar: jar}, conn: conn, q: q, clock: c, disk: disk}
}

// do 送出 JSON 請求（body 為 nil 時不帶 body），自動帶上允許的 Origin。
func (h *harness) do(method, path string, body any) *http.Response {
	h.t.Helper()
	var rd io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			h.t.Fatal(err)
		}
		rd = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, h.srv.URL+path, rd)
	if err != nil {
		h.t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Origin", testOrigin)
	return h.send(req)
}

func (h *harness) send(req *http.Request) *http.Response {
	h.t.Helper()
	resp, err := h.client.Do(req)
	if err != nil {
		h.t.Fatal(err)
	}
	h.t.Cleanup(func() { resp.Body.Close() })
	return resp
}

// login 建立帳號並登入，之後的請求會帶 cookie。
func (h *harness) login() {
	h.t.Helper()
	if _, err := auth.CreateUser(context.Background(), h.q, "admin", "管理者", "password1"); err != nil {
		h.t.Fatal(err)
	}
	resp := h.do("POST", "/api/v1/admin/auth/login", map[string]string{"account": "admin", "password": "password1"})
	expectStatus(h.t, resp, 200)
}

func expectStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode != want {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("%s %s: status %d, want %d, body %s", resp.Request.Method, resp.Request.URL.Path, resp.StatusCode, want, b)
	}
}

func decode[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	var v T
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return v
}

type errorResponse struct {
	Error struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		RequestID string `json:"request_id"`
	} `json:"error"`
}

func expectError(t *testing.T, resp *http.Response, status int, code string) errorResponse {
	t.Helper()
	expectStatus(t, resp, status)
	e := decode[errorResponse](t, resp)
	if e.Error.Code != code || e.Error.Message == "" {
		t.Fatalf("錯誤回應 %+v，want code %s", e, code)
	}
	return e
}

func mustURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
