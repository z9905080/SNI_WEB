import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, RouterProvider, useLocation } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { retryUnlessClientError } from '@/shared/api'
import AdminApp from './AdminApp'
import { safeRedirect } from './auth'

describe('safeRedirect', () => {
  it.each([
    [null, '/admin/groups'],
    ['/admin/pages/3?x=1', '/admin/pages/3?x=1'],
    ['/admin', '/admin/groups'],
    ['//evil.example/admin/', '/admin/groups'],
    ['https://evil.example/admin/x', '/admin/groups'],
    ['/page/3', '/admin/groups'],
  ])('%s → %s', (input, want) => {
    expect(safeRedirect(input)).toBe(want)
  })
})

function Where() {
  const l = useLocation()
  return <p data-testid="where">{l.pathname + l.search}</p>
}

function renderAdmin(path: string) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: retryUnlessClientError } } })
  const router = createMemoryRouter(
    [{ path: '/admin/*', element: <><AdminApp /><Where /></> }],
    { initialEntries: [path] },
  )
  render(
    <QueryClientProvider client={qc}>
      <RouterProvider router={router} />
    </QueryClientProvider>,
  )
}

const json = (status: number, body: unknown) => new Response(JSON.stringify(body), { status })

afterEach(() => vi.unstubAllGlobals())

describe('登入流程', () => {
  it('未登入時導向登入頁並帶 redirect，登入後回到原頁', async () => {
    let loggedIn = false
    vi.stubGlobal('fetch', vi.fn(async (url: string, init?: RequestInit) => {
      if (url.endsWith('/auth/me')) {
        return loggedIn
          ? json(200, { user: { id: 1, account: 'admin', name: '管理者' } })
          : json(401, { error: { code: 'unauthorized', message: '請先登入' } })
      }
      if (url.endsWith('/auth/login')) {
        const body = JSON.parse(String(init?.body))
        if (body.password !== 'password1') return json(401, { error: { code: 'invalid_credentials', message: '帳號或密碼錯誤' } })
        loggedIn = true
        return json(200, { user: { id: 1, account: 'admin', name: '管理者' } })
      }
      if (url.endsWith('/admin/settings')) {
        return json(200, { settings: { web_title: '生長之家', web_sub_title: '', facebook_url: '' } })
      }
      return json(404, { error: { code: 'not_found', message: '找不到' } })
    }))

    renderAdmin('/admin/settings')
    // 導向是 401 錯誤經由 queryCache 訂閱非同步觸發，findByTestId 只等元素存在（一開始就存在），
    // 不會等內容變成導向後的網址，這裡改用 waitFor 包住斷言本身，直到內容符合或逾時
    await waitFor(() =>
      expect(screen.getByTestId('where')).toHaveTextContent('/admin/login?redirect=%2Fadmin%2Fsettings'),
    )

    const user = userEvent.setup()
    await user.type(screen.getByLabelText('帳號'), 'admin')
    await user.type(screen.getByLabelText('密碼'), 'wrong')
    await user.click(screen.getByRole('button', { name: '登入' }))
    expect(await screen.findByRole('alert')).toHaveTextContent('帳號或密碼錯誤')

    await user.clear(screen.getByLabelText('密碼'))
    await user.type(screen.getByLabelText('密碼'), 'password1')
    await user.click(screen.getByRole('button', { name: '登入' }))
    // 手機與桌機版面各顯示一次使用者名稱（jsdom 不套用 CSS）
    expect(await screen.findAllByText('管理者')).not.toHaveLength(0)
    expect(screen.getByTestId('where')).toHaveTextContent('/admin/settings')
  })
})
