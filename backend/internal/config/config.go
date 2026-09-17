// Package config 從環境變數載入並驗證服務設定。
package config

import (
	"errors"
	"strconv"
	"strings"
)

type S3 struct {
	Endpoint      string
	Region        string
	Bucket        string
	AccessKey     string
	SecretKey     string
	PublicBaseURL string
}

type Config struct {
	Port             string
	DatabaseURL      string
	PublicBaseURL    string
	CookieSecure     bool
	DevOrigins       []string
	StorageDriver    string
	StorageDiskDir   string
	S3               S3
	GAMeasurementID  string
	CounterScriptURL string
}

// Load 讀取環境變數；缺少必填值或值不合法時，錯誤訊息會一次列出所有問題。
func Load(getenv func(string) string) (Config, error) {
	get := func(key, def string) string {
		if v := strings.TrimSpace(getenv(key)); v != "" {
			return v
		}
		return def
	}

	c := Config{
		Port:           get("PORT", "8080"),
		DatabaseURL:    get("DATABASE_URL", ""),
		PublicBaseURL:  strings.TrimRight(get("PUBLIC_BASE_URL", ""), "/"),
		StorageDriver:  get("STORAGE_DRIVER", "disk"),
		StorageDiskDir: get("STORAGE_DISK_DIR", "/data/picture"),
		S3: S3{
			Endpoint:      get("S3_ENDPOINT", ""),
			Region:        get("S3_REGION", ""),
			Bucket:        get("S3_BUCKET", ""),
			AccessKey:     get("S3_ACCESS_KEY", ""),
			SecretKey:     get("S3_SECRET_KEY", ""),
			PublicBaseURL: strings.TrimRight(get("S3_PUBLIC_BASE_URL", ""), "/"),
		},
		GAMeasurementID:  get("GA_MEASUREMENT_ID", ""),
		CounterScriptURL: get("COUNTER_SCRIPT_URL", ""),
	}

	var problems, missing []string

	secure, err := strconv.ParseBool(get("COOKIE_SECURE", "true"))
	if err != nil {
		problems = append(problems, "COOKIE_SECURE 必須是 true 或 false")
	}
	c.CookieSecure = secure

	for o := range strings.SplitSeq(getenv("DEV_ORIGINS"), ",") {
		if o = strings.TrimRight(strings.TrimSpace(o), "/"); o != "" {
			c.DevOrigins = append(c.DevOrigins, o)
		}
	}

	if c.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if c.PublicBaseURL == "" {
		missing = append(missing, "PUBLIC_BASE_URL")
	}

	switch c.StorageDriver {
	case "disk":
	case "s3":
		for _, kv := range []struct{ key, val string }{
			{"S3_ENDPOINT", c.S3.Endpoint},
			{"S3_REGION", c.S3.Region},
			{"S3_BUCKET", c.S3.Bucket},
			{"S3_ACCESS_KEY", c.S3.AccessKey},
			{"S3_SECRET_KEY", c.S3.SecretKey},
			{"S3_PUBLIC_BASE_URL", c.S3.PublicBaseURL},
		} {
			if kv.val == "" {
				missing = append(missing, kv.key)
			}
		}
	default:
		problems = append(problems, "STORAGE_DRIVER 必須是 disk 或 s3")
	}

	if len(missing) > 0 {
		problems = append(problems, "缺少環境變數："+strings.Join(missing, ", "))
	}
	if len(problems) > 0 {
		return Config{}, errors.New(strings.Join(problems, "；"))
	}
	return c, nil
}

// AllowedOrigins 回傳後台非 GET 請求允許的 Origin。
func (c Config) AllowedOrigins() []string {
	return append([]string{c.PublicBaseURL}, c.DevOrigins...)
}
