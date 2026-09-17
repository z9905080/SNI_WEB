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
		if err := s.Ping(r.Context()); err != nil {
			writeError(w, r, http.StatusServiceUnavailable, "not_ready", "資料庫無法連線")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /api/v1/site", s.getSite)
	mux.HandleFunc("GET /api/v1/home", s.getHome)
	mux.HandleFunc("GET /api/v1/pages/{id}", s.getPage)

	// admin routes（後續 task 加入）

	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, r, http.StatusNotFound, "not_found", "找不到此 API")
	})

	return mux
}
