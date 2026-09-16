package content

import (
	"crypto/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	_ "time/tzdata" // distroless 映像沒有時區資料
)

const (
	MaxImageBytes  = 5 << 20
	ImageURLPrefix = "/php/picture/"
)

var (
	taipei      = mustLocation("Asia/Taipei")
	imageExts   = map[string]string{"image/jpeg": "jpg", "image/png": "png", "image/gif": "gif", "image/webp": "webp"}
	imageTypes  = map[string]string{"jpg": "image/jpeg", "png": "image/png", "gif": "image/gif", "webp": "image/webp"}
	validNameRe = regexp.MustCompile(`^[A-Za-z0-9._-]{1,255}$`)
)

func mustLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

// DetectImageExt 依檔案開頭內容判斷圖片類型，不信任客戶端提供的 MIME。
func DetectImageExt(head []byte) (string, bool) {
	if len(head) == 0 {
		return "", false
	}
	ext, ok := imageExts[http.DetectContentType(head)]
	return ext, ok
}

func ImageContentType(ext string) string {
	if t, ok := imageTypes[ext]; ok {
		return t
	}
	return "application/octet-stream"
}

func ImageName(now time.Time, ext string) string {
	return now.In(taipei).Format("2006-01-02_15-04-05") + "." + ext
}

func ImageNameWithSuffix(now time.Time, ext, suffix string) string {
	return now.In(taipei).Format("2006-01-02_15-04-05") + "-" + suffix + "." + ext
}

// RandomSuffix 回傳 6 碼 [a-z2-7]（base32 小寫）。
func RandomSuffix() string {
	return strings.ToLower(rand.Text()[:6])
}

func ValidImageName(name string) bool {
	return validNameRe.MatchString(name) && !strings.HasPrefix(name, ".")
}

func ImageURL(name string) string {
	return ImageURLPrefix + name
}
