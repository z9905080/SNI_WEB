package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/z9905080/SNI_WEB/backend/internal/store"
)

const maxJSONBytes = 1 << 20

type apiError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

func writeError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	e := apiError{Code: code, Message: message}
	if status >= 500 {
		e.RequestID = RequestID(r.Context())
	}
	writeJSON(w, status, map[string]apiError{"error": e})
}

func (s *Server) internalError(w http.ResponseWriter, r *http.Request, err error) {
	slog.ErrorContext(r.Context(), "internal error", "err", err, "request_id", RequestID(r.Context()),
		"method", r.Method, "path", r.URL.Path)
	writeError(w, r, http.StatusInternalServerError, "internal", "系統發生錯誤，請稍後再試")
}

func (s *Server) storeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, r, http.StatusNotFound, "not_found", "找不到資料")
	case errors.Is(err, store.ErrGroupNotEmpty):
		writeError(w, r, http.StatusConflict, "group_not_empty", "此頁籤仍有頁面，請先移除或刪除頁面")
	case errors.Is(err, store.ErrGroupProtected):
		writeError(w, r, http.StatusConflict, "group_protected", "首頁頁籤不可刪除")
	case errors.Is(err, store.ErrOrderMismatch):
		writeError(w, r, http.StatusBadRequest, "invalid_order", "排序資料與現有項目不一致，請重新整理後再試")
	case errors.Is(err, store.ErrInvalidGroup):
		writeError(w, r, http.StatusBadRequest, "invalid_group", "所屬頁籤不存在")
	case errors.Is(err, store.ErrContentTooLong):
		writeError(w, r, http.StatusBadRequest, "content_too_long", "內容過長（上限約 65KB），請分成多個頁面")
	default:
		s.internalError(w, r, err)
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBytes)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		var tooBig *http.MaxBytesError
		if errors.As(err, &tooBig) {
			writeError(w, r, http.StatusRequestEntityTooLarge, "too_large", "資料太大")
		} else {
			writeError(w, r, http.StatusBadRequest, "invalid_json", "資料格式錯誤")
		}
		return false
	}
	return true
}

func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, r, http.StatusBadRequest, "invalid_id", "編號格式錯誤")
		return 0, false
	}
	return id, true
}

func validationError(w http.ResponseWriter, r *http.Request, message string) {
	writeError(w, r, http.StatusBadRequest, "validation", message)
}

// cleanText 去除前後空白並檢查長度（以字元計）。
func cleanText(s string, min, max int) (string, bool) {
	s = strings.TrimSpace(s)
	n := utf8.RuneCountInString(s)
	return s, n >= min && n <= max
}
