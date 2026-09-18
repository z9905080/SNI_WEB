import '@testing-library/jest-dom/vitest'
import { cleanup } from '@testing-library/react'
import { afterEach } from 'vitest'

// vitest 沒有預設全域 afterEach，RTL 的自動 cleanup 偵測不到，需手動註冊
// 否則同一個 test 檔案裡多個 it() 的 render() 會疊加，造成重複元素
afterEach(() => cleanup())

// jsdom 沒有的瀏覽器 API（sonner 讀取 prefers-color-scheme、Radix Popper 量測尺寸）
if (!window.matchMedia) {
  window.matchMedia = (query: string) =>
    ({
      matches: false,
      media: query,
      onchange: null,
      addEventListener: () => {},
      removeEventListener: () => {},
      addListener: () => {},
      removeListener: () => {},
      dispatchEvent: () => false,
    }) as MediaQueryList
}
window.ResizeObserver ??= class {
  observe() {}
  unobserve() {}
  disconnect() {}
}
Element.prototype.scrollIntoView ??= () => {}
