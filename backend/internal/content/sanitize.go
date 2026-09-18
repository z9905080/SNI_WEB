package content

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

// IframeHosts 為允許嵌入的 iframe 網域；前端 web/src/admin/editor/iframeHosts.ts 需同步。
var IframeHosts = []string{
	"www.youtube.com",
	"youtube.com",
	"www.youtube-nocookie.com",
	"players.brightcove.net",
	"drive.google.com",
	"www.facebook.com",
}

var policy = newPolicy()

func iframeSrcPattern() *regexp.Regexp {
	quoted := make([]string, len(IframeHosts))
	for i, h := range IframeHosts {
		quoted[i] = regexp.QuoteMeta(h)
	}
	return regexp.MustCompile(`^(?i)(?:https?:)?//(?:` + strings.Join(quoted, "|") + `)(?:[/?#]|$)`)
}

func newPolicy() *bluemonday.Policy {
	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()
	p.AllowStandardAttributes()
	p.AllowImages() // 內部會呼叫 AllowStandardURLs()，因此 RequireNoFollowOnLinks(false) 需放在其後才生效
	p.AllowLists()
	p.AllowTables()
	p.RequireNoFollowOnLinks(false)
	p.AllowElements(
		"p", "br", "hr", "h1", "h2", "h3", "h4", "h5", "h6",
		"strong", "b", "em", "i", "u", "s", "strike", "del", "ins", "sub", "sup", "small", "big", "mark",
		"span", "div", "blockquote", "pre", "code", "figure", "figcaption", "center", "font",
		"section", "article", "header", "footer",
		"table", "thead", "tbody", "tfoot", "tr", "td", "th",
	)
	p.AllowAttrs("href").OnElements("a")
	p.AllowAttrs("target").Matching(regexp.MustCompile(`^_(?:blank|self)$`)).OnElements("a")
	p.AllowAttrs("rel").OnElements("a")
	// 舊內容大量使用行內樣式；bluemonday 在沒有設定 style 規則時會原樣保留 style 屬性
	p.AllowAttrs("style", "class").Globally()
	p.AllowAttrs("border", "cellpadding", "cellspacing", "bgcolor").OnElements("table", "td", "th", "tr")
	p.AllowAttrs("color", "face", "size").OnElements("font")
	p.AllowAttrs("src").Matching(iframeSrcPattern()).OnElements("iframe")
	p.AllowAttrs("width", "height").Matching(bluemonday.NumberOrPercent).OnElements("iframe")
	p.AllowAttrs("frameborder", "allowfullscreen", "allow", "scrolling").OnElements("iframe")
	return p
}

// SanitizeHTML 過濾頁面 HTML，並移除 src 不在白名單內（因而被拿掉 src）的 iframe。
func SanitizeHTML(s string) string {
	return dropIframesWithoutSrc(policy.Sanitize(s))
}

func dropIframesWithoutSrc(s string) string {
	z := html.NewTokenizer(strings.NewReader(s))
	var b strings.Builder
	skipping := false
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return b.String()
		}
		raw := string(z.Raw())
		switch tt {
		case html.StartTagToken, html.SelfClosingTagToken:
			name, hasAttr := z.TagName()
			if string(name) == "iframe" && !hasSrc(z, hasAttr) {
				skipping = tt == html.StartTagToken
				continue
			}
		case html.EndTagToken:
			if name, _ := z.TagName(); skipping && string(name) == "iframe" {
				skipping = false
				continue
			}
		}
		if !skipping {
			b.WriteString(raw)
		}
	}
}

func hasSrc(z *html.Tokenizer, more bool) bool {
	for more {
		var key, val []byte
		key, val, more = z.TagAttr()
		if string(key) == "src" && len(val) > 0 {
			return true
		}
	}
	return false
}

// PlainText 取出文字內容（合併空白）並截斷為 maxRunes 個字。
func PlainText(s string, maxRunes int) string {
	z := html.NewTokenizer(strings.NewReader(s))
	var parts []string
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.TextToken {
			parts = append(parts, strings.Fields(string(z.Text()))...) // Fields 也會切開 &nbsp;（U+00A0）
		}
	}
	text := strings.Join(parts, " ")
	if utf8.RuneCountInString(text) <= maxRunes {
		return text
	}
	return string([]rune(text)[:maxRunes])
}

func FirstImageSrc(s string) string {
	z := html.NewTokenizer(strings.NewReader(s))
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return ""
		}
		if tt != html.StartTagToken && tt != html.SelfClosingTagToken {
			continue
		}
		name, more := z.TagName()
		if string(name) != "img" {
			continue
		}
		for more {
			var key, val []byte
			key, val, more = z.TagAttr()
			if string(key) == "src" && len(val) > 0 {
				return string(val)
			}
		}
	}
}
