package api

import (
	"net/http"

	"github.com/z9905080/SNI_WEB/backend/internal/store"
)

const nameMessage = "名稱需為 1–50 個字"

type nameRequest struct {
	Name string `json:"name"`
}

type orderRequest struct {
	IDs *[]int64 `json:"ids"`
}

func (s *Server) listGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := s.Store.Nav(r.Context())
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"groups": groups})
}

func (s *Server) createGroup(w http.ResponseWriter, r *http.Request) {
	var req nameRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	name, ok := cleanText(req.Name, 1, 50)
	if !ok {
		validationError(w, r, nameMessage)
		return
	}
	g, err := s.Store.CreateGroup(r.Context(), name)
	if err != nil {
		s.storeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"group": g})
}

func (s *Server) renameGroup(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var req nameRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	name, ok := cleanText(req.Name, 1, 50)
	if !ok {
		validationError(w, r, nameMessage)
		return
	}
	if err := s.Store.RenameGroup(r.Context(), id, name); err != nil {
		s.storeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"group": map[string]any{"id": id, "name": name}})
}

func (s *Server) deleteGroup(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := s.Store.DeleteGroup(r.Context(), id); err != nil {
		s.storeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func decodeOrder(w http.ResponseWriter, r *http.Request) ([]int64, bool) {
	var req orderRequest
	if !decodeJSON(w, r, &req) {
		return nil, false
	}
	if req.IDs == nil {
		validationError(w, r, "缺少排序資料")
		return nil, false
	}
	return *req.IDs, true
}

func (s *Server) setGroupOrder(w http.ResponseWriter, r *http.Request) {
	ids, ok := decodeOrder(w, r)
	if !ok {
		return
	}
	if err := s.Store.SetGroupOrder(r.Context(), ids); err != nil {
		s.storeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) setPageOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	ids, ok := decodeOrder(w, r)
	if !ok {
		return
	}
	if err := s.Store.SetPageOrder(r.Context(), id, ids); err != nil {
		s.storeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type pageRequest struct {
	Name    *string `json:"name"`
	GroupID *int64  `json:"group_id"`
	HTML    *string `json:"html"`
}

// apply 把請求中有提供的欄位套用到 p；驗證失敗時已寫出回應。
func (req pageRequest) apply(w http.ResponseWriter, r *http.Request, p *store.Page) bool {
	if req.Name != nil {
		name, ok := cleanText(*req.Name, 1, 50)
		if !ok {
			validationError(w, r, "頁面"+nameMessage)
			return false
		}
		p.Name = name
	}
	if req.GroupID != nil {
		if *req.GroupID <= 0 {
			validationError(w, r, "請選擇所屬頁籤")
			return false
		}
		p.GroupID = *req.GroupID
	}
	if req.HTML != nil {
		p.HTML = *req.HTML
	}
	return true
}

func (s *Server) createPage(w http.ResponseWriter, r *http.Request) {
	var req pageRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.Name == nil || req.GroupID == nil {
		validationError(w, r, "請填寫頁面名稱與所屬頁籤")
		return
	}
	var p store.Page
	if !req.apply(w, r, &p) {
		return
	}
	created, err := s.Store.CreatePage(r.Context(), p)
	if err != nil {
		s.storeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"page": created})
}

func (s *Server) getAdminPage(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	p, err := s.Store.Page(r.Context(), id)
	if err != nil {
		s.storeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"page": p})
}

func (s *Server) updatePage(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var req pageRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	p, err := s.Store.Page(r.Context(), id)
	if err != nil {
		s.storeError(w, r, err)
		return
	}
	if !req.apply(w, r, &p) {
		return
	}
	updated, err := s.Store.UpdatePage(r.Context(), p)
	if err != nil {
		s.storeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"page": updated})
}

func (s *Server) deletePage(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := s.Store.DeletePage(r.Context(), id); err != nil {
		s.storeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
