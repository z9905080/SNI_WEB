import { describe, expect, it } from 'vitest'
import { removeStyleProps, roundTrip } from './extensions'
import { detectLoss } from './lossDetection'

// 取自舊站實際內容的寫法
const legacy = `<p><a title="傳道協會FB" href="https://www.facebook.com/seichonoie.tw"><img src="/php/picture/2019-12-21_18-41-17.jpg"></a>&nbsp;</p>
<p style="text-align: left;"><span style="font-size: 14pt; color: #e03e2d; font-family: arial, helvetica, sans-serif;">●官方網站重新改版</span></p>
<p><span style="background-color: #fbeeb8;">每週日上午舉行。</span></p>
<h2 style="text-align: center;">標題</h2>
<table style="border-collapse: collapse; width: 100%;" border="1"><tbody>
<tr style="height: 34px;"><td style="width: 100%; background-color: #18a085; text-align: center;"><strong><span style="color: #ffffff;">國際本部發行的影片</span></strong></td></tr>
<tr><td><iframe src="//www.youtube.com/embed/V4LhnGIPZt4" width="100%" height="225" allowfullscreen="allowfullscreen"></iframe></td></tr>
</tbody></table>
<ul><li>項目一</li><li><u>項目二</u></li></ul>`

describe('roundTrip', () => {
  it('舊站常見格式不會遺失', () => {
    expect(detectLoss(legacy, roundTrip(legacy))).toEqual([])
  })

  it('不在白名單的 iframe 會被移除', () => {
    const out = roundTrip('<iframe src="https://evil.example/x"></iframe><p>後面</p>')
    expect(out).not.toContain('iframe')
    expect(detectLoss('<iframe src="https://evil.example/x"></iframe>', out)).toContain('<iframe>')
  })

  it('無法保留的結構會被偵測到', () => {
    const html = '<div class="box"><p>框內</p></div><p><font color="red">舊字型</font></p>'
    const lost = detectLoss(html, roundTrip(html))
    expect(lost).toEqual(expect.arrayContaining(['<div>', '<font>']))
  })

  it('連結不會被加上 target 與 rel', () => {
    expect(roundTrip('<p><a href="/page/3">練成會</a></p>')).toBe('<p><a href="/page/3">練成會</a></p>')
  })
})

describe('removeStyleProps', () => {
  it.each([
    [null, ['color'], null],
    ['color: red; font-weight: bold;', ['color'], 'font-weight: bold;'],
    ['text-align: left', ['text-align'], null],
    ['COLOR:red;;  width:1px', ['color'], 'width:1px;'],
    ['background: url(a;b)', [], 'background: url(a;b);'],
  ])('%s', (style, props, want) => {
    expect(removeStyleProps(style, props)).toBe(want)
  })
})
