package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/z9905080/SNI_WEB/backend/internal/content"
)

type ctxKey int

const requestIDKey ctxKey = iota

func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	if s.status == 0 {
		s.status = code
	}
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	if s.status == 0 {
		s.status = http.StatusOK
	}
	return s.ResponseWriter.Write(b)
}

func (s *statusRecorder) Unwrap() http.ResponseWriter { return s.ResponseWriter }

// Wrap 套用所有請求共用的 middleware。
func Wrap(h http.Handler, csp string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := make([]byte, 8)
		rand.Read(b)
		id := hex.EncodeToString(b)
		r = r.WithContext(context.WithValue(r.Context(), requestIDKey, id))

		hdr := w.Header()
		hdr.Set("X-Request-ID", id)
		hdr.Set("X-Content-Type-Options", "nosniff")
		hdr.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		hdr.Set("Content-Security-Policy", csp)

		rec := &statusRecorder{ResponseWriter: w}
		start := time.Now()
		defer func() {
			if v := recover(); v != nil {
				if err, ok := v.(error); ok && errors.Is(err, http.ErrAbortHandler) {
					panic(v)
				}
				slog.ErrorContext(r.Context(), "panic", "err", fmt.Sprint(v), "request_id", id)
				if rec.status == 0 {
					writeError(rec, r, http.StatusInternalServerError, "internal", "系統發生錯誤，請稍後再試")
				}
			}
			slog.InfoContext(r.Context(), "http",
				"method", r.Method, "path", r.URL.Path, "status", rec.status,
				"duration_ms", time.Since(start).Milliseconds(), "request_id", id)
		}()
		h.ServeHTTP(rec, r)
	})
}

// BuildCSP 產生 Content-Security-Policy；計數器 script 會在 iframe srcdoc 中執行，其來源需加入 script-src。
func BuildCSP(counterScriptURL string) string {
	scripts := []string{"'self'", "https://www.googletagmanager.com"}
	if u, err := url.Parse(counterScriptURL); err == nil && u.Scheme != "" && u.Host != "" {
		scripts = append(scripts, u.Scheme+"://"+u.Host)
	}
	frames := []string{"'self'"}
	for _, h := range content.IframeHosts {
		frames = append(frames, "https://"+h)
	}
	return strings.Join([]string{
		"default-src 'self'",
		"script-src " + strings.Join(scripts, " "),
		"connect-src 'self' https://*.google-analytics.com https://*.analytics.google.com https://www.googletagmanager.com",
		"img-src 'self' data: blob: https: http:",
		"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
		"font-src 'self' data: https://fonts.gstatic.com",
		"frame-src " + strings.Join(frames, " "),
		"object-src 'none'",
		"base-uri 'self'",
		"frame-ancestors 'self'",
		"upgrade-insecure-requests",
	}, "; ")
}
