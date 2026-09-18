import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import { afterEach, expect, it, vi } from 'vitest'
import Banner from './Banner'

const items = [{ id: 1, image: '/img/1.jpg', url: '/page/3' }]

// jsdom 沒有 matchMedia／ResizeObserver／IntersectionObserver，embla 內部都會用到
vi.stubGlobal('matchMedia', (query: string) => ({
  matches: false,
  media: query,
  addEventListener: () => {},
  removeEventListener: () => {},
}))
class ObserverStub {
  observe() {}
  unobserve() {}
  disconnect() {}
}
vi.stubGlobal('ResizeObserver', ObserverStub)
vi.stubGlobal('IntersectionObserver', ObserverStub)

afterEach(() => vi.unstubAllGlobals())

it('輪播連結有可讀名稱', () => {
  render(
    <MemoryRouter>
      <Banner items={items} />
    </MemoryRouter>,
  )
  const link = screen.getByRole('link')
  expect(link).toHaveAccessibleName()
})
