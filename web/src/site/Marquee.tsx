import { useEffect, useState, useSyncExternalStore } from 'react'
import type { Marquee as MarqueeItem } from '@/shared/types'

const reducedQuery = '(prefers-reduced-motion: reduce)'

export function usePrefersReducedMotion() {
  return useSyncExternalStore(
    (onChange) => {
      const mq = window.matchMedia?.(reducedQuery)
      mq?.addEventListener('change', onChange)
      return () => mq?.removeEventListener('change', onChange)
    },
    () => window.matchMedia?.(reducedQuery).matches ?? false,
    () => false,
  )
}

export default function Marquee({ items }: { items: MarqueeItem[] }) {
  const reduced = usePrefersReducedMotion()
  if (items.length === 0) return null
  return (
    <div role="region" aria-label="最新公告" className="bg-brand-dark text-sm font-medium sm:text-base">
      {reduced ? <Rotator items={items} /> : <Scroller items={items} />}
    </div>
  )
}

function Scroller({ items }: { items: MarqueeItem[] }) {
  const row = (copy: string) =>
    items.map((m) => (
      <span key={`${copy}-${m.id}`} style={{ color: m.color }} className="px-10">
        {m.text}
      </span>
    ))
  // 內容重複兩份、動畫位移 -50%，形成無縫循環
  return (
    <div className="group overflow-hidden py-2">
      <div
        className="marquee-track flex w-max whitespace-nowrap group-hover:[animation-play-state:paused]"
        style={{ animationDuration: `${Math.max(20, items.length * 12)}s` }}
      >
        {row('a')}
        <span aria-hidden="true" className="flex">
          {row('b')}
        </span>
      </div>
    </div>
  )
}

function Rotator({ items }: { items: MarqueeItem[] }) {
  const [index, setIndex] = useState(0)
  useEffect(() => {
    if (items.length < 2) return
    const timer = setInterval(() => setIndex((i) => (i + 1) % items.length), 5000)
    return () => clearInterval(timer)
  }, [items.length])
  const m = items[index % items.length]
  return (
    <p aria-live="polite" className="px-4 py-2 text-center" style={{ color: m.color }}>
      {m.text}
    </p>
  )
}
