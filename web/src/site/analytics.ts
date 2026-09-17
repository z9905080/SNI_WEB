declare global {
  interface Window {
    dataLayer?: unknown[]
    gtag?: (...args: unknown[]) => void
  }
}

let loaded = false

export function initAnalytics(id: string) {
  if (!id || loaded) return
  loaded = true
  window.dataLayer = window.dataLayer ?? []
  // gtag.js 需要 push arguments 物件，不能改成陣列
  window.gtag = function gtag() {
    window.dataLayer!.push(arguments)
  }
  window.gtag('js', new Date())
  window.gtag('config', id, { send_page_view: false })
  const script = document.createElement('script')
  script.async = true
  script.src = `https://www.googletagmanager.com/gtag/js?id=${encodeURIComponent(id)}`
  document.head.append(script)
}

export function trackPageView(path: string) {
  window.gtag?.('event', 'page_view', {
    page_path: path,
    page_location: window.location.origin + path,
    page_title: document.title,
  })
}
