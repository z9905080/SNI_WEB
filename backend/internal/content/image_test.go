package content

import (
	"regexp"
	"testing"
	"time"
)

func TestDetectImageExt(t *testing.T) {
	cases := []struct {
		name string
		head []byte
		ext  string
		ok   bool
	}{
		{"jpeg", []byte("\xFF\xD8\xFF\xE0\x00\x10JFIF"), "jpg", true},
		{"png", []byte("\x89PNG\r\n\x1a\n\x00\x00"), "png", true},
		{"gif", []byte("GIF89a\x01\x00"), "gif", true},
		{"webp", []byte("RIFF\x24\x00\x00\x00WEBPVP8 "), "webp", true},
		{"html 偽裝", []byte("<html><script>"), "", false},
		{"svg", []byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`), "", false},
		{"空", nil, "", false},
	}
	for _, c := range cases {
		ext, ok := DetectImageExt(c.head)
		if ext != c.ext || ok != c.ok {
			t.Errorf("%s: got (%q,%v)", c.name, ext, ok)
		}
	}
}

func TestImageContentType(t *testing.T) {
	for ext, want := range map[string]string{"jpg": "image/jpeg", "png": "image/png", "gif": "image/gif", "webp": "image/webp", "exe": "application/octet-stream"} {
		if got := ImageContentType(ext); got != want {
			t.Errorf("%s → %s", ext, got)
		}
	}
}

func TestImageNames(t *testing.T) {
	now := time.Date(2026, 9, 16, 1, 2, 3, 0, time.UTC) // 台北時間 09:02:03
	if got := ImageName(now, "jpg"); got != "2026-09-16_09-02-03.jpg" {
		t.Errorf("ImageName = %q", got)
	}
	if got := ImageNameWithSuffix(now, "png", "a1b2c3"); got != "2026-09-16_09-02-03-a1b2c3.png" {
		t.Errorf("ImageNameWithSuffix = %q", got)
	}
	re := regexp.MustCompile(`^[a-z0-9]{6}$`)
	seen := map[string]bool{}
	for range 20 {
		s := RandomSuffix()
		if !re.MatchString(s) {
			t.Fatalf("RandomSuffix = %q", s)
		}
		seen[s] = true
	}
	if len(seen) < 15 {
		t.Fatalf("RandomSuffix 重複太多：%v", seen)
	}
}

func TestValidImageName(t *testing.T) {
	for _, ok := range []string{"2019-12-25_08-27-33.jpg", "a.B-c_d.png"} {
		if !ValidImageName(ok) {
			t.Errorf("%q 應合法", ok)
		}
	}
	for _, bad := range []string{"", ".", "..", ".hidden", "../x.jpg", "a/b.jpg", `a\b.jpg`, "中文.jpg", "a b.jpg", "a%2e.jpg", string(make([]byte, 256))} {
		if ValidImageName(bad) {
			t.Errorf("%q 應不合法", bad)
		}
	}
}

func TestImageURL(t *testing.T) {
	if got := ImageURL("a.jpg"); got != "/php/picture/a.jpg" {
		t.Fatal(got)
	}
}
