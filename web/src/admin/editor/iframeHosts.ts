// 需與後端 backend/internal/content/sanitize.go 的 IframeHosts 一致
export const IFRAME_HOSTS = [
  'www.youtube.com',
  'youtube.com',
  'www.youtube-nocookie.com',
  'players.brightcove.net',
  'drive.google.com',
  'www.facebook.com',
]

export function isAllowedIframe(src: string): boolean {
  try {
    // 相對網址會解析到 invalid 主機，自然不在白名單內
    const url = new URL(src, 'https://relative.invalid')
    return /^https?:$/.test(url.protocol) && IFRAME_HOSTS.includes(url.hostname)
  } catch {
    return false
  }
}

export function toEmbedUrl(input: string): string {
  const trimmed = input.trim()
  const raw = /<iframe[^>]*\ssrc=["']([^"']+)["']/i.exec(trimmed)?.[1] ?? trimmed
  let url: URL
  try {
    url = new URL(raw)
  } catch {
    return raw
  }
  const host = url.hostname.replace(/^(www|m)\./, '')
  if (host === 'youtube.com' && url.pathname === '/watch' && url.searchParams.get('v')) {
    return `https://www.youtube.com/embed/${url.searchParams.get('v')}`
  }
  if (host === 'youtube.com' && url.pathname.startsWith('/shorts/')) {
    return `https://www.youtube.com/embed/${url.pathname.slice('/shorts/'.length)}`
  }
  if (host === 'youtu.be') return `https://www.youtube.com/embed${url.pathname}`
  return raw
}
