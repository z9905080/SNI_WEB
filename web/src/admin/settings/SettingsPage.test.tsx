import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, expect, it, vi } from 'vitest'
import { Toaster } from '@/components/ui/sonner'
import SettingsPage from './SettingsPage'

const json = (status: number, body: unknown) => new Response(JSON.stringify(body), { status })

afterEach(() => vi.unstubAllGlobals())

function setup(putResponse: Response) {
  const fetchMock = vi.fn(async (_url: string, init?: RequestInit) =>
    init?.method === 'PUT'
      ? putResponse
      : json(200, { settings: { web_title: '生長之家', web_sub_title: '一句簡單的感謝', facebook_url: '' } }),
  )
  vi.stubGlobal('fetch', fetchMock)
  render(
    <QueryClientProvider client={new QueryClient()}>
      <SettingsPage />
      <Toaster />
    </QueryClientProvider>,
  )
  return { fetchMock, user: userEvent.setup() }
}

it('載入並儲存網站設定', async () => {
  const saved = { web_title: '生長之家台灣', web_sub_title: '一句簡單的感謝', facebook_url: 'https://www.facebook.com/seichonoie.tw' }
  const { fetchMock, user } = setup(json(200, { settings: saved }))

  const title = await screen.findByLabelText('網站名稱')
  expect(title).toHaveValue('生長之家')
  await user.clear(title)
  await user.type(title, '生長之家台灣')
  await user.type(screen.getByLabelText('Facebook 粉絲專頁網址'), 'https://www.facebook.com/seichonoie.tw')
  await user.click(screen.getByRole('button', { name: '儲存設定' }))

  expect(await screen.findByText('已儲存設定')).toBeInTheDocument()
  expect(fetchMock).toHaveBeenLastCalledWith('/api/v1/admin/settings', expect.objectContaining({ method: 'PUT', body: JSON.stringify(saved) }))
})

it('顯示伺服器的驗證訊息', async () => {
  const { user } = setup(json(400, { error: { code: 'validation', message: 'Facebook 網址需為 http(s):// 開頭' } }))
  await user.click(await screen.findByRole('button', { name: '儲存設定' }))
  expect(await screen.findByText('Facebook 網址需為 http(s):// 開頭')).toBeInTheDocument()
})
