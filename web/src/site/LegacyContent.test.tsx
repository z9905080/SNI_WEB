import { fireEvent, render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes, useLocation } from 'react-router'
import { expect, it } from 'vitest'
import LegacyContent from './LegacyContent'

function Where() {
  return <p data-testid="where">{useLocation().pathname}</p>
}

function setup(html: string) {
  render(
    <MemoryRouter initialEntries={['/page/1']}>
      <LegacyContent html={html} />
      <Routes>
        <Route path="*" element={<Where />} />
      </Routes>
    </MemoryRouter>,
  )
}

it('包裝表格', () => {
  setup('<table><tr><td>1</td></tr></table>')
  expect(document.querySelector('.legacy-content .table-scroll > table')).not.toBeNull()
})

it('站內舊連結改用前端路由', () => {
  setup(`<p><a href="http://${location.hostname}/#/3"><strong>練成會</strong></a></p>`)
  const notPrevented = fireEvent.click(screen.getByText('練成會'))
  expect(notPrevented).toBe(false)
  expect(screen.getByTestId('where')).toHaveTextContent('/page/3')
})

it.each([
  ['外部連結', '<a href="https://www.facebook.com/x">連結</a>', {}],
  ['開新分頁', '<a href="/#/3" target="_blank">連結</a>', {}],
  ['按住 Ctrl', '<a href="/#/3">連結</a>', { ctrlKey: true }],
  ['同頁錨點', '<a href="#top">連結</a>', {}],
])('%s維持瀏覽器預設行為', (_, html, init) => {
  setup(html)
  const notPrevented = fireEvent.click(screen.getByText('連結'), init)
  expect(notPrevented).toBe(true)
  expect(screen.getByTestId('where')).toHaveTextContent('/page/1')
})
