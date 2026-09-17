package api

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/z9905080/SNI_WEB/backend/internal/store"
)

var pngBytes = append([]byte("\x89PNG\r\n\x1a\n"), make([]byte, 64)...)

type namedFile struct {
	name string
	data []byte
}

func (h *harness) upload(files ...namedFile) *http.Response {
	h.t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	for _, f := range files {
		fw, err := mw.CreateFormFile("files", f.name)
		if err != nil {
			h.t.Fatal(err)
		}
		fw.Write(f.data)
	}
	mw.Close()
	req, _ := http.NewRequest("POST", h.srv.URL+"/api/v1/admin/images", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Origin", testOrigin)
	return h.send(req)
}

type imagesBody struct {
	Images []struct {
		Name    string    `json:"name"`
		URL     string    `json:"url"`
		Size    int64     `json:"size"`
		ModTime time.Time `json:"mod_time"`
	} `json:"images"`
}

func TestUploadListDeleteImages(t *testing.T) {
	h := newHarness(t)
	h.login()

	resp := h.upload(namedFile{"a.png", pngBytes}, namedFile{"b.png", pngBytes})
	expectStatus(t, resp, 201)
	up := decode[imagesBody](t, resp)
	if len(up.Images) != 2 {
		t.Fatalf("%+v", up)
	}
	first, second := up.Images[0].Name, up.Images[1].Name
	if !regexp.MustCompile(`^\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}\.png$`).MatchString(first) {
		t.Errorf("檔名格式錯誤：%s", first)
	}
	if !regexp.MustCompile(`^\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}-[a-z2-7]{6}\.png$`).MatchString(second) {
		t.Errorf("同一秒的第二個檔案應加亂數：%s", second)
	}
	if up.Images[0].URL != "/php/picture/"+first || up.Images[0].Size != int64(len(pngBytes)) {
		t.Errorf("%+v", up.Images[0])
	}

	list := decode[imagesBody](t, h.do("GET", "/api/v1/admin/images", nil))
	if len(list.Images) != 2 {
		t.Fatalf("%+v", list)
	}

	expectStatus(t, h.do("DELETE", "/api/v1/admin/images/"+first, nil), 204)
	expectError(t, h.do("DELETE", "/api/v1/admin/images/"+first, nil), 404, "not_found")
	expectError(t, h.do("DELETE", "/api/v1/admin/images/.hidden", nil), 400, "invalid_name")
}

func TestUploadRejections(t *testing.T) {
	h := newHarness(t)
	h.login()

	expectError(t, h.upload(), 400, "no_files")
	expectError(t, h.upload(namedFile{"x.svg", []byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`)}), 400, "unsupported_type")
	big := append(append([]byte{}, pngBytes...), make([]byte, 5<<20)...)
	e := expectError(t, h.upload(namedFile{"big.png", big}), 413, "file_too_large")
	if !strings.Contains(e.Error.Message, "big.png") {
		t.Errorf("訊息應包含檔名：%s", e.Error.Message)
	}

	var many []namedFile
	for range 11 {
		many = append(many, namedFile{"a.png", pngBytes})
	}
	expectError(t, h.upload(many...), 400, "too_many_files")

	// 整批拒絕：前面合格的檔案也不能寫入
	expectError(t, h.upload(namedFile{"ok.png", pngBytes}, namedFile{"bad.txt", []byte("hello")}), 400, "unsupported_type")
	if list := decode[imagesBody](t, h.do("GET", "/api/v1/admin/images", nil)); len(list.Images) != 0 {
		t.Fatalf("不應寫入任何檔案：%+v", list)
	}

	req, _ := http.NewRequest("POST", h.srv.URL+"/api/v1/admin/images", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", testOrigin)
	expectError(t, h.send(req), 400, "invalid_upload")
}

func TestImageUsagesEndpoint(t *testing.T) {
	h := newHarness(t,
		`INSERT INTO page_content (id, page_group_id, page_name, html_context) VALUES (1, 1, '用到', '<img src="/php/picture/a.jpg">')`)
	h.login()
	b := decode[struct {
		Usages store.ImageUsages `json:"usages"`
	}](t, h.do("GET", "/api/v1/admin/images/a.jpg/usages", nil))
	if len(b.Usages.Pages) != 1 || b.Usages.Carousels == nil {
		t.Fatalf("%+v", b)
	}
	expectError(t, h.do("GET", "/api/v1/admin/images/a%20b.jpg/usages", nil), 400, "invalid_name")
}
