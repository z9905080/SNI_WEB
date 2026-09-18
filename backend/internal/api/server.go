// Package api 實作 /api/v1 的 HTTP handlers。
package api

import (
	"context"
	"net/http"
	"time"

	"github.com/z9905080/SNI_WEB/backend/internal/auth"
	"github.com/z9905080/SNI_WEB/backend/internal/storage"
	"github.com/z9905080/SNI_WEB/backend/internal/store"
)

type Deps struct {
	Store          *store.Store
	Sessions       *auth.Sessions
	Storage        storage.Storage
	CookieSecure   bool
	AllowedOrigins []string
	Now            func() time.Time
	Ping           func(context.Context) error
}

type Server struct {
	Deps
	ipLimiter      *auth.Limiter
	accountLimiter *auth.Limiter
}

func New(d Deps) *Server {
	if d.Now == nil {
		d.Now = time.Now
	}
	return &Server{
		Deps:           d,
		ipLimiter:      auth.NewLimiter(20, 15*time.Minute, d.Now),
		accountLimiter: auth.NewLimiter(10, 15*time.Minute, d.Now),
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := s.Ping(ctx); err != nil {
			writeError(w, r, http.StatusServiceUnavailable, "not_ready", "資料庫無法連線")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /api/v1/site", s.getSite)
	mux.HandleFunc("GET /api/v1/home", s.getHome)
	mux.HandleFunc("GET /api/v1/pages/{id}", s.getPage)

	mux.Handle("POST /api/v1/admin/auth/login", s.csrf(http.HandlerFunc(s.login)))
	s.admin(mux, "POST /api/v1/admin/auth/logout", s.logout)
	s.admin(mux, "GET /api/v1/admin/auth/me", s.me)
	s.admin(mux, "GET /api/v1/admin/groups", s.listGroups)
	s.admin(mux, "POST /api/v1/admin/groups", s.createGroup)
	s.admin(mux, "PATCH /api/v1/admin/groups/{id}", s.renameGroup)
	s.admin(mux, "DELETE /api/v1/admin/groups/{id}", s.deleteGroup)
	s.admin(mux, "PUT /api/v1/admin/groups/order", s.setGroupOrder)
	s.admin(mux, "PUT /api/v1/admin/groups/{id}/pages/order", s.setPageOrder)
	s.admin(mux, "POST /api/v1/admin/pages", s.createPage)
	s.admin(mux, "GET /api/v1/admin/pages/{id}", s.getAdminPage)
	s.admin(mux, "PATCH /api/v1/admin/pages/{id}", s.updatePage)
	s.admin(mux, "DELETE /api/v1/admin/pages/{id}", s.deletePage)
	s.admin(mux, "GET /api/v1/admin/carousels", s.listCarousels)
	s.admin(mux, "POST /api/v1/admin/carousels", s.createCarousel)
	s.admin(mux, "PATCH /api/v1/admin/carousels/{id}", s.updateCarousel)
	s.admin(mux, "DELETE /api/v1/admin/carousels/{id}", s.deleteCarousel)
	s.admin(mux, "GET /api/v1/admin/marquees", s.listMarquees)
	s.admin(mux, "POST /api/v1/admin/marquees", s.createMarquee)
	s.admin(mux, "PATCH /api/v1/admin/marquees/{id}", s.updateMarquee)
	s.admin(mux, "DELETE /api/v1/admin/marquees/{id}", s.deleteMarquee)
	s.admin(mux, "GET /api/v1/admin/settings", s.getSettings)
	s.admin(mux, "PUT /api/v1/admin/settings", s.putSettings)
	s.admin(mux, "GET /api/v1/admin/images", s.listImages)
	s.admin(mux, "POST /api/v1/admin/images", s.uploadImages)
	s.admin(mux, "GET /api/v1/admin/images/{name}/usages", s.imageUsages)
	s.admin(mux, "DELETE /api/v1/admin/images/{name}", s.deleteImage)

	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, r, http.StatusNotFound, "not_found", "找不到此 API")
	})

	return mux
}
