import { beforeEach, expect, it, vi } from 'vitest'

beforeEach(() => {
  vi.resetModules()
  document.head.innerHTML = ''
  delete window.dataLayer
  delete window.gtag
})

it('沒有 ID 時不載入', async () => {
  const { initAnalytics, trackPageView } = await import('./analytics')
  initAnalytics('')
  trackPageView('/')
  expect(document.querySelector('script')).toBeNull()
  expect(window.dataLayer).toBeUndefined()
})

it('載入 gtag 一次並送出 page_view', async () => {
  const { initAnalytics, trackPageView } = await import('./analytics')
  initAnalytics('G-TEST')
  initAnalytics('G-TEST')
  const scripts = document.querySelectorAll('script')
  expect(scripts).toHaveLength(1)
  expect(scripts[0].src).toBe('https://www.googletagmanager.com/gtag/js?id=G-TEST')

  trackPageView('/page/3')
  const events = window.dataLayer!.map((args) => Array.from(args as ArrayLike<unknown>))
  expect(events).toContainEqual(['config', 'G-TEST', { send_page_view: false }])
  expect(events).toContainEqual(['event', 'page_view', expect.objectContaining({ page_path: '/page/3' })])
})
