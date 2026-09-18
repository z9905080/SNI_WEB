import { describe, expect, it } from 'vitest'
import { readPublicConfig } from './config'

function doc(html: string) {
  return new DOMParser().parseFromString(`<html><head>${html}</head></html>`, 'text/html')
}

describe('readPublicConfig', () => {
  it('讀取注入的設定', () => {
    const d = doc('<script id="sni-config" type="application/json">{"gaMeasurementId":"G-1","counterScriptUrl":"https://c/x"}</script>')
    expect(readPublicConfig(d)).toEqual({ gaMeasurementId: 'G-1', counterScriptUrl: 'https://c/x' })
  })

  it('沒有注入（vite dev）或格式錯誤時回傳空設定', () => {
    const empty = { gaMeasurementId: '', counterScriptUrl: '' }
    expect(readPublicConfig(doc('<!--sni:head-->'))).toEqual(empty)
    expect(readPublicConfig(doc('<script id="sni-config" type="application/json">{oops</script>'))).toEqual(empty)
  })
})
