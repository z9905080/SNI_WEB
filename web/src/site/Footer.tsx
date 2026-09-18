type Props = { title: string; facebookUrl: string; counterScriptUrl: string }

const escapeAttr = (s: string) =>
  s.replaceAll('&', '&amp;').replaceAll('"', '&quot;').replaceAll('<', '&lt;').replaceAll('>', '&gt;')

export default function Footer({ title, facebookUrl, counterScriptUrl }: Props) {
  return (
    <footer className="mt-20 bg-brand-dark text-slate-200">
      <div className="mx-auto flex max-w-6xl flex-col gap-6 px-4 py-12 text-sm sm:flex-row sm:items-end sm:justify-between">
        <div className="space-y-2">
          {title && <p className="font-serif text-xl font-bold tracking-wider text-white">{title}</p>}
          <p>© {new Date().getFullYear()} Seicho-No-Ie R.O.C Missionary Headquarters All Rights Reserved</p>
        </div>
        <div className="flex flex-col items-start gap-3 sm:items-end">
          {facebookUrl && (
            <a
              href={facebookUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 underline-offset-4 hover:underline"
            >
              <svg viewBox="0 0 24 24" className="h-4 w-4" fill="currentColor" aria-hidden="true">
                <path d="M14 8h3V4h-3c-2.8 0-4 1.7-4 4.3V10H7v4h3v8h4v-8h3l1-4h-4V8.5c0-.3.2-.5.5-.5Z" />
              </svg>
              Facebook
            </a>
          )}
          {counterScriptUrl && (
            // 舊式計數器以 document.write 輸出，放在沙箱 iframe 內執行
            <iframe
              title="瀏覽次數"
              sandbox="allow-scripts"
              className="h-6 w-56 border-0"
              srcDoc={`<!doctype html><meta charset="utf-8"><body style="margin:0;font:13px sans-serif;color:#e2e8f0;white-space:nowrap">本站瀏覽次數：<script src="${escapeAttr(counterScriptUrl)}"></script></body>`}
            />
          )}
        </div>
      </div>
    </footer>
  )
}
