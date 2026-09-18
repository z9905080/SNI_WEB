package db

import (
	"strings"
	"testing"
)

func TestNormalizeDSN(t *testing.T) {
	cases := map[string][]string{
		"root:pw@tcp(127.0.0.1:3306)/sniweb":            {"root:pw@tcp(127.0.0.1:3306)/sniweb?", "parseTime=true", "clientFoundRows=true", "charset=utf8mb4", "timeout=5s", "readTimeout=30s", "writeTimeout=30s"},
		"mysql://root:p%40ss@mysql.zeabur:3306/zeabur":  {"root:p@ss@tcp(mysql.zeabur:3306)/zeabur?", "parseTime=true", "charset=utf8mb4", "timeout=5s", "readTimeout=30s", "writeTimeout=30s"},
		"u:p@tcp(h:1)/d?charset=utf8mb4&timeout=5s":     {"charset=utf8mb4", "timeout=5s", "parseTime=true", "readTimeout=30s", "writeTimeout=30s"},
		"u:p@tcp(h:1)/d?readTimeout=1m&writeTimeout=2m": {"readTimeout=1m0s", "writeTimeout=2m0s", "timeout=5s"},
	}
	for in, wants := range cases {
		got, err := NormalizeDSN(in)
		if err != nil {
			t.Fatalf("%s: %v", in, err)
		}
		for _, w := range wants {
			if !strings.Contains(got, w) {
				t.Errorf("%s → %s，缺少 %s", in, got, w)
			}
		}
		if strings.Count(got, "charset=") != 1 {
			t.Errorf("%s → %s，charset 應只出現一次", in, got)
		}
	}
	if _, err := NormalizeDSN("::::"); err == nil {
		t.Error("無效 DSN 應回錯誤")
	}
}
