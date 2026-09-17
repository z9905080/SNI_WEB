package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/z9905080/SNI_WEB/backend/internal/content"
	"github.com/z9905080/SNI_WEB/backend/internal/storage"
)

const maxUploadFiles = 10

type imageItem struct {
	Name    string    `json:"name"`
	URL     string    `json:"url"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

func (s *Server) listImages(w http.ResponseWriter, r *http.Request) {
	objs, err := s.Storage.List(r.Context())
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	out := make([]imageItem, 0, len(objs))
	for _, o := range objs {
		out = append(out, imageItem{Name: o.Name, URL: content.ImageURL(o.Name), Size: o.Size, ModTime: o.ModTime})
	}
	writeJSON(w, http.StatusOK, map[string]any{"images": out})
}

func uploadReadError(w http.ResponseWriter, r *http.Request) {
	writeError(w, r, http.StatusBadRequest, "invalid_upload", "上傳資料格式錯誤")
}

func (s *Server) uploadImages(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadFiles*(content.MaxImageBytes+64<<10))
	mr, err := r.MultipartReader()
	if err != nil {
		uploadReadError(w, r)
		return
	}

	type upload struct {
		data []byte
		ext  string
	}
	var files []upload
	for {
		part, err := mr.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		var tooBig *http.MaxBytesError
		if errors.As(err, &tooBig) {
			writeError(w, r, http.StatusRequestEntityTooLarge, "too_large", "上傳資料太大")
			return
		}
		if err != nil {
			uploadReadError(w, r)
			return
		}
		if part.FormName() != "files" {
			continue
		}
		if len(files) == maxUploadFiles {
			writeError(w, r, http.StatusBadRequest, "too_many_files", fmt.Sprintf("一次最多上傳 %d 個檔案", maxUploadFiles))
			return
		}
		data, err := io.ReadAll(io.LimitReader(part, content.MaxImageBytes+1))
		if err != nil {
			uploadReadError(w, r)
			return
		}
		if len(data) > content.MaxImageBytes {
			writeError(w, r, http.StatusRequestEntityTooLarge, "file_too_large", fmt.Sprintf("「%s」超過 5MB", part.FileName()))
			return
		}
		ext, ok := content.DetectImageExt(data)
		if !ok {
			writeError(w, r, http.StatusBadRequest, "unsupported_type",
				fmt.Sprintf("「%s」不是支援的圖片格式（jpg、png、gif、webp）", part.FileName()))
			return
		}
		files = append(files, upload{data: data, ext: ext})
	}
	if len(files) == 0 {
		writeError(w, r, http.StatusBadRequest, "no_files", "請選擇要上傳的圖片")
		return
	}

	out := make([]imageItem, 0, len(files))
	for _, f := range files {
		name, err := s.saveImage(r.Context(), f.data, f.ext)
		if err != nil {
			s.internalError(w, r, err)
			return
		}
		out = append(out, imageItem{Name: name, URL: content.ImageURL(name), Size: int64(len(f.data)), ModTime: s.Now()})
	}
	writeJSON(w, http.StatusCreated, map[string]any{"images": out})
}

// saveImage 以時間命名；同名已存在時改用加上亂數的檔名重試。
func (s *Server) saveImage(ctx context.Context, data []byte, ext string) (string, error) {
	now := s.Now()
	name := content.ImageName(now, ext)
	for range 5 {
		err := s.Storage.Put(ctx, name, bytes.NewReader(data), int64(len(data)), content.ImageContentType(ext))
		if !errors.Is(err, storage.ErrExists) {
			return name, err
		}
		name = content.ImageNameWithSuffix(now, ext, content.RandomSuffix())
	}
	return "", errors.New("無法產生不重複的檔名")
}

func imageName(w http.ResponseWriter, r *http.Request) (string, bool) {
	name := r.PathValue("name")
	if !content.ValidImageName(name) {
		writeError(w, r, http.StatusBadRequest, "invalid_name", "檔名格式錯誤")
		return "", false
	}
	return name, true
}

func (s *Server) imageUsages(w http.ResponseWriter, r *http.Request) {
	name, ok := imageName(w, r)
	if !ok {
		return
	}
	u, err := s.Store.ImageUsages(r.Context(), name)
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"usages": u})
}

func (s *Server) deleteImage(w http.ResponseWriter, r *http.Request) {
	name, ok := imageName(w, r)
	if !ok {
		return
	}
	err := s.Storage.Delete(r.Context(), name)
	if errors.Is(err, storage.ErrNotFound) {
		writeError(w, r, http.StatusNotFound, "not_found", "找不到此圖片")
		return
	}
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
