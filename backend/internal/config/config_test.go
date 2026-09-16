package config

import (
	"slices"
	"strings"
	"testing"
)

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestLoadDefaults(t *testing.T) {
	c, err := Load(env(map[string]string{
		"DATABASE_URL":    "u:p@tcp(db:3306)/x",
		"PUBLIC_BASE_URL": "https://example.org/",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if c.Port != "8080" || !c.CookieSecure || c.StorageDriver != "disk" || c.StorageDiskDir != "/data/picture" {
		t.Fatalf("預設值錯誤：%+v", c)
	}
	if c.PublicBaseURL != "https://example.org" {
		t.Fatalf("PublicBaseURL 應去掉結尾斜線，得到 %q", c.PublicBaseURL)
	}
	if len(c.DevOrigins) != 0 {
		t.Fatalf("DevOrigins 應為空，得到 %v", c.DevOrigins)
	}
}

func TestLoadMissingRequired(t *testing.T) {
	_, err := Load(env(map[string]string{}))
	if err == nil {
		t.Fatal("應回傳錯誤")
	}
	for _, k := range []string{"DATABASE_URL", "PUBLIC_BASE_URL"} {
		if !strings.Contains(err.Error(), k) {
			t.Errorf("錯誤訊息應包含 %s：%v", k, err)
		}
	}
}

func TestLoadS3RequiresAllFields(t *testing.T) {
	_, err := Load(env(map[string]string{
		"DATABASE_URL":    "x",
		"PUBLIC_BASE_URL": "https://a",
		"STORAGE_DRIVER":  "s3",
		"S3_BUCKET":       "b",
	}))
	if err == nil {
		t.Fatal("應回傳錯誤")
	}
	for _, k := range []string{"S3_ENDPOINT", "S3_REGION", "S3_ACCESS_KEY", "S3_SECRET_KEY", "S3_PUBLIC_BASE_URL"} {
		if !strings.Contains(err.Error(), k) {
			t.Errorf("錯誤訊息應包含 %s：%v", k, err)
		}
	}
	if strings.Contains(err.Error(), "S3_BUCKET") {
		t.Errorf("S3_BUCKET 已提供，不應列出：%v", err)
	}
}

func TestLoadInvalidValues(t *testing.T) {
	_, err := Load(env(map[string]string{
		"DATABASE_URL":    "x",
		"PUBLIC_BASE_URL": "https://a",
		"STORAGE_DRIVER":  "ftp",
		"COOKIE_SECURE":   "maybe",
	}))
	if err == nil || !strings.Contains(err.Error(), "STORAGE_DRIVER") || !strings.Contains(err.Error(), "COOKIE_SECURE") {
		t.Fatalf("應同時指出 STORAGE_DRIVER 與 COOKIE_SECURE：%v", err)
	}
}

func TestAllowedOrigins(t *testing.T) {
	c, err := Load(env(map[string]string{
		"DATABASE_URL":    "x",
		"PUBLIC_BASE_URL": "https://a.org",
		"COOKIE_SECURE":   "false",
		"DEV_ORIGINS":     " http://localhost:5173/ , ,http://127.0.0.1:5173",
	}))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"https://a.org", "http://localhost:5173", "http://127.0.0.1:5173"}
	if got := c.AllowedOrigins(); !slices.Equal(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
	if c.CookieSecure {
		t.Fatal("COOKIE_SECURE=false 應為 false")
	}
}
