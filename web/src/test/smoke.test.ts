import { readFileSync } from 'node:fs'
import { expect, test } from 'vitest'

test('index.html 保留 meta 佔位符且沒有自帶 title', () => {
  // 注意：不可寫成 `new URL('../../index.html', import.meta.url)` 字面量寫法——
  // Vite/Vitest 會將這個寫死的 pattern 當成靜態資源 URL 特例處理，
  // 在測試環境下會解析成 http://localhost:3000/index.html 而不是實際檔案路徑。
  // 先把 import.meta.url 存到變數，繞過這個靜態分析。
  const testFileUrl = import.meta.url
  const html = readFileSync(new URL('../../index.html', testFileUrl), 'utf8')
  expect(html).toContain('<!--sni:head-->')
  expect(html).not.toMatch(/<title>/i)
})
