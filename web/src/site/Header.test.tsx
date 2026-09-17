import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes, useLocation } from 'react-router'
import { expect, it } from 'vitest'
import Header from './Header'
import type { Group } from '@/shared/types'

const menu: Group[] = [
  { id: 1, name: '首頁', pages: [{ id: 1, name: '最新消息' }] },
  { id: 2, name: '主要行事活動', pages: [{ id: 3, name: '練成會' }, { id: 4, name: '誌友會' }] },
]

function Where() {
  return <p data-testid="where">{useLocation().pathname}</p>
}

function setup(path = '/') {
  const user = userEvent.setup()
  render(
    <MemoryRouter initialEntries={[path]}>
      <Header title="生長之家" subtitle="一句簡單的感謝" menu={menu} />
      <Routes>
        <Route path="*" element={<Where />} />
      </Routes>
    </MemoryRouter>,
  )
  return user
}

it('滑過群組顯示頁面，點擊後導覽並關閉', async () => {
  const user = setup()
  const nav = screen.getByRole('navigation', { name: '主選單' })
  const trigger = within(nav).getByRole('button', { name: '主要行事活動' })
  expect(trigger).toHaveAttribute('aria-expanded', 'false')

  await user.hover(trigger)
  expect(trigger).toHaveAttribute('aria-expanded', 'true')
  await user.click(within(nav).getByRole('link', { name: '誌友會' }))
  expect(screen.getByTestId('where')).toHaveTextContent('/page/4')
  expect(within(nav).queryByRole('link', { name: '誌友會' })).toBeNull()
})

it('標示目前頁面所在的頁籤', () => {
  setup('/page/4')
  const nav = screen.getByRole('navigation', { name: '主選單' })
  expect(within(nav).getByRole('button', { name: '主要行事活動' })).toHaveAttribute('data-active', 'true')
  expect(within(nav).getByRole('button', { name: '首頁' })).not.toHaveAttribute('data-active')
  expect(screen.getByText('一句簡單的感謝')).toBeInTheDocument()
})

it('鍵盤：Enter 開啟、方向鍵進入選單、Escape 關閉並回到按鈕', async () => {
  const user = setup()
  const nav = screen.getByRole('navigation', { name: '主選單' })
  const trigger = within(nav).getByRole('button', { name: '主要行事活動' })
  trigger.focus()
  await user.keyboard('{ArrowDown}')
  await expect.poll(() => document.activeElement?.textContent).toBe('練成會')
  await user.keyboard('{Escape}')
  expect(trigger).toHaveAttribute('aria-expanded', 'false')
  expect(trigger).toHaveFocus()

  await user.keyboard('{Enter}')
  expect(trigger).toHaveAttribute('aria-expanded', 'true')
})

it('手機抽屜：手風琴展開群組並導覽', async () => {
  const user = setup()
  await user.click(screen.getByRole('button', { name: '開啟選單' }))
  const drawer = screen.getByRole('dialog', { name: '選單' })
  expect(within(drawer).queryByRole('link', { name: '練成會' })).toBeNull()
  await user.click(within(drawer).getByRole('button', { name: '主要行事活動' }))
  await user.click(within(drawer).getByRole('link', { name: '練成會' }))
  expect(screen.getByTestId('where')).toHaveTextContent('/page/3')
  expect(screen.queryByRole('dialog')).toBeNull()
})
