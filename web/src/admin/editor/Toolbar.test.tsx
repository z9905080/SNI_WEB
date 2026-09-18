import { EditorContent, useEditor } from '@tiptap/react'
import { render, screen } from '@testing-library/react'
import userEvent, { PointerEventsCheckLevel } from '@testing-library/user-event'
import { expect, it, vi } from 'vitest'
import { createExtensions } from './extensions'
import './testJsdomPolyfills'
import Toolbar from './Toolbar'

const setupUser = () => userEvent.setup({ pointerEventsCheck: PointerEventsCheckLevel.Never })

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

// 「表格」下拉選單的開關互動獨立放在 ToolbarTableMenu.test.tsx，理由見該檔案開頭的註解。
it('點擊粗體會切換該按鈕的 active 狀態', async () => {
  const user = setupUser()
  render(<Harness />)
  const bold = screen.getByRole('button', { name: '粗體' })
  expect(bold).toHaveAttribute('aria-pressed', 'false')
  await user.click(bold)
  expect(bold).toHaveAttribute('aria-pressed', 'true')
})

it('插入圖片／插入影片按鈕會呼叫對應的 callback', async () => {
  const user = setupUser()
  const onInsertImage = vi.fn()
  const onInsertVideo = vi.fn()
  render(<Harness onInsertImage={onInsertImage} onInsertVideo={onInsertVideo} />)
  await user.click(screen.getByRole('button', { name: '插入圖片' }))
  await user.click(screen.getByRole('button', { name: '插入影片' }))
  expect(onInsertImage).toHaveBeenCalledTimes(1)
  expect(onInsertVideo).toHaveBeenCalledTimes(1)
})
