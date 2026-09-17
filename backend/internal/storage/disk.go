package storage

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/z9905080/SNI_WEB/backend/internal/content"
)

// Disk 以 os.Root 限制所有檔案操作都在目錄內（檔名驗證之外的第二道防線）。
type Disk struct {
	root *os.Root
}

func NewDisk(dir string) (*Disk, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, err
	}
	return &Disk{root: root}, nil
}

func (d *Disk) Put(_ context.Context, name string, r io.ReadSeeker, _ int64, _ string) error {
	if !content.ValidImageName(name) {
		return ErrInvalidName
	}
	f, err := d.root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, fs.ErrExist) {
		return ErrExists
	}
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		d.root.Remove(name)
	}
	return err
}

func (d *Disk) Delete(_ context.Context, name string) error {
	if !content.ValidImageName(name) {
		return ErrInvalidName
	}
	err := d.root.Remove(name)
	if errors.Is(err, fs.ErrNotExist) {
		return ErrNotFound
	}
	return err
}

func (d *Disk) List(_ context.Context) ([]Object, error) {
	entries, err := fs.ReadDir(d.root.FS(), ".")
	if err != nil {
		return nil, err
	}
	out := make([]Object, 0, len(entries))
	for _, e := range entries {
		if !e.Type().IsRegular() || !content.ValidImageName(e.Name()) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		out = append(out, Object{Name: e.Name(), Size: info.Size(), ModTime: info.ModTime()})
	}
	slices.SortFunc(out, func(a, b Object) int { return strings.Compare(b.Name, a.Name) })
	return out, nil
}

func (d *Disk) Serve(w http.ResponseWriter, r *http.Request, name string) {
	if !content.ValidImageName(name) {
		http.NotFound(w, r)
		return
	}
	f, err := d.root.Open(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", immutableCache)
	http.ServeContent(w, r, name, info.ModTime(), f)
}
