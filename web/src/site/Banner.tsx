import { useEffect, useState, type ReactNode } from 'react'
import { Link } from 'react-router'
import useEmblaCarousel from 'embla-carousel-react'
import Autoplay from 'embla-carousel-autoplay'
import type { Carousel } from '@/shared/types'
import { internalPath } from './legacyLinks'
import { usePrefersReducedMotion } from './Marquee'

export default function Banner({ items }: { items: Carousel[] }) {
  const reduced = usePrefersReducedMotion()
  const [viewportRef, embla] = useEmblaCarousel({ loop: true }, [
    Autoplay({ delay: 5000, active: !reduced, stopOnInteraction: false, stopOnMouseEnter: true }),
  ])
  const [selected, setSelected] = useState(0)

  useEffect(() => {
    if (!embla) return
    const onSelect = () => setSelected(embla.selectedScrollSnap())
    onSelect()
    embla.on('select', onSelect).on('reInit', onSelect)
    return () => {
      embla.off('select', onSelect).off('reInit', onSelect)
    }
  }, [embla])

  if (items.length === 0) return null
  return (
    <section aria-roledescription="carousel" aria-label="輪播圖" className="relative mx-auto max-w-6xl sm:px-4 sm:pt-6">
      <div ref={viewportRef} className="overflow-hidden sm:rounded-md">
        <div className="flex">
          {items.map((c, i) => (
            <div key={c.id} aria-roledescription="slide" aria-label={`第 ${i + 1} 張，共 ${items.length} 張`} className="min-w-0 flex-[0_0_100%]">
              <SlideLink url={c.url}>
                <img
                  src={c.image}
                  alt=""
                  loading={i === 0 ? 'eager' : 'lazy'}
                  className="aspect-[16/7] w-full bg-brand-soft object-cover"
                />
              </SlideLink>
            </div>
          ))}
        </div>
      </div>
      {items.length > 1 && (
        <>
          <ArrowButton label="上一張" side="left" onClick={() => embla?.scrollPrev()} />
          <ArrowButton label="下一張" side="right" onClick={() => embla?.scrollNext()} />
          <div className="absolute inset-x-0 bottom-3 flex justify-center gap-2">
            {items.map((c, i) => (
              <button
                key={c.id}
                type="button"
                aria-label={`第 ${i + 1} 張`}
                aria-current={i === selected}
                onClick={() => embla?.scrollTo(i)}
                className={`h-2.5 rounded-full transition-all ${i === selected ? 'w-6 bg-white' : 'w-2.5 bg-white/60'}`}
              />
            ))}
          </div>
        </>
      )}
    </section>
  )
}

function ArrowButton({ label, side, onClick }: { label: string; side: 'left' | 'right'; onClick: () => void }) {
  return (
    <button
      type="button"
      aria-label={label}
      onClick={onClick}
      className={`absolute top-1/2 hidden h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-black/30 text-2xl text-white hover:bg-black/50 sm:flex ${side === 'left' ? 'left-6' : 'right-6'}`}
    >
      {side === 'left' ? '‹' : '›'}
    </button>
  )
}

function SlideLink({ url, children }: { url: string; children: ReactNode }) {
  if (!url) return children
  const to = internalPath(url, window.location)
  if (to) return <Link to={to}>{children}</Link>
  return (
    <a href={url} target="_blank" rel="noopener noreferrer">
      {children}
    </a>
  )
}
