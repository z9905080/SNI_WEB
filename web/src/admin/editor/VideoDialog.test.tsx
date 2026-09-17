import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, expect, it, vi } from 'vitest'
import VideoDialog from './VideoDialog'

window.matchMedia ??= () =>
  ({ matches: false, media: '', addEventListener: () => {}, removeEventListener: () => {} }) as unknown as MediaQueryList
Element.prototype.hasPointerCapture ??= () => false
Element.prototype.scrollIntoView ??= () => {}

afterEach(() => vi.unstubAllGlobals())

it('允許的網址可以插入，並清空輸入框；不允許的網址顯示錯誤且無法送出', async () => {
  const user = userEvent.setup()
  const onInsert = vi.fn()
  render(<VideoDialog open onOpenChange={() => {}} onInsert={onInsert} />)

  const input = screen.getByLabelText('影片網址')
  await user.type(input, 'https://evil.example/x')
  expect(screen.getByRole('alert')).toHaveTextContent('無法嵌入這個網址')
  expect(screen.getByRole('button', { name: '插入影片' })).toBeDisabled()

  await user.clear(input)
  await user.type(input, 'https://www.youtube.com/watch?v=abc123')
  expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  const submit = screen.getByRole('button', { name: '插入影片' })
  expect(submit).toBeEnabled()
  await user.click(submit)
  expect(onInsert).toHaveBeenCalledWith('https://www.youtube.com/embed/abc123')
  expect(input).toHaveValue('')
})
