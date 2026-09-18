import { useLayoutEffect, useRef, type MouseEvent } from 'react'
import { useNavigate } from 'react-router'
import { internalPath } from './legacyLinks'
import { wrapTables } from './wrapTables'

// html 已由後端過濾，這裡原樣輸出
export default function LegacyContent({ html }: { html: string }) {
  const ref = useRef<HTMLDivElement>(null)
  const navigate = useNavigate()

  useLayoutEffect(() => {
    if (ref.current) wrapTables(ref.current)
  }, [html])

  const onClick = (e: MouseEvent<HTMLDivElement>) => {
    if (e.defaultPrevented || e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) return
    const a = (e.target as Element).closest('a')
    if (!a || (a.target && a.target !== '_self') || a.hasAttribute('download')) return
    const to = internalPath(a.getAttribute('href') ?? '', window.location)
    if (!to) return
    e.preventDefault()
    navigate(to)
  }

  return <div ref={ref} className="legacy-content" onClick={onClick} dangerouslySetInnerHTML={{ __html: html }} />
}
