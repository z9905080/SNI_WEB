export type PublicConfig = { gaMeasurementId: string; counterScriptUrl: string }

const empty: PublicConfig = { gaMeasurementId: '', counterScriptUrl: '' }

// 設定由 Go 服務注入 index.html；vite dev 時不存在
export function readPublicConfig(doc: Document = document): PublicConfig {
  const text = doc.getElementById('sni-config')?.textContent
  if (!text) return empty
  try {
    return { ...empty, ...JSON.parse(text) }
  } catch {
    return empty
  }
}

export const publicConfig = readPublicConfig()
