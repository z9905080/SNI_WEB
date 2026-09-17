import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, RouterProvider } from 'react-router'
import { afterEach, beforeEach, expect, it, vi } from 'vitest'
import { Toaster } from '@/components/ui/sonner'
import PageEditorPage from './PageEditorPage'

// jsdom 沒有 ProseMirror／Radix Popper 需要的量測 API
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
Element.prototype.hasPointerCapture ??= () => false
Element.prototype.setPointerCapture ??= () => {}
Element.prototype.releasePointerCapture ??= () => {}

const json = (status: number, body: unknown) => new Response(body === null ? null : JSON.stringify(body), { status })

const groups = [
  { id: 1, name: '首頁', pages: [] },
  { id: 2, name: '主要行事活動', pages: [] },
]
// id 5：Tiptap 無法完整保留（div 結構），載入時應以原始碼模式開啟
const lossyPage = { id: 5, group_id: 1, name: '練成會', html: '<div class="box"><p>x</p></div>' }
// id 6：可以完整往返，載入時是視覺模式
const cleanPage = { id: 6, group_id: 1, name: '誌友會', html: '<p>hello</p>' }

let fetchMock: ReturnType<typeof vi.fn>

beforeEach(() => {
  fetchMock = vi.fn(async (url: string, init?: RequestInit) => {
    const method = init?.method ?? 'GET'
    if (url === '/api/v1/admin/groups' && method === 'GET') return json(200, { groups })
    if (url === '/api/v1/admin/pages/5' && method === 'GET') return json(200, { page: lossyPage })
    if (url === '/api/v1/admin/pages/6' && method === 'GET') return json(200, { page: cleanPage })
    if (url === '/api/v1/admin/pages' && method === 'POST') {
      const body = JSON.parse(String(init?.body))
      return json(201, { page: { id: 9, ...body } })
    }
    if (url === '/api/v1/admin/pages/6' && method === 'PATCH') {
      const body = JSON.parse(String(init?.body))
      return json(200, { page: { ...cleanPage, ...body } })
    }
    return json(404, { error: { code: 'not_found', message: '找不到' } })
  })
  vi.stubGlobal('fetch', fetchMock)
})

afterEach(() => vi.unstubAllGlobals())

function renderPage(initialPath: string) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  const router = createMemoryRouter(
    [
      { path: '/admin/pages/new', element: <PageEditorPage /> },
      { path: '/admin/pages/:id', element: <PageEditorPage /> },
      { path: '/admin/groups', element: <p>頁籤與頁面</p> },
    ],
    { initialEntries: [initialPath] },
  )
  render(
    <QueryClientProvider client={qc}>
      <RouterProvider router={router} />
      <Toaster />
    </QueryClientProvider>,
  )
  return { user: userEvent.setup(), router }
}

it('新增頁面：預設群組來自網址參數，儲存成功後導向新網址且不會被攔截', async () => {
  const { user, router } = renderPage('/admin/pages/new?group=2')
  const nameInput = await screen.findByLabelText('頁面名稱')
  expect(screen.getByLabelText('所屬頁籤')).toHaveValue('2')

  await user.type(nameInput, '關於我們')
  await user.click(screen.getByRole('button', { name: '儲存' }))

  expect(await screen.findByText('已儲存')).toBeInTheDocument()
  expect(fetchMock).toHaveBeenCalledWith(
    '/api/v1/admin/pages',
    expect.objectContaining({ method: 'POST', body: JSON.stringify({ name: '關於我們', group_id: 2, html: '' }) }),
  )
  expect(router.state.location.pathname).toBe('/admin/pages/9')
  // 儲存成功後 dirty 已重置，不應該出現未儲存變更的攔截對話框
  expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument()
})

it('載入時偵測到遺失格式，預設進入原始碼模式並顯示警告', async () => {
  renderPage('/admin/pages/5')
  const status = await screen.findByRole('status')
  expect(status).toHaveTextContent('HTML 原始碼模式')
  expect(status).toHaveTextContent('<div>')
  expect(screen.getByRole('button', { name: '視覺編輯' })).toHaveAttribute('aria-pressed', 'false')
  expect(screen.getByRole('button', { name: 'HTML 原始碼' })).toHaveAttribute('aria-pressed', 'true')
})

it('從原始碼模式切換到視覺編輯時，格式遺失需要先確認', async () => {
  const { user } = renderPage('/admin/pages/5')
  await screen.findByRole('status')

  await user.click(screen.getByRole('button', { name: '視覺編輯' }))
  const dialog = await screen.findByRole('alertdialog')
  expect(within(dialog).getByText('切換到視覺編輯？')).toBeInTheDocument()

  await user.click(within(dialog).getByRole('button', { name: '取消' }))
  expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'HTML 原始碼' })).toHaveAttribute('aria-pressed', 'true')

  await user.click(screen.getByRole('button', { name: '視覺編輯' }))
  await user.click(within(await screen.findByRole('alertdialog')).getByRole('button', { name: '切換到視覺編輯' }))
  expect(screen.getByRole('button', { name: '視覺編輯' })).toHaveAttribute('aria-pressed', 'true')
})

it('修改後嘗試離開頁面會被攔截，確認捨棄後才真的離開', async () => {
  const { user, router } = renderPage('/admin/pages/6')
  const nameInput = await screen.findByLabelText('頁面名稱')
  await user.type(nameInput, '！')

  await user.click(screen.getByRole('link', { name: '返回列表' }))
  const dialog = await screen.findByRole('alertdialog')
  expect(within(dialog).getByText('還有變更沒有儲存')).toBeInTheDocument()
  expect(router.state.location.pathname).toBe('/admin/pages/6')

  await user.click(within(dialog).getByRole('button', { name: '取消' }))
  expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument()
  expect(router.state.location.pathname).toBe('/admin/pages/6')

  await user.click(screen.getByRole('link', { name: '返回列表' }))
  await user.click(within(await screen.findByRole('alertdialog')).getByRole('button', { name: '捨棄變更並離開' }))
  expect(router.state.location.pathname).toBe('/admin/groups')
})
