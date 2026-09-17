import { describe, expect, it } from 'vitest'
import { detectLoss } from './lossDetection'

describe('detectLoss', () => {
  it('等價的 HTML 不算遺失', () => {
    expect(detectLoss('<p><b>粗</b>&nbsp;<i>斜</i></p>', '<p><strong>粗</strong> <em>斜</em></p>')).toEqual([])
    expect(
      detectLoss(
        '<p style="color: #e03e2d; font-size: 14pt">x</p>',
        '<p style="font-size:14pt;color:rgb(224, 62, 45)">x</p>',
      ),
    ).toEqual([])
    expect(detectLoss('<table><tr><td>1</td></tr></table>', '<table><colgroup><col></colgroup><tbody><tr><td><p>1</p></td></tr></tbody></table>')).toEqual([])
    expect(detectLoss('<p class="a b">x</p>', '<p class="b a">x</p>')).toEqual([])
  })

  it('找出遺失的標籤、屬性、樣式與文字', () => {
    expect(detectLoss('<div class="box"><p>x</p></div>', '<p>x</p>')).toEqual(['<div>', 'class="box"'])
    expect(detectLoss('<table border="1"><tr><td>x</td></tr></table>', '<table><tr><td>x</td></tr></table>')).toEqual(['border="1"'])
    expect(detectLoss('<p style="font-weight: bold">x</p>', '<p>x</p>')).toEqual(['font-weight: bold'])
    expect(detectLoss('<p>一段文字</p>', '<p>一段</p>')).toEqual(['文字內容'])
    expect(detectLoss('<p><span>a</span><span>b</span></p>', '<p><span>ab</span></p>')).toEqual(['<span>'])
  })
})
