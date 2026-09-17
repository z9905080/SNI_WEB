package site

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
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
	store   *store.Store
	images  storage.Storage
	baseURL string
	public  PublicConfig
}

func New(dist fs.FS, st *store.Store, images storage.Storage, baseURL string, pc PublicConfig) *Handler {
	index, _ := fs.ReadFile(dist, "index.html")
	return &Handler{dist: dist, index: index, store: st, images: images, baseURL: baseURL, public: pc}
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
			w.Header().Set("Cache-Control", cache)
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
	w.Write(RenderIndex(h.index, m, h.public))
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
