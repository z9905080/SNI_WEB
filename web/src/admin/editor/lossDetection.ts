// 視為相同的標籤（Tiptap 會把 b 輸出成 strong 等）
const aliases: Record<string, string> = { b: 'strong', i: 'em', strike: 's', del: 's' }
// 解析或 Tiptap 會自動補上／調整的結構
const ignoredTags = new Set(['tbody', 'colgroup', 'col'])

// 交給瀏覽器正規化 CSS（#e03e2d 與 rgb(224, 62, 45) 視為相同）
function styleDeclarations(style: string): string[] {
  const el = document.createElement('div')
  el.style.cssText = style
  return Array.from(el.style, (prop) => `${prop}: ${el.style.getPropertyValue(prop)}`)
}

function signature(html: string) {
  const doc = new DOMParser().parseFromString(`<body>${html}</body>`, 'text/html')
  const counts = new Map<string, number>()
  const add = (key: string) => counts.set(key, (counts.get(key) ?? 0) + 1)
  for (const el of doc.body.querySelectorAll('*')) {
    const tag = aliases[el.localName] ?? el.localName
    if (ignoredTags.has(tag)) continue
    add(`<${tag}>`)
    for (const attr of el.attributes) {
      if (attr.name === 'style') styleDeclarations(attr.value).forEach(add)
      else if (attr.name === 'class') attr.value.split(/\s+/).filter(Boolean).forEach((c) => add(`class="${c}"`))
      else add(`${attr.name}="${attr.value}"`)
    }
  }
  // 結構容器（body/table/tbody/…）之間的縮排換行是純排版雜訊，直接忽略；
  // 其他容器（p/span/div/td/li/標題……）內的純空白文字節點仍可能是有意義的分隔（例如英數字之間的空白）或
  // 舊站常見的尾端 &nbsp;，要保留，只把連續空白正規化成一個空格。
  const structuralContainers = new Set(['body', 'table', 'tbody', 'thead', 'tfoot', 'tr', 'ul', 'ol', 'colgroup'])
  const walker = document.createTreeWalker(doc.body, NodeFilter.SHOW_TEXT)
  let text = ''
  for (let node = walker.nextNode(); node; node = walker.nextNode()) {
    const value = node.nodeValue ?? ''
    const isWhitespaceOnly = value.trim() === ''
    if (isWhitespaceOnly && structuralContainers.has(node.parentElement?.localName ?? '')) continue
    text += value.replace(/\s+/g, ' ')
  }
  return { counts, text }
}

export function detectLoss(original: string, roundTripped: string): string[] {
  const before = signature(original)
  const after = signature(roundTripped)
  const lost: string[] = []
  for (const [key, n] of before.counts) {
    if ((after.counts.get(key) ?? 0) < n) lost.push(key)
  }
  if (before.text !== after.text) lost.push('文字內容')
  return lost
}
