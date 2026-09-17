package storage

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDisk(t *testing.T) {
	ctx := context.Background()
	dir := filepath.Join(t.TempDir(), "picture")
	d, err := NewDisk(dir)
	if err != nil {
		t.Fatal(err)
	}

	if err := d.Put(ctx, "2026-01-01_00-00-00.png", strings.NewReader("png-a"), 5, "image/png"); err != nil {
		t.Fatal(err)
	}
	if err := d.Put(ctx, "2026-01-02_00-00-00.jpg", strings.NewReader("jpg-b"), 5, "image/jpeg"); err != nil {
		t.Fatal(err)
	}
	if err := d.Put(ctx, "2026-01-01_00-00-00.png", strings.NewReader("x"), 1, "image/png"); !errors.Is(err, ErrExists) {
		t.Fatalf("重複檔名應回 ErrExists：%v", err)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "2026-01-01_00-00-00.png")); string(b) != "png-a" {
		t.Fatalf("重複上傳不應覆蓋：%q", b)
	}
	if err := d.Put(ctx, "../evil.png", strings.NewReader("x"), 1, "image/png"); !errors.Is(err, ErrInvalidName) {
		t.Fatalf("got %v", err)
	}

	os.WriteFile(filepath.Join(dir, ".DS_Store"), []byte("x"), 0o644)
	os.Mkdir(filepath.Join(dir, "sub"), 0o755)

	objs, err := d.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(objs) != 2 || objs[0].Name != "2026-01-02_00-00-00.jpg" || objs[1].Size != 5 || objs[0].ModTime.IsZero() {
		t.Fatalf("List = %+v", objs)
	}

	w := httptest.NewRecorder()
	d.Serve(w, httptest.NewRequest(http.MethodGet, "/php/picture/2026-01-02_00-00-00.jpg", nil), "2026-01-02_00-00-00.jpg")
	if w.Code != 200 || w.Body.String() != "jpg-b" || w.Header().Get("Content-Type") != "image/jpeg" ||
		w.Header().Get("Cache-Control") != "public, max-age=31536000, immutable" {
		t.Fatalf("Serve: %d %q %v", w.Code, w.Body.String(), w.Header())
	}

	for _, name := range []string{"missing.jpg", "..", ".DS_Store"} {
		w = httptest.NewRecorder()
		d.Serve(w, httptest.NewRequest(http.MethodGet, "/", nil), name)
		if w.Code != 404 {
			t.Errorf("%s: 應 404，得到 %d", name, w.Code)
		}
	}

	if err := d.Delete(ctx, "2026-01-02_00-00-00.jpg"); err != nil {
		t.Fatal(err)
	}
	if err := d.Delete(ctx, "2026-01-02_00-00-00.jpg"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
	if err := d.Delete(ctx, "../x"); !errors.Is(err, ErrInvalidName) {
		t.Fatalf("got %v", err)
	}
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error)       { return 0, errors.New("boom") }
func (failingReader) Seek(int64, int) (int64, error) { return 0, nil }

func TestDiskPutRemovesPartialFile(t *testing.T) {
	dir := t.TempDir()
	d, _ := NewDisk(dir)
	if err := d.Put(context.Background(), "a.png", failingReader{}, 1, "image/png"); err == nil {
		t.Fatal("應回錯誤")
	}
	if _, err := os.Stat(filepath.Join(dir, "a.png")); !os.IsNotExist(err) {
		t.Fatalf("失敗時應刪除殘檔：%v", err)
	}
}
