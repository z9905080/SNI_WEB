import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, RouterProvider } from 'react-router'
import { afterEach, beforeEach, expect, it, vi } from 'vitest'
import { Toaster } from '@/components/ui/sonner'
import GroupsPage from './GroupsPage'

const json = (status: number, body: unknown) => new Response(body === null ? null : JSON.stringify(body), { status })

let groups: { id: number; name: string; pages: { id: number; name: string }[] }[]
let fetchMock: ReturnType<typeof vi.fn>

beforeEach(() => {
  groups = [
    { id: 1, name: '首頁', pages: [{ id: 1, name: '最新消息' }] },
    { id: 2, name: '主要行事活動', pages: [{ id: 3, name: '練成會' }, { id: 4, name: '誌友會' }] },
  ]
  fetchMock = vi.fn(async (url: string, init?: RequestInit) => {
    const method = init?.method ?? 'GET'
    if (url === '/api/v1/admin/groups' && method === 'GET') return json(200, { groups })
    if (url === '/api/v1/admin/groups' && method === 'POST') {
      const g = { id: 9, name: JSON.parse(String(init?.body)).name, pages: [] }
      groups = [...groups, g]
      return json(201, { group: g })
    }
    if (url === '/api/v1/admin/groups/2' && method === 'DELETE') {
      return json(409, { error: { code: 'group_not_empty', message: '此頁籤仍有頁面，請先移除或刪除頁面' } })
    }
    return json(404, { error: { code: 'not_found', message: '找不到' } })
  })
  vi.stubGlobal('fetch', fetchMock)
})

afterEach(() => vi.unstubAllGlobals())

function renderPage() {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  const router = createMemoryRouter([{ path: '/admin/groups', element: <GroupsPage /> }], {
    initialEntries: ['/admin/groups'],
  })
  render(
    <QueryClientProvider client={qc}>
      <RouterProvider router={router} />
      <Toaster />
    </QueryClientProvider>,
  )
  return userEvent.setup()
}

it('列出頁籤，選擇後顯示該頁籤的頁面', async () => {
  const user = renderPage()
  const groupList = await screen.findByRole('list', { name: '頁籤' })
  expect(within(groupList).getAllByRole('listitem')).toHaveLength(2)
  expect(screen.getByRole('link', { name: '最新消息' })).toHaveAttribute('href', '/admin/pages/1')

  await user.click(within(groupList).getByRole('button', { name: /^主要行事活動/ }))
  const pageList = screen.getByRole('list', { name: '頁面' })
  expect(within(pageList).getByRole('link', { name: '練成會' })).toHaveAttribute('href', '/admin/pages/3')
  expect(within(pageList).getByRole('link', { name: '誌友會' })).toHaveAttribute('href', '/admin/pages/4')
  expect(screen.getByRole('link', { name: '新增頁面' })).toHaveAttribute('href', '/admin/pages/new?group=2')
})

it('新增頁籤後自動選取', async () => {
  const user = renderPage()
  await user.click(await screen.findByRole('button', { name: '新增頁籤' }))
  await user.type(screen.getByLabelText('名稱'), '關於我們')
  await user.click(screen.getByRole('button', { name: '儲存' }))
  expect(await screen.findByText('已新增頁籤')).toBeInTheDocument()
  expect(await screen.findByText('此頁籤還沒有頁面。')).toBeInTheDocument()
  expect(fetchMock).toHaveBeenCalledWith('/api/v1/admin/groups', expect.objectContaining({ method: 'POST', body: '{"name":"關於我們"}' }))
})

it('刪除失敗時顯示 API 訊息；首頁頁籤不可刪除', async () => {
  const user = renderPage()
  expect(await screen.findByRole('button', { name: '刪除 首頁' })).toBeDisabled()
  await user.click(screen.getByRole('button', { name: '刪除 主要行事活動' }))
  await user.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: '刪除' }))
  expect(await screen.findByText('此頁籤仍有頁面，請先移除或刪除頁面')).toBeInTheDocument()
})
