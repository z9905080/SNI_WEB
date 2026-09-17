import { act, render, screen } from '@testing-library/react'
import { afterEach, expect, it, vi } from 'vitest'
import Marquee from './Marquee'

const items = [
  { id: 1, text: '合掌感謝！', color: '#1EFF00' },
  { id: 2, text: '改版上線', color: '#FFFFFF' },
]

function mockReducedMotion(matches: boolean) {
  vi.stubGlobal('matchMedia', (query: string) => ({
    matches,
    media: query,
    addEventListener: () => {},
    removeEventListener: () => {},
  }))
}

afterEach(() => {
  vi.unstubAllGlobals()
  vi.useRealTimers()
})

it('捲動模式：每則使用自己的顏色', () => {
  mockReducedMotion(false)
  render(<Marquee items={items} />)
  const region = screen.getByRole('region', { name: '最新公告' })
  const first = region.querySelector('span')!
  expect(first).toHaveTextContent('合掌感謝！')
  expect(first).toHaveStyle({ color: '#1EFF00' })
  expect(region.querySelector('.marquee-track')).not.toBeNull()
})

it('減少動態時改為靜態輪替', () => {
  vi.useFakeTimers()
  mockReducedMotion(true)
  render(<Marquee items={items} />)
  expect(screen.getByText('合掌感謝！')).toBeInTheDocument()
  expect(document.querySelector('.marquee-track')).toBeNull()
  act(() => vi.advanceTimersByTime(5000))
  expect(screen.getByText('改版上線')).toBeInTheDocument()
  expect(screen.queryByText('合掌感謝！')).toBeNull()
})

it('沒有資料時不顯示', () => {
  mockReducedMotion(false)
  const { container } = render(<Marquee items={[]} />)
  expect(container).toBeEmptyDOMElement()
})
