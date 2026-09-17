package site

import (
	"strings"
	"testing"
)

const indexTmpl = `<!doctype html><html><head><meta charset="utf-8"><!--sni:head--></head><body><div id="root"></div></body></html>`

func TestRenderIndex(t *testing.T) {
	out := string(RenderIndex([]byte(indexTmpl), Meta{
		Title:       `練成會｜生長之家 </title><script>alert(1)</script>`,
		Description: `說明 "引號"`,
		Image:       "https://example.org/php/picture/a.jpg",
		URL:         "https://example.org/page/3",
	}, PublicConfig{GAMeasurementID: "G-TEST", CounterScriptURL: "https://c.example/x.php?a=1&b=</script>"}))

	for _, want := range []string{
		`<title>練成會｜生長之家 &lt;/title&gt;&lt;script&gt;alert(1)&lt;/script&gt;</title>`,
		`<meta name="description" content="說明 &#34;引號&#34;">`,
		`<meta property="og:title" content="練成會｜生長之家 &lt;/title&gt;`,
		`<meta property="og:url" content="https://example.org/page/3">`,
		`<meta property="og:type" content="website">`,
		`<meta property="og:image" content="https://example.org/php/picture/a.jpg">`,
		`<meta name="twitter:card" content="summary_large_image">`,
		`<script id="sni-config" type="application/json">{"gaMeasurementId":"G-TEST","counterScriptUrl":"https://c.example/x.php?a=1&b=</script>"}</script>`,
		`<div id="root"></div>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("缺少 %s\n---\n%s", want, out)
		}
	}
	if strings.Contains(out, HeadPlaceholder) || strings.Contains(out, "robots") {
		t.Errorf("不應殘留佔位符或 robots：%s", out)
	}
}

func TestRenderIndexWithoutImageAndNoIndex(t *testing.T) {
	out := string(RenderIndex([]byte(indexTmpl), Meta{Title: "後台", NoIndex: true}, PublicConfig{}))
	if strings.Contains(out, "og:image") || !strings.Contains(out, `content="summary"`) || !strings.Contains(out, `<meta name="robots" content="noindex">`) {
		t.Fatal(out)
	}
}

func TestAbsoluteURL(t *testing.T) {
	base := "https://www.seicho-no-ie.org.tw"
	cases := map[string]string{
		"":                           "",
		"/php/picture/a.jpg":         base + "/php/picture/a.jpg",
		"picture/a.jpg":              base + "/picture/a.jpg",
		"//www.youtube.com/x.jpg":    "https://www.youtube.com/x.jpg",
		"http://old.example/a.jpg":   "http://old.example/a.jpg",
		"javascript:alert(1)":        "",
		"data:image/png;base64,AAAA": "",
	}
	for in, want := range cases {
		if got := absoluteURL(base, in); got != want {
			t.Errorf("absoluteURL(%q) = %q, want %q", in, got, want)
		}
	}
}
