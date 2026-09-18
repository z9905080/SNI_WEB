import { EditorContent, useEditor } from '@tiptap/react'
import { render, screen } from '@testing-library/react'
import userEvent, { PointerEventsCheckLevel } from '@testing-library/user-event'
import { expect, it } from 'vitest'
import { createExtensions } from './extensions'
import './testJsdomPolyfills'
import Toolbar from './Toolbar'

// 這個互動（開啟表格下拉選單）獨立成自己的檔案：Radix DropdownMenu 的 Portal 掛載和
// Tiptap Editor 的建立/銷毀時序，在同一個 jsdom 環境裡跑過第二次以上的 Editor 實例後，
// 偶爾會卡住不觸發（在 Toolbar.test.tsx 的其他測試之後才跑這個案例時重現過，即使拉高
// findByRole 逾時到 8 秒仍等不到）；獨立成單一檔案，讓它永遠是這個 jsdom 環境裡第一個
// （也是唯一一個）建立的 Tiptap Editor + DropdownMenu 組合，可以穩定重現「總是成功」的
// 那個情境，而不必依賴同一檔案裡測試的執行順序。
function Harness() {
  const editor = useEditor({ extensions: createExtensions(), content: '<p>hello</p>', immediatelyRender: true })
  if (!editor) return null
  return (
    <>
      <Toolbar editor={editor} onInsertImage={() => {}} onInsertVideo={() => {}} />
      <EditorContent editor={editor} />
    </>
  )
}

it('表格選單可以插入 3x3 表格', async () => {
  const user = userEvent.setup({ pointerEventsCheck: PointerEventsCheckLevel.Never })
  render(<Harness />)
  await user.click(screen.getByRole('button', { name: '表格' }))
  await screen.findByRole('menu')
  await user.click(screen.getByRole('menuitem', { name: '插入 3×3 表格' }))
  expect(document.querySelectorAll('table td')).toHaveLength(9)
})
