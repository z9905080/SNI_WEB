// Package site 提供 SPA 靜態檔、meta 注入與 /php/picture 圖片路由。
package site

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"strings"
)

const HeadPlaceholder = "<!--sni:head-->"

type Meta struct {
	Title       string
	Description string
	Image       string
	URL         string
	NoIndex     bool
}

// PublicConfig 是前端需要的公開設定，前端從 #sni-config 讀取。
type PublicConfig struct {
	GAMeasurementID  string `json:"gaMeasurementId"`
	CounterScriptURL string `json:"counterScriptUrl"`
}

func RenderIndex(index []byte, m Meta, pc PublicConfig) []byte {
	var b strings.Builder
	tag := func(attr, key, val string) {
		fmt.Fprintf(&b, "<meta %s=%q content=\"%s\">\n", attr, key, html.EscapeString(val))
	}
	fmt.Fprintf(&b, "<title>%s</title>\n", html.EscapeString(m.Title))
	tag("name", "description", m.Description)
	tag("property", "og:title", m.Title)
	tag("property", "og:description", m.Description)
	tag("property", "og:url", m.URL)
	tag("property", "og:type", "website")
	card := "summary"
	if m.Image != "" {
		tag("property", "og:image", m.Image)
		card = "summary_large_image"
	}
	tag("name", "twitter:card", card)
	if m.NoIndex {
		tag("name", "robots", "noindex")
	}
	// 用 Encoder 並關閉 HTML escaping：json.Marshal 預設會把 <、>、& 轉成 < 等，
	// 但 counterScriptUrl 這類外部設定值需要原樣輸出（見 meta_test.go 的期望值）。
	var cfg bytes.Buffer
	enc := json.NewEncoder(&cfg)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(pc)
	fmt.Fprintf(&b, `<script id="sni-config" type="application/json">%s</script>`, bytes.TrimRight(cfg.Bytes(), "\n"))
	return bytes.Replace(index, []byte(HeadPlaceholder), []byte(b.String()), 1)
}

func absoluteURL(base, ref string) string {
	if ref == "" {
		return ""
	}
	b, err := url.Parse(base + "/")
	if err != nil {
		return ""
	}
	u, err := url.Parse(ref)
	if err != nil {
		return ""
	}
	abs := b.ResolveReference(u)
	if abs.Scheme != "http" && abs.Scheme != "https" {
		return ""
	}
	return abs.String()
}
