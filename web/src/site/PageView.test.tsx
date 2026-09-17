import type { UseQueryResult } from '@tanstack/react-query'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router'
import { expect, it, vi } from 'vitest'
import { ApiError } from '@/shared/api'
import type { PageData } from '@/shared/types'
import PageView from './PageView'

vi.mock('./Banner', () => ({ default: () => <div data-testid="banner" /> }))

const q = (over: Partial<UseQueryResult<PageData>>) =>
  ({ isPending: false, isError: false, data: undefined, error: null, refetch: vi.fn(), ...over }) as UseQueryResult<PageData>

const renderView = (query: UseQueryResult<PageData>, props: { greeting?: string; groupName?: string } = { groupName: '主要行事活動' }) =>
  render(
    <MemoryRouter>
      <PageView query={query} {...props} />
    </MemoryRouter>,
  )

it('載入中顯示骨架', () => {
  renderView(q({ isPending: true }))
  expect(screen.getByLabelText('載入中')).toBeInTheDocument()
})

it('錯誤時顯示訊息並可重試', async () => {
  const query = q({ isError: true, error: new ApiError(0, 'network', '無法連線到伺服器') })
  renderView(query)
  expect(screen.getByRole('alert')).toHaveTextContent('無法連線到伺服器')
  await userEvent.click(screen.getByRole('button', { name: '重新載入' }))
  expect(query.refetch).toHaveBeenCalled()
})

it('內頁顯示所在頁籤、標題與內容', () => {
  renderView(q({ data: { page: { id: 3, group_id: 2, name: '練成會', html: '<p>日程</p>' }, carousels: [], marquees: [] } }))
  expect(screen.getByRole('heading', { name: '練成會' })).toBeInTheDocument()
  expect(screen.getByText('主要行事活動')).toBeInTheDocument()
  expect(screen.getByText('日程')).toBeInTheDocument()
})

it('首頁顯示問候語而不顯示頁面標題；沒有內容時顯示提示', () => {
  renderView(q({ data: { page: null, carousels: [], marquees: [] } }), { greeting: '一句簡單的感謝' })
  expect(screen.getByText('一句簡單的感謝')).toBeInTheDocument()
  expect(screen.queryByRole('heading')).toBeNull()
  expect(screen.getByText('這裡還沒有內容。')).toBeInTheDocument()
})
