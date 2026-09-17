import { describe, expect, it } from 'vitest'
import { internalPath, legacyHashTarget } from './legacyLinks'

describe('legacyHashTarget', () => {
  it.each([
    ['#/5', '/page/5'],
    ['#/75/', '/page/75'],
    ['#/', '/'],
    ['#', null],
    ['', null],
    ['#/abc', null],
    ['#top', null],
  ])('%s → %s', (hash, want) => {
    expect(legacyHashTarget(hash)).toBe(want)
  })
})

describe('internalPath', () => {
  const current = { href: 'https://www.seicho-no-ie.org.tw/page/3', hostname: 'www.seicho-no-ie.org.tw' }

  it.each([
    ['舊站絕對網址（http）', 'http://www.seicho-no-ie.org.tw/#/5', '/page/5'],
    ['舊站首頁', 'http://www.seicho-no-ie.org.tw/#/', '/'],
    ['相對舊 hash', '/#/20', '/page/20'],
    ['新網址', '/page/7', '/page/7'],
    ['新網址帶查詢字串', 'https://www.seicho-no-ie.org.tw/page/7?x=1', '/page/7'],
    ['站內首頁', '/', '/'],
    ['同頁錨點', '#section', null],
    ['其他站內路徑', '/php/picture/a.jpg', null],
    ['後台', '/admin/groups', null],
    ['外部網站', 'https://www.facebook.com/seichonoie.tw', null],
    ['其他子網域', 'http://seicho-no-ie.org.tw/#/5', null],
    ['mailto', 'mailto:a@b.c', null],
    ['javascript', 'javascript:alert(1)', null],
    ['空字串', '', null],
  ])('%s', (_, href, want) => {
    expect(internalPath(href, current)).toBe(want)
  })
})
