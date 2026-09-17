import '@testing-library/jest-dom/vitest'
import { cleanup } from '@testing-library/react'
import { afterEach } from 'vitest'

// vitest 沒有預設全域 afterEach，RTL 的自動 cleanup 偵測不到，需手動註冊
// 否則同一個 test 檔案裡多個 it() 的 render() 會疊加，造成重複元素
afterEach(() => cleanup())
