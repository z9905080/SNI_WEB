package site

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/z9905080/SNI_WEB/backend/internal/content"
	"github.com/z9905080/SNI_WEB/backend/internal/storage"
	"github.com/z9905080/SNI_WEB/backend/internal/store"
)

var pagePath = regexp.MustCompile(`^/page/([1-9][0-9]{0,17})$`)

type Handler struct {
	dist    fs.FS
	index   []byte // nil 代表前端尚未建置
	gzipped map[string]gzipAsset
	store   *store.Store
	images  storage.Storage
	baseURL string
	public  PublicConfig
}

type gzipAsset struct {
	body        []byte
	contentType string
}

func New(dist fs.FS, st *store.Store, images storage.Storage, baseURL string, pc PublicConfig) *Handler {
	index, _ := fs.ReadFile(dist, "index.html")
	return &Handler{
		dist: dist, index: index, gzipped: precompress(dist),
		store: st, images: images, baseURL: baseURL, public: pc,
	}
}

// precompress 在啟動時壓好靜態資源。內容在建置時就固定，先壓完可以省下每次請求的 CPU；
// 後台的 JS chunk 有 1MB 以上，未壓縮傳輸在連線不穩時會讓畫面空白數秒。
func precompress(dist fs.FS) map[string]gzipAsset {
	types := map[string]string{
		".js":   "text/javascript; charset=utf-8",
		".css":  "text/css; charset=utf-8",
		".svg":  "image/svg+xml",
		".json": "application/json",
	}
	out := map[string]gzipAsset{}
	_ = fs.WalkDir(dist, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil //nolint:nilerr // 個別檔案讀不到就跳過，不影響服務啟動
		}
		ct, ok := types[strings.ToLower(path.Ext(name))]
		if !ok {
			return nil
		}
		raw, err := fs.ReadFile(dist, name)
		if err != nil {
			return nil //nolint:nilerr
		}
		var buf bytes.Buffer
		zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
		if err != nil {
			return nil //nolint:nilerr
		}
		if _, err := zw.Write(raw); err != nil || zw.Close() != nil {
			return nil //nolint:nilerr
		}
		// 壓不小就不留，直接走原本的檔案服務
		if buf.Len() >= len(raw) {
			return nil
		}
		out[name] = gzipAsset{body: buf.Bytes(), contentType: ct}
		return nil
	})
	return out
}

func acceptsGzip(r *http.Request) bool {
	for enc := range strings.SplitSeq(r.Header.Get("Accept-Encoding"), ",") {
		if name, _, _ := strings.Cut(enc, ";"); strings.TrimSpace(name) == "gzip" {
			return true
		}
	}
	return false
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	p := r.URL.Path
	if name, ok := strings.CutPrefix(p, "/php/picture/"); ok {
		h.images.Serve(w, r, name)
		return
	}
	if name := strings.TrimPrefix(p, "/"); name != "" && name != "index.html" {
		if info, err := fs.Stat(h.dist, name); err == nil && !info.IsDir() {
			cache := "public, max-age=3600"
			if strings.HasPrefix(name, "assets/") {
				cache = "public, max-age=31536000, immutable"
			}
			hdr := w.Header()
			hdr.Set("Cache-Control", cache)
			hdr.Set("Vary", "Accept-Encoding")
			if a, ok := h.gzipped[name]; ok && acceptsGzip(r) {
				hdr.Set("Content-Type", a.contentType)
				hdr.Set("Content-Encoding", "gzip")
				hdr.Set("Content-Length", strconv.Itoa(len(a.body)))
				w.WriteHeader(http.StatusOK)
				if r.Method != http.MethodHead {
					_, _ = w.Write(a.body)
				}
				return
			}
			http.ServeFileFS(w, r, h.dist, name)
			return
		}
	}
	h.serveIndex(w, r)
}

func (h *Handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	if h.index == nil {
		http.Error(w, "前端尚未建置，請執行 make web", http.StatusServiceUnavailable)
		return
	}
	m, status := h.meta(r.Context(), r.URL.Path)
	hdr := w.Header()
	hdr.Set("Content-Type", "text/html; charset=utf-8")
	hdr.Set("Cache-Control", "no-cache")
	if m.NoIndex {
		hdr.Set("X-Robots-Tag", "noindex")
	}
	w.WriteHeader(status)
	_, _ = w.Write(RenderIndex(h.index, m, h.public))
}

// meta 依路徑組出 meta 與 HTTP status；資料庫錯誤只記 log，仍回傳可用的頁面。
func (h *Handler) meta(ctx context.Context, path string) (Meta, int) {
	settings, err := h.store.Settings(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "load settings", "err", err)
	}
	m := Meta{Title: settings.WebTitle, Description: settings.WebSubTitle, URL: h.baseURL + path}
	status := http.StatusOK
	bannerImage := func() string {
		cs, err := h.store.Carousels(ctx)
		if err != nil || len(cs) == 0 {
			return ""
		}
		return cs[0].Image
	}
	describe := func(html string) {
		if d := content.PlainText(html, 120); d != "" {
			m.Description = d
		}
	}

	switch match := pagePath.FindStringSubmatch(path); {
	case path == "/":
		if settings.WebSubTitle != "" {
			m.Title = settings.WebTitle + "｜" + settings.WebSubTitle
		}
		if p, ok, err := h.store.HomePage(ctx); err == nil && ok {
			describe(p.HTML)
		}
		m.Image = bannerImage()
	case match != nil:
		id, _ := strconv.ParseInt(match[1], 10, 64)
		p, err := h.store.Page(ctx, id)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				status = http.StatusNotFound
			} else {
				slog.ErrorContext(ctx, "load page", "err", err, "id", id)
			}
			m.Image = bannerImage()
			break
		}
		m.Title = p.Name + "｜" + settings.WebTitle
		describe(p.HTML)
		if m.Image = content.FirstImageSrc(p.HTML); m.Image == "" {
			m.Image = bannerImage()
		}
	case path == "/admin" || strings.HasPrefix(path, "/admin/"):
		m.NoIndex = true
	default:
		status = http.StatusNotFound
		m.Image = bannerImage()
	}
	m.Image = absoluteURL(h.baseURL, m.Image)
	return m, status
}
