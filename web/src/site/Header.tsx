import { type ReactNode, useEffect, useLayoutEffect, useRef, useState } from 'react'
import { Link, useLocation } from 'react-router'
import type { Group } from '@/shared/types'

type Props = { title: string; subtitle: string; menu: Group[] }

const ITEM_CLASS = 'whitespace-nowrap border-b-2 border-transparent px-3 py-2 font-medium'
const LINK_CLASS =
  'block whitespace-nowrap px-4 py-2 text-ink hover:bg-brand-soft hover:text-brand focus-visible:bg-brand-soft focus-visible:outline-none'
const GAP = 4 // 對應 ul 的 gap-1

function MenuPanel({
  id,
  onClose,
  className = '',
  children,
}: {
  id: string
  onClose: () => void
  className?: string
  children: ReactNode
}) {
  return (
    <ul
      id={id}
      className={`absolute right-0 top-full min-w-52 rounded-md border border-slate-200 bg-white py-2 shadow-lg ${className}`}
      onKeyDown={(e) => {
        const links = [...e.currentTarget.querySelectorAll<HTMLElement>('a')]
        const i = links.indexOf(document.activeElement as HTMLElement)
        if (e.key === 'Escape') {
          onClose()
          ;(e.currentTarget.previousElementSibling as HTMLElement | null)?.focus()
        } else if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
          e.preventDefault()
          const next = e.key === 'ArrowDown' ? i + 1 : i - 1
          links[(next + links.length) % links.length]?.focus()
        }
      }}
    >
      {children}
    </ul>
  )
}

export default function Header({ title, subtitle, menu }: Props) {
  const [open, setOpen] = useState<number | 'more' | null>(null)
  const [drawer, setDrawer] = useState(false)
  const [compact, setCompact] = useState(false)
  const navRef = useRef<HTMLElement>(null)
  const barRef = useRef<HTMLDivElement>(null)
  const rulerRef = useRef<HTMLUListElement>(null)
  const [shown, setShown] = useState(menu.length)
  const { pathname } = useLocation()
  const activeGroupId = menu.find((g) => g.pages.some((p) => pathname === `/page/${p.id}`))?.id

  // 換頁時關閉所有選單
  const [lastPath, setLastPath] = useState(pathname)
  if (lastPath !== pathname) {
    setLastPath(pathname)
    setOpen(null)
    setDrawer(false)
  }

  // 依可用寬度決定橫列放得下幾個群組，其餘收進「更多」。
  // 寬度量不到時（jsdom、字型還沒載好）一律全部顯示。
  useLayoutEffect(() => {
    const bar = barRef.current
    const ruler = rulerRef.current
    if (!bar || !ruler) return
    const count = () => {
      const items = [...ruler.querySelectorAll<HTMLElement>('[data-navitem]')]
      const widths = items.map((el) => el.offsetWidth + GAP)
      const avail = bar.clientWidth
      const more = widths.pop() ?? 0
      const total = widths.reduce((a, b) => a + b, 0)
      if (!avail || !total || total <= avail) {
        setShown(menu.length)
        return
      }
      let used = more
      let n = 0
      for (const w of widths) {
        if (used + w > avail) break
        used += w
        n++
      }
      setShown(n)
    }
    count()
    const ro = new ResizeObserver(count)
    ro.observe(bar)
    document.fonts?.ready.then(count).catch(() => {})
    return () => ro.disconnect()
  }, [menu])

  useEffect(() => {
    const onScroll = () => setCompact(window.scrollY > 40)
    onScroll()
    window.addEventListener('scroll', onScroll, { passive: true })
    return () => window.removeEventListener('scroll', onScroll)
  }, [])

  // 觸控裝置沒有 mouseleave，點選單外面時關閉
  useEffect(() => {
    if (open === null) return
    const onPointer = (e: PointerEvent) => {
      if (!navRef.current?.contains(e.target as Node)) setOpen(null)
    }
    document.addEventListener('pointerdown', onPointer)
    return () => document.removeEventListener('pointerdown', onPointer)
  }, [open])

  const focusFirst = (id: number | 'more') =>
    requestAnimationFrame(() => document.querySelector<HTMLElement>(`#menu-${id} a`)?.focus())

  const visible = menu.slice(0, shown)
  const overflow = menu.slice(shown)

  return (
    <header
      className={`sticky top-0 z-40 border-b transition-[border-color,background-color] ${
        compact ? 'border-slate-200 bg-white' : 'border-transparent bg-gradient-to-b from-sky/20 via-sky/5 to-white'
      }`}
    >
      <div className={`mx-auto flex max-w-6xl items-center justify-between gap-4 px-4 transition-[height] ${compact ? 'h-14' : 'h-20'}`}>
        <Link to="/" className="flex min-w-0 items-center gap-3">
          <img src="/logo.png" alt="" className={`w-auto transition-[height] ${compact ? 'h-9' : 'h-12'}`} />
          <span className="truncate font-serif text-2xl font-bold tracking-wider text-brand-dark">{title}</span>
          {subtitle && <span className="hidden font-serif text-sm text-slate-500 xl:inline">{subtitle}</span>}
        </Link>

        <div ref={barRef} data-navbar className="hidden min-w-0 flex-1 justify-end lg:flex">
          <nav ref={navRef} aria-label="主選單">
            <ul className="flex flex-nowrap gap-1">
              {visible.map((g) => (
                <li
                  key={g.id}
                  className="relative shrink-0"
                  onMouseEnter={() => setOpen(g.id)}
                  // relatedTarget 不是實際的節點時代表滑鼠真的離開了整個文件（例如切到別的視窗），
                  // React 會把它標成 window；而不是移到選單內的連結。
                  // @testing-library/user-event 的 click() 在模擬滑鼠移動時不會帶正確的
                  // relatedTarget，若不判斷會在點擊連結前就先關閉選單
                  onMouseLeave={(e) => e.relatedTarget instanceof Node && setOpen(null)}
                >
                  <button
                    type="button"
                    data-active={g.id === activeGroupId || undefined}
                    aria-expanded={g.pages.length > 0 ? open === g.id : undefined}
                    aria-controls={g.pages.length > 0 ? `menu-${g.id}` : undefined}
                    onClick={() => setOpen(g.id)}
                    onKeyDown={(e) => {
                      if (e.key === 'ArrowDown') {
                        e.preventDefault()
                        setOpen(g.id)
                        focusFirst(g.id)
                      }
                    }}
                    className={`${ITEM_CLASS} text-ink hover:text-brand aria-expanded:text-brand data-active:border-olive`}
                  >
                    {g.name}
                  </button>
                  {open === g.id && g.pages.length > 0 && (
                    <MenuPanel id={`menu-${g.id}`} onClose={() => setOpen(null)}>
                      {g.pages.map((p) => (
                        <li key={p.id}>
                          <Link to={`/page/${p.id}`} onClick={() => setOpen(null)} className={LINK_CLASS}>
                            {p.name}
                          </Link>
                        </li>
                      ))}
                    </MenuPanel>
                  )}
                </li>
              ))}

              {overflow.length > 0 && (
                <li
                  className="relative shrink-0"
                  onMouseEnter={() => setOpen('more')}
                  onMouseLeave={(e) => e.relatedTarget instanceof Node && setOpen(null)}
                >
                  <button
                    type="button"
                    data-active={overflow.some((g) => g.id === activeGroupId) || undefined}
                    aria-expanded={open === 'more'}
                    aria-controls="menu-more"
                    onClick={() => setOpen('more')}
                    onKeyDown={(e) => {
                      if (e.key === 'ArrowDown') {
                        e.preventDefault()
                        setOpen('more')
                        focusFirst('more')
                      }
                    }}
                    className={`${ITEM_CLASS} text-ink hover:text-brand aria-expanded:text-brand data-active:border-olive`}
                  >
                    更多
                  </button>
                  {open === 'more' && (
                    <MenuPanel id="menu-more" onClose={() => setOpen(null)} className="max-h-[70vh] overflow-y-auto">
                      {overflow.map((g) => (
                        <li key={g.id}>
                          <p className="px-4 pb-1 pt-3 text-xs font-medium text-slate-500">{g.name}</p>
                          <ul>
                            {g.pages.map((p) => (
                              <li key={p.id}>
                                <Link to={`/page/${p.id}`} onClick={() => setOpen(null)} className={`${LINK_CLASS} pl-7`}>
                                  {p.name}
                                </Link>
                              </li>
                            ))}
                          </ul>
                        </li>
                      ))}
                    </MenuPanel>
                  )}
                </li>
              )}
            </ul>
          </nav>
        </div>

        <button
          type="button"
          aria-label="開啟選單"
          aria-expanded={drawer}
          onClick={() => setDrawer(true)}
          className="rounded-md p-2 text-brand-dark lg:hidden"
        >
          <svg viewBox="0 0 24 24" className="h-7 w-7" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
            <path d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
      </div>
      {/* 量寬度用的隱藏清單：永遠是完整的選單，收折後才還量得到原本的寬度 */}
      <div aria-hidden="true" className="pointer-events-none invisible absolute h-0 overflow-hidden">
        <ul ref={rulerRef} className="flex w-max flex-nowrap gap-1">
          {menu.map((g) => (
            <li key={g.id} data-navitem className={ITEM_CLASS}>
              {g.name}
            </li>
          ))}
          <li data-navitem className={ITEM_CLASS}>
            更多
          </li>
        </ul>
      </div>
      {drawer && <MobileDrawer menu={menu} onClose={() => setDrawer(false)} />}
    </header>
  )
}

function MobileDrawer({ menu, onClose }: { menu: Group[]; onClose: () => void }) {
  const [expanded, setExpanded] = useState<number | null>(null)
  const closeRef = useRef<HTMLButtonElement>(null)

  useEffect(() => {
    closeRef.current?.focus()
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && onClose()
    window.addEventListener('keydown', onKey)
    document.body.style.overflow = 'hidden'
    return () => {
      window.removeEventListener('keydown', onKey)
      document.body.style.overflow = ''
    }
  }, [onClose])

  return (
    <div className="fixed inset-0 z-50 lg:hidden" role="dialog" aria-modal="true" aria-label="選單">
      <div className="absolute inset-0 bg-black/40" onClick={onClose} />
      <nav className="absolute right-0 top-0 flex h-full w-72 max-w-[85vw] flex-col bg-white shadow-xl">
        <div className="flex justify-end p-2">
          <button ref={closeRef} type="button" aria-label="關閉選單" onClick={onClose} className="rounded-md p-2 text-2xl leading-none">
            ×
          </button>
        </div>
        <ul className="overflow-y-auto px-4 pb-8">
          {menu.map((g) => (
            <li key={g.id} className="border-b border-slate-100">
              <button
                type="button"
                aria-expanded={expanded === g.id}
                onClick={() => setExpanded(expanded === g.id ? null : g.id)}
                className="flex w-full items-center justify-between py-3 text-left font-medium text-slate-800"
              >
                {g.name}
                <span aria-hidden="true" className="text-brand">{expanded === g.id ? '−' : '+'}</span>
              </button>
              {expanded === g.id && (
                <ul className="pb-2">
                  {g.pages.map((p) => (
                    <li key={p.id}>
                      <Link to={`/page/${p.id}`} onClick={onClose} className="block py-2 pl-4 text-slate-600">
                        {p.name}
                      </Link>
                    </li>
                  ))}
                </ul>
              )}
            </li>
          ))}
        </ul>
      </nav>
    </div>
  )
}
