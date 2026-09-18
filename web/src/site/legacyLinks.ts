// 舊站使用 hash 路由：/#/{頁面 id}
export function legacyHashTarget(hash: string): string | null {
  if (hash === '#/') return '/'
  const m = /^#\/(\d+)\/?$/.exec(hash)
  return m ? `/page/${m[1]}` : null
}

// internalPath 判斷內文連結是否該用前端路由切換；是的話回傳目標路徑。
// 只比對主機名稱，讓舊內容中的 http:// 絕對網址在 https 站上也能使用。
export function internalPath(href: string, current: { href: string; hostname: string }): string | null {
  if (!href) return null
  let url: URL
  try {
    url = new URL(href, current.href)
  } catch {
    return null
  }
  if (!/^https?:$/.test(url.protocol) || url.hostname !== current.hostname) return null
  if (url.hash) {
    return url.pathname === '/' ? legacyHashTarget(url.hash) : null
  }
  if (url.pathname === '/') return '/'
  return /^\/page\/\d+$/.test(url.pathname) ? url.pathname : null
}
