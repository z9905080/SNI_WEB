import { render, screen } from '@testing-library/react'
import { expect, it } from 'vitest'
import Footer from './Footer'

it('有設定時顯示 Facebook 連結與計數器', () => {
  render(<Footer title="生長之家" facebookUrl="https://www.facebook.com/x" counterScriptUrl={'https://c.example/s.php?a=1&b="x"'} />)
  const fb = screen.getByRole('link', { name: /Facebook/ })
  expect(fb).toHaveAttribute('href', 'https://www.facebook.com/x')
  expect(fb).toHaveAttribute('target', '_blank')
  expect(fb).toHaveAttribute('rel', 'noopener noreferrer')
  const frame = screen.getByTitle('瀏覽次數')
  expect(frame.getAttribute('srcdoc')).toContain('src="https://c.example/s.php?a=1&amp;b=&quot;x&quot;"')
  expect(frame).toHaveAttribute('sandbox', 'allow-scripts')
})

it('沒有設定時不顯示', () => {
  render(<Footer title="生長之家" facebookUrl="" counterScriptUrl="" />)
  expect(screen.queryByRole('link')).toBeNull()
  expect(screen.queryByTitle('瀏覽次數')).toBeNull()
  expect(screen.getByText(/Seicho-No-Ie/)).toBeInTheDocument()
})
