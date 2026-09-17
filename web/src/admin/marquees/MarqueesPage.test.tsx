import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { fireEvent, render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router'
import { afterEach, expect, it, vi } from 'vitest'
import { Toaster } from '@/components/ui/sonner'
import MarqueesPage from './MarqueesPage'

const json = (status: number, body: unknown) => new Response(JSON.stringify(body), { status })

afterEach(() => vi.unstubAllGlobals())

it('新增跑馬燈時即時預覽文字與顏色，儲存後顯示結果', async () => {
  let marquees = [{ id: 1, text: '合掌感謝！', color: '#1EFF00' }]
  const fetchMock = vi.fn(async (_url: string, init?: RequestInit) => {
    if (init?.method === 'POST') {
      const m = { id: 2, ...JSON.parse(String(init.body)) }
      marquees = [...marquees, m]
      return json(201, { marquee: m })
    }
    return json(200, { marquees })
  })
  vi.stubGlobal('fetch', fetchMock)

  render(
    <QueryClientProvider client={new QueryClient()}>
      <MemoryRouter>
        <MarqueesPage />
      </MemoryRouter>
      <Toaster />
    </QueryClientProvider>,
  )
  const user = userEvent.setup()
  expect(await screen.findByText('合掌感謝！')).toBeInTheDocument()

  await user.click(screen.getByRole('button', { name: '新增跑馬燈' }))
  const dialog = screen.getByRole('dialog')
  await user.type(within(dialog).getByLabelText('文字'), '歡迎參加練成會')
  fireEvent.input(within(dialog).getByLabelText('顏色'), { target: { value: '#ffcc00' } })
  const preview = within(dialog).getByTestId('marquee-preview')
  expect(preview).toHaveTextContent('歡迎參加練成會')
  expect(preview).toHaveStyle({ color: '#ffcc00' })

  await user.click(within(dialog).getByRole('button', { name: '儲存' }))
  expect(await screen.findByText('已儲存跑馬燈')).toBeInTheDocument()
  expect(fetchMock).toHaveBeenCalledWith(
    '/api/v1/admin/marquees',
    expect.objectContaining({ method: 'POST', body: JSON.stringify({ text: '歡迎參加練成會', color: '#ffcc00' }) }),
  )
  expect(await screen.findByText('歡迎參加練成會')).toBeInTheDocument()
})
