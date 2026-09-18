package api

import (
	"net/http"

	"github.com/z9905080/SNI_WEB/backend/internal/store"
)

type siteResponse struct {
	store.Settings
	Menu []store.Group `json:"menu"`
}

type pageResponse struct {
	Page      *store.Page      `json:"page"`
	Carousels []store.Carousel `json:"carousels"`
	Marquees  []store.Marquee  `json:"marquees"`
}

func (s *Server) getSite(w http.ResponseWriter, r *http.Request) {
	settings, err := s.Store.Settings(r.Context())
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	menu, err := s.Store.Nav(r.Context())
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	writeJSON(w, http.StatusOK, siteResponse{Settings: settings, Menu: menu})
}

func (s *Server) writePage(w http.ResponseWriter, r *http.Request, page *store.Page) {
	carousels, err := s.Store.Carousels(r.Context())
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	marquees, err := s.Store.Marquees(r.Context())
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	writeJSON(w, http.StatusOK, pageResponse{Page: page, Carousels: carousels, Marquees: marquees})
}

func (s *Server) getHome(w http.ResponseWriter, r *http.Request) {
	page, ok, err := s.Store.HomePage(r.Context())
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	if !ok {
		s.writePage(w, r, nil)
		return
	}
	s.writePage(w, r, &page)
}

func (s *Server) getPage(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	page, err := s.Store.Page(r.Context(), id)
	if err != nil {
		s.storeError(w, r, err)
		return
	}
	s.writePage(w, r, &page)
}
