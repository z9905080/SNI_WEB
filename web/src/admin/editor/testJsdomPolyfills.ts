// jsdom 沒有實作 ProseMirror／Radix 在編輯器測試中會用到的幾個瀏覽器 API。
// 只給受影響的測試檔（Toolbar/VideoDialog/PageEditorPage）明確 import，不放進共用的 test/setup.ts。
import { afterEach } from 'vitest'

// jsdom 原生就有 document.createRange（`??=` 守不到它），但回傳的 Range 沒有
// getClientRects／getBoundingClientRect：ProseMirror 的 coordsAtPos／scrollToSelection
// 在每次選取改變時都會呼叫，缺了會直接丟出 TypeError。要包住「回傳值」而不是整個函式。
const nativeCreateRange = document.createRange.bind(document)
document.createRange = () => {
  const range = nativeCreateRange()
  range.getBoundingClientRect ??= () =>
    ({ x: 0, y: 0, top: 0, left: 0, right: 0, bottom: 0, width: 0, height: 0, toJSON() { return this } }) as DOMRect
  range.getClientRects ??= () => ({ item: () => null, length: 0, [Symbol.iterator]: Array.prototype[Symbol.iterator] }) as unknown as DOMRectList
  return range
}

Element.prototype.scrollIntoView ??= () => {}
window.matchMedia ??= (query: string) =>
  ({ matches: false, media: query, addEventListener: () => {}, removeEventListener: () => {} }) as unknown as MediaQueryList
// Radix DropdownMenu／Dialog 觸發判斷用的 pointer capture，jsdom 完全沒有實作
Element.prototype.hasPointerCapture ??= () => false
Element.prototype.setPointerCapture ??= () => {}
Element.prototype.releasePointerCapture ??= () => {}

// Radix `@radix-ui/react-dismissable-layer` 的 DismissableLayerContext 在這個 app 裡沒有 Provider，
// 是跨測試共用的模組級單例；理論上一旦某個 dropdown/dialog（例如 disableOutsidePointerEvents
// 的選單）在測試結束時沒有正常收尾，就可能把 document.body.style.pointerEvents 卡在 'none'，
// 讓後面的 userEvent 點擊全部失效。每個 it() 結束後強制清掉，當作保險。
afterEach(() => {
  document.body.style.pointerEvents = ''
})
