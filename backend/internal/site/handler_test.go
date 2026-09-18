package site

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/z9905080/SNI_WEB/backend/internal/storage"
	"github.com/z9905080/SNI_WEB/backend/internal/store"
	"github.com/z9905080/SNI_WEB/backend/internal/testutil"
)

func newTestHandler(t *testing.T, dist fstest.MapFS) *Handler {
	t.Helper()
	conn, _ := testutil.MySQL(t)
	testutil.Exec(t, conn,
		`INSERT INTO web_config (data_key, data_value) VALUES ('web_title', '生長之家'), ('web_sub_title', '感謝')`,
		`INSERT INTO page_group (id, group_name, page_sort) VALUES (1, '首頁', ''), (2, '活動', '')`,
		`INSERT INTO page_content (id, page_group_id, page_name, html_context) VALUES
		 (1, 1, '最新消息', '<p>首頁的&nbsp;內容</p>'),
		 (3, 2, '練成會', '<p>日程</p><img src="/php/picture/c.jpg">'),
		 (4, 2, '<b>誌友會</b>', '<p>無圖</p>')`,
		`INSERT INTO carousel (id, image, url) VALUES (1, '/php/picture/banner.jpg', '')`)
	disk, err := storage.NewDisk(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := disk.Put(context.Background(), "c.jpg", strings.NewReader("jpg"), 3, "image/jpeg"); err != nil {
		t.Fatal(err)
	}
	return New(dist, store.New(conn), disk, "https://example.org", PublicConfig{GAMeasurementID: "G-1"})
}

var testDist = fstest.MapFS{
	"index.html":      {Data: []byte(indexTmpl)},
	"assets/app-1.js": {Data: []byte("console.log(1)")},
	"favicon.png":     {Data: []byte("png")},
}

func get(t *testing.T, h http.Handler, method, path string) (*http.Response, string) {
	t.Helper()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	resp := w.Result()
	b, _ := io.ReadAll(resp.Body)
	return resp, string(b)
}

func TestSPAHandler(t *testing.T) {
	h := newTestHandler(t, testDist)

	tests := []struct {
		name, path  string
		status      int
		contains    []string
		header, val string
	}{
		{"首頁", "/", 200, []string{
			"<title>生長之家｜感謝</title>",
			`content="首頁的 內容"`,
			`<meta property="og:image" content="https://example.org/php/picture/banner.jpg">`,
			`<meta property="og:url" content="https://example.org/">`,
			`"gaMeasurementId":"G-1"`,
		}, "Cache-Control", "no-cache"},
		{"內頁用內容第一張圖", "/page/3", 200, []string{
			"<title>練成會｜生長之家</title>",
			`content="日程"`,
			`content="https://example.org/php/picture/c.jpg"`,
		}, "", ""},
		{"內頁無圖用輪播圖且名稱跳脫", "/page/4", 200, []string{
			"<title>&lt;b&gt;誌友會&lt;/b&gt;｜生長之家</title>",
			`content="https://example.org/php/picture/banner.jpg"`,
		}, "", ""},
		{"不存在的頁面", "/page/99", 404, []string{`<div id="root">`}, "", ""},
		{"頁面編號格式錯誤", "/page/abc", 404, []string{`<div id="root">`}, "", ""},
		{"未知路徑", "/nope", 404, []string{"<title>生長之家</title>", `content="感謝"`}, "", ""},
		{"後台", "/admin/groups", 200, []string{`<meta name="robots" content="noindex">`}, "X-Robots-Tag", "noindex"},
		{"assets 長效快取", "/assets/app-1.js", 200, []string{"console.log(1)"}, "Cache-Control", "public, max-age=31536000, immutable"},
		{"其他靜態檔", "/favicon.png", 200, []string{"png"}, "Cache-Control", "public, max-age=3600"},
		{"index.html 不直接回傳", "/index.html", 404, []string{"<title>"}, "", ""},
		{"圖片", "/php/picture/c.jpg", 200, []string{"jpg"}, "Content-Type", "image/jpeg"},
		{"圖片不存在", "/php/picture/none.jpg", 404, nil, "", ""},
		{"路徑穿越", "/assets/../../etc/passwd", 404, nil, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := get(t, h, "GET", tt.path)
			if resp.StatusCode != tt.status {
				t.Fatalf("status %d, want %d\n%s", resp.StatusCode, tt.status, body)
			}
			for _, c := range tt.contains {
				if !strings.Contains(body, c) {
					t.Errorf("缺少 %s\n%s", c, body)
				}
			}
			if tt.header != "" && resp.Header.Get(tt.header) != tt.val {
				t.Errorf("%s = %q, want %q", tt.header, resp.Header.Get(tt.header), tt.val)
			}
		})
	}

	if resp, _ := get(t, h, "POST", "/"); resp.StatusCode != 405 {
		t.Errorf("POST 應回 405，得到 %d", resp.StatusCode)
	}
}

func TestSPAHandlerWithoutBuild(t *testing.T) {
	h := newTestHandler(t, fstest.MapFS{})
	if resp, _ := get(t, h, "GET", "/"); resp.StatusCode != 503 {
		t.Fatalf("前端未建置應回 503，得到 %d", resp.StatusCode)
	}
	if resp, _ := get(t, h, "GET", "/php/picture/c.jpg"); resp.StatusCode != 200 {
		t.Fatalf("圖片仍應可用，得到 %d", resp.StatusCode)
	}
}

func TestStaticAssetGzip(t *testing.T) {
	// 需要夠大且可壓縮的內容，太小的檔案壓完反而變大，precompress 會略過
	big := strings.Repeat("export const x = 'seicho-no-ie';\n", 500)
	dist := fstest.MapFS{
		"index.html":        {Data: []byte(indexTmpl)},
		"assets/big-1.js":   {Data: []byte(big)},
		"assets/tiny-1.js":  {Data: []byte("console.log(1)")},
		"assets/logo-1.png": {Data: []byte("\x89PNG not really")},
	}
	h := newTestHandler(t, dist)

	t.Run("接受 gzip 時回傳壓縮內容", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/assets/big-1.js", nil)
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if got := rec.Header().Get("Content-Encoding"); got != "gzip" {
			t.Fatalf("Content-Encoding = %q，want gzip", got)
		}
		if got := rec.Header().Get("Vary"); got != "Accept-Encoding" {
			t.Errorf("Vary = %q，want Accept-Encoding", got)
		}
		if rec.Body.Len() >= len(big) {
			t.Errorf("壓縮後 %d bytes 沒有比原始 %d bytes 小", rec.Body.Len(), len(big))
		}
		// Content-Length 必須是壓縮後的長度，否則用戶端會讀不完整
		if got := rec.Header().Get("Content-Length"); got != strconv.Itoa(rec.Body.Len()) {
			t.Errorf("Content-Length = %q，實際 body %d bytes", got, rec.Body.Len())
		}
		zr, err := gzip.NewReader(rec.Body)
		if err != nil {
			t.Fatal(err)
		}
		out, err := io.ReadAll(zr)
		if err != nil {
			t.Fatal(err)
		}
		if string(out) != big {
			t.Error("解壓後的內容與原始檔案不符")
		}
	})

	t.Run("不接受 gzip 時回傳原始內容", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/assets/big-1.js", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got := rec.Header().Get("Content-Encoding"); got != "" {
			t.Errorf("Content-Encoding = %q，未要求壓縮時應為空", got)
		}
		if rec.Body.String() != big {
			t.Error("內容與原始檔案不符")
		}
	})

	t.Run("壓不小的檔案與圖片不壓縮", func(t *testing.T) {
		for _, p := range []string{"/assets/tiny-1.js", "/assets/logo-1.png"} {
			req := httptest.NewRequest("GET", p, nil)
			req.Header.Set("Accept-Encoding", "gzip")
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if got := rec.Header().Get("Content-Encoding"); got != "" {
				t.Errorf("%s 的 Content-Encoding = %q，應為空", p, got)
			}
		}
	})
}

// 部署後舊的雜湊檔名會消失。若掉進 SPA fallback 回傳 HTML，拿著舊 index.html 的
// 瀏覽器會把它當成 JS 解析，錯誤訊息看不出真正原因。
func TestMissingAssetIsNotHTML(t *testing.T) {
	h := newTestHandler(t, testDist)
	resp, body := get(t, h, "GET", "/assets/index-gone9.js")

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d，want 404", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); strings.Contains(ct, "text/html") {
		t.Errorf("Content-Type = %q，不該回傳 HTML", ct)
	}
	if strings.Contains(body, "<div id=\"root\">") || strings.Contains(body, "<!doctype html") {
		t.Errorf("body 回傳了 SPA 的 HTML：%.80s", body)
	}
}
