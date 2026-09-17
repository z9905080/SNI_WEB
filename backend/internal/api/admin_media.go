package api

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/z9905080/SNI_WEB/backend/internal/content"
	"github.com/z9905080/SNI_WEB/backend/internal/store"
)

var colorRe = regexp.MustCompile(`^#(?:[0-9A-Fa-f]{3}|[0-9A-Fa-f]{6})$`)

func isHTTPURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

func validImageRef(s string) bool {
	return strings.HasPrefix(s, content.ImageURLPrefix) || isHTTPURL(s)
}

func validLink(s string) bool {
	if s == "" || isHTTPURL(s) {
		return true
	}
	return strings.HasPrefix(s, "/") && !strings.HasPrefix(s, "//")
}

// ---- 輪播圖 ----

type carouselRequest struct {
	Image *string `json:"image"`
	URL   *string `json:"url"`
}

func (req carouselRequest) apply(w http.ResponseWriter, r *http.Request, c *store.Carousel) bool {
	if req.Image != nil {
		c.Image = strings.TrimSpace(*req.Image)
	}
	if req.URL != nil {
		c.URL = strings.TrimSpace(*req.URL)
	}
	if !validImageRef(c.Image) {
		validationError(w, r, "請選擇圖片")
		return false
	}
	if !validLink(c.URL) {
		validationError(w, r, "連結網址需為 http(s):// 開頭或站內路徑")
		return false
	}
	return true
}

func (s *Server) listCarousels(w http.ResponseWriter, r *http.Request) {
	list, err := s.Store.Carousels(r.Context())
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"carousels": list})
}

func (s *Server) createCarousel(w http.ResponseWriter, r *http.Request) {
	var req carouselRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	var c store.Carousel
	if !req.apply(w, r, &c) {
		return
	}
	created, err := s.Store.CreateCarousel(r.Context(), c.Image, c.URL)
	if err != nil {
		s.storeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"carousel": created})
}

func (s *Server) updateCarousel(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var req carouselRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	c, err := s.Store.Carousel(r.Context(), id)
	if err != nil {
		s.storeError(w, r, err)
		return
	}
	if !req.apply(w, r, &c) {
		return
	}
	if err := s.Store.UpdateCarousel(r.Context(), c); err != nil {
		s.storeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"carousel": c})
}

func (s *Server) deleteCarousel(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := s.Store.DeleteCarousel(r.Context(), id); err != nil {
		s.storeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---- 跑馬燈 ----

type marqueeRequest struct {
	Text  *string `json:"text"`
	Color *string `json:"color"`
}

func (req marqueeRequest) apply(w http.ResponseWriter, r *http.Request, m *store.Marquee) bool {
	if req.Text != nil {
		m.Text = *req.Text
	}
	if req.Color != nil {
		m.Color = strings.TrimSpace(*req.Color)
	}
	text, ok := cleanText(m.Text, 1, 500)
	if !ok {
		validationError(w, r, "跑馬燈文字需為 1–500 個字")
		return false
	}
	m.Text = text
	if !colorRe.MatchString(m.Color) {
		validationError(w, r, "顏色格式需為 #RRGGBB")
		return false
	}
	return true
}

func (s *Server) listMarquees(w http.ResponseWriter, r *http.Request) {
	list, err := s.Store.Marquees(r.Context())
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"marquees": list})
}

func (s *Server) createMarquee(w http.ResponseWriter, r *http.Request) {
	var req marqueeRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	var m store.Marquee
	if !req.apply(w, r, &m) {
		return
	}
	created, err := s.Store.CreateMarquee(r.Context(), m.Text, m.Color)
	if err != nil {
		s.storeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"marquee": created})
}

func (s *Server) updateMarquee(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var req marqueeRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	m, err := s.Store.Marquee(r.Context(), id)
	if err != nil {
		s.storeError(w, r, err)
		return
	}
	if !req.apply(w, r, &m) {
		return
	}
	if err := s.Store.UpdateMarquee(r.Context(), m); err != nil {
		s.storeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"marquee": m})
}

func (s *Server) deleteMarquee(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := s.Store.DeleteMarquee(r.Context(), id); err != nil {
		s.storeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---- 網站設定 ----

func (s *Server) getSettings(w http.ResponseWriter, r *http.Request) {
	v, err := s.Store.Settings(r.Context())
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"settings": v})
}

func (s *Server) putSettings(w http.ResponseWriter, r *http.Request) {
	var v store.Settings
	if !decodeJSON(w, r, &v) {
		return
	}
	var ok bool
	if v.WebTitle, ok = cleanText(v.WebTitle, 1, 100); !ok {
		validationError(w, r, "網站標題需為 1–100 個字")
		return
	}
	if v.WebSubTitle, ok = cleanText(v.WebSubTitle, 0, 200); !ok {
		validationError(w, r, "副標題不可超過 200 個字")
		return
	}
	v.FacebookURL = strings.TrimSpace(v.FacebookURL)
	if v.FacebookURL != "" && !isHTTPURL(v.FacebookURL) {
		validationError(w, r, "Facebook 網址需為 http(s):// 開頭")
		return
	}
	if err := s.Store.UpdateSettings(r.Context(), v); err != nil {
		s.internalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"settings": v})
}
