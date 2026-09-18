package auth

import (
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

const CookieName = "sni_session"

func SetCookie(w http.ResponseWriter, raw string, expires time.Time, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    raw,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// OriginAllowed 檢查 Origin（沒有時改用 Referer）是否在允許清單中。
func OriginAllowed(r *http.Request, allowed []string) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		u, err := url.Parse(r.Header.Get("Referer"))
		if err != nil || u.Scheme == "" || u.Host == "" {
			return false
		}
		origin = u.Scheme + "://" + u.Host
	}
	origin = strings.ToLower(strings.TrimRight(origin, "/"))
	return slices.ContainsFunc(allowed, func(a string) bool { return strings.ToLower(a) == origin })
}

// ClientIP 取反向代理加入的 X-Forwarded-For 最右側一筆，沒有時使用連線來源。
// security: 只有服務前面一定有反向代理（Zeabur）時最右側一筆才可信；直連時可偽造，
// 但帳號層級的限流仍然有效。
func ClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if ip := strings.TrimSpace(parts[len(parts)-1]); ip != "" {
			return ip
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
