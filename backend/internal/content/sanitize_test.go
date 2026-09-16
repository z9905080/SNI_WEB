package content

import (
	"strings"
	"testing"
)

// 取自舊站實際內容的片段
const legacySample = `<p><a title="傳道協會FB" href="https://www.facebook.com/seichonoie.tw"><img src="http://www.seicho-no-ie.org.tw/php/picture/2019-12-21_18-41-17.jpg" /></a> &nbsp;</p>
<p><span style="font-size: 14pt; color: #e03e2d; font-family: arial, helvetica, sans-serif;">●官方網站重新改版</span></p>
<table style="border-collapse: collapse; width: 100%;" border="1"><tbody><tr style="height: 34px;"><td style="width: 100%; background-color: #18a085; text-align: center;"><strong>國際本部</strong></td></tr>
<tr><td><iframe src="https://players.brightcove.net/5049773559001/gjrCcgkq9_default/index.html?videoId=6001446559001" width="100%" height="225" frameborder="0" allowfullscreen="allowfullscreen"></iframe></td></tr>
<tr><td><iframe src="//www.youtube.com/embed/V4LhnGIPZt4" width="100%" height="225" allowfullscreen="allowfullscreen"></iframe></td></tr></tbody></table>`

func TestSanitizeKeepsLegacyFormatting(t *testing.T) {
	out := SanitizeHTML(legacySample)
	for _, want := range []string{
		`style="font-size: 14pt; color: #e03e2d; font-family: arial, helvetica, sans-serif;"`,
		`border="1"`,
		`style="border-collapse: collapse; width: 100%;"`,
		`<img src="http://www.seicho-no-ie.org.tw/php/picture/2019-12-21_18-41-17.jpg"`,
		`href="https://www.facebook.com/seichonoie.tw"`,
		`title="傳道協會FB"`,
		`src="https://players.brightcove.net/5049773559001/gjrCcgkq9_default/index.html?videoId=6001446559001"`,
		`src="//www.youtube.com/embed/V4LhnGIPZt4"`,
		`allowfullscreen="allowfullscreen"`,
		`width="100%"`,
		`<strong>國際本部</strong>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("輸出缺少 %s\n---\n%s", want, out)
		}
	}
	if strings.Contains(out, "nofollow") {
		t.Errorf("不應自動加入 rel=nofollow：%s", out)
	}
}

func TestSanitizeRemovesDangerousContent(t *testing.T) {
	in := `<p onclick="alert(1)" class="lead">hi</p><script>alert(2)</script>` +
		`<a href="javascript:alert(3)">x</a><img src="x.jpg" onerror="alert(4)">` +
		`<iframe src="https://evil.example.com/embed"></iframe>` +
		`<iframe src="https://www.youtube.com.evil.com/embed"></iframe>` +
		`<iframe src="javascript:alert(5)"></iframe>`
	out := SanitizeHTML(in)
	for _, bad := range []string{"onclick", "<script", "alert(2)", "javascript:", "onerror", "evil", "<iframe"} {
		if strings.Contains(out, bad) {
			t.Errorf("輸出不應包含 %q：%s", bad, out)
		}
	}
	if !strings.Contains(out, `<p class="lead">hi</p>`) {
		t.Errorf("應保留 class：%s", out)
	}
}

func TestSanitizeAllowsWhitelistedIframes(t *testing.T) {
	for _, src := range []string{
		"https://www.youtube.com/embed/a",
		"https://youtube.com/embed/a",
		"https://www.youtube-nocookie.com/embed/a",
		"https://drive.google.com/file/d/x/preview",
		"https://www.facebook.com/plugins/video.php?href=x",
	} {
		out := SanitizeHTML(`<iframe src="` + src + `"></iframe>`)
		if !strings.Contains(out, "<iframe") {
			t.Errorf("應允許 %s：%s", src, out)
		}
	}
}

func TestPlainText(t *testing.T) {
	in := "<p>第一段&nbsp;文字</p>\n<table><tr><td>表格</td></tr></table><p>  多   空白 </p>"
	if got := PlainText(in, 120); got != "第一段 文字 表格 多 空白" {
		t.Fatalf("got %q", got)
	}
	if got := PlainText("<p>一二三四五</p>", 3); got != "一二三" {
		t.Fatalf("截斷錯誤：%q", got)
	}
	if got := PlainText("", 10); got != "" {
		t.Fatalf("空字串：%q", got)
	}
}

func TestFirstImageSrc(t *testing.T) {
	in := `<p>x</p><img alt="a"><img src="/php/picture/a.jpg"><img src="/php/picture/b.jpg">`
	if got := FirstImageSrc(in); got != "/php/picture/a.jpg" {
		t.Fatalf("got %q", got)
	}
	if got := FirstImageSrc("<p>none</p>"); got != "" {
		t.Fatalf("got %q", got)
	}
}
