package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/z9905080/SNI_WEB/backend/internal/auth"
)

func (s *Server) csrf(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
		default:
			if !auth.OriginAllowed(r, s.AllowedOrigins) {
				writeError(w, r, http.StatusForbidden, "forbidden_origin", "來源不被允許")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(auth.CookieName)
		if err != nil {
			writeError(w, r, http.StatusUnauthorized, "unauthorized", "請先登入")
			return
		}
		u, renewed, err := s.Sessions.Validate(r.Context(), c.Value)
		if errors.Is(err, auth.ErrInvalidSession) {
			auth.ClearCookie(w, s.CookieSecure)
			writeError(w, r, http.StatusUnauthorized, "unauthorized", "登入已逾時，請重新登入")
			return
		}
		if err != nil {
			s.internalError(w, r, err)
			return
		}
		if !renewed.IsZero() {
			auth.SetCookie(w, c.Value, renewed, s.CookieSecure)
		}
		next.ServeHTTP(w, r.WithContext(auth.WithUser(r.Context(), u)))
	})
}

func (s *Server) admin(mux *http.ServeMux, pattern string, h http.HandlerFunc) {
	mux.Handle(pattern, s.csrf(s.requireAuth(h)))
}

type loginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	account := strings.TrimSpace(req.Account)
	if account == "" || req.Password == "" {
		validationError(w, r, "請輸入帳號與密碼")
		return
	}

	ip := auth.ClientIP(r)
	if s.ipLimiter.Blocked(ip) || s.accountLimiter.Blocked(account) {
		writeError(w, r, http.StatusTooManyRequests, "too_many_requests", "嘗試次數過多，請 15 分鐘後再試")
		return
	}
	s.ipLimiter.Add(ip)

	u, raw, exp, err := s.Sessions.Login(r.Context(), account, req.Password)
	if errors.Is(err, auth.ErrInvalidCredentials) {
		s.accountLimiter.Add(account)
		writeError(w, r, http.StatusUnauthorized, "invalid_credentials", "帳號或密碼錯誤")
		return
	}
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	s.accountLimiter.Reset(account)
	auth.SetCookie(w, raw, exp, s.CookieSecure)
	writeJSON(w, http.StatusOK, map[string]auth.User{"user": u})
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.UserFrom(r.Context())
	if err := s.Sessions.Logout(r.Context(), u.ID); err != nil {
		s.internalError(w, r, err)
		return
	}
	auth.ClearCookie(w, s.CookieSecure)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.UserFrom(r.Context())
	writeJSON(w, http.StatusOK, map[string]auth.User{"user": u})
}
