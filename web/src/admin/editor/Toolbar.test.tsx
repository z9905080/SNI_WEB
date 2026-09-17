import { EditorContent, useEditor } from '@tiptap/react'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, expect, it, vi } from 'vitest'
import { createExtensions } from './extensions'
import Toolbar from './Toolbar'

// jsdom 沒有 ProseMirror 需要的量測 API
document.createRange ??= () => {
  const range = new Range()
  range.getBoundingClientRect = () => ({ top: 0, left: 0, bottom: 0, right: 0, width: 0, height: 0 }) as DOMRect
  range.getClientRects = () => [] as unknown as DOMRectList
  return range
}
Element.prototype.getBoundingClientRect ??= () =>
  ({ top: 0, left: 0, bottom: 0, right: 0, width: 0, height: 0 }) as DOMRect
Element.prototype.scrollIntoView ??= () => {}
window.matchMedia ??= () =>
  ({ matches: false, media: '', addEventListener: () => {}, removeEventListener: () => {} }) as unknown as MediaQueryList
// Radix DropdownMenu 用 pointer capture 判斷觸發互動，jsdom 沒有實作
Element.prototype.hasPointerCapture ??= () => false
Element.prototype.setPointerCapture ??= () => {}
Element.prototype.releasePointerCapture ??= () => {}

afterEach(() => vi.unstubAllGlobals())

function Harness({ onInsertImage = () => {}, onInsertVideo = () => {} } = {}) {
  const editor = useEditor({ extensions: createExtensions(), content: '<p>hello</p>', immediatelyRender: true })
  if (!editor) return null
  return (
    <>
      <Toolbar editor={editor} onInsertImage={onInsertImage} onInsertVideo={onInsertVideo} />
      <EditorContent editor={editor} />
    </>
  )
}

it('工具列有可存取名稱的按鈕', async () => {
  render(<Harness />)
  expect(screen.getByRole('toolbar', { name: '編輯工具' })).toBeInTheDocument()
  for (const label of ['粗體', '斜體', '底線', '刪除線', '連結', '插入圖片', '插入影片', '表格', '復原', '重做']) {
    expect(screen.getByRole('button', { name: label })).toBeInTheDocument()
  }
})

// 這個測試要放在任何「先點擊過按鈕」的案例之前：Radix DropdownMenu 在 jsdom 下，
// 若同一個檔案先前的測試已經 userEvent.click 過任何按鈕，觸發選單的 pointerdown
// 判斷會失敗（選單開不起來），實際瀏覽器沒有這個問題，純粹是 jsdom + Radix 的測試環境限制。
it('表格選單可以插入 3x3 表格', async () => {
  const user = userEvent.setup()
  render(<Harness />)
  await user.click(screen.getByRole('button', { name: '表格' }))
  await screen.findByRole('menu')
  await user.click(screen.getByRole('menuitem', { name: '插入 3×3 表格' }))
  expect(document.querySelectorAll('table td')).toHaveLength(9)
})

it('點擊粗體會切換該按鈕的 active 狀態', async () => {
  const user = userEvent.setup()
  render(<Harness />)
  const bold = screen.getByRole('button', { name: '粗體' })
  expect(bold).toHaveAttribute('aria-pressed', 'false')
  await user.click(bold)
  expect(bold).toHaveAttribute('aria-pressed', 'true')
})

it('插入圖片／插入影片按鈕會呼叫對應的 callback', async () => {
  const user = userEvent.setup()
  const onInsertImage = vi.fn()
  const onInsertVideo = vi.fn()
  render(<Harness onInsertImage={onInsertImage} onInsertVideo={onInsertVideo} />)
  await user.click(screen.getByRole('button', { name: '插入圖片' }))
  await user.click(screen.getByRole('button', { name: '插入影片' }))
  expect(onInsertImage).toHaveBeenCalledTimes(1)
  expect(onInsertVideo).toHaveBeenCalledTimes(1)
})

