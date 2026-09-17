export class ApiError extends Error {
  status: number
  code: string
  requestId?: string

  constructor(status: number, code: string, message: string, requestId?: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.requestId = requestId
  }
}

type Options = { method?: string; body?: unknown; signal?: AbortSignal }

export async function api<T>(path: string, { method = 'GET', body, signal }: Options = {}): Promise<T> {
  const isForm = body instanceof FormData
  let res: Response
  try {
    res = await fetch(`/api/v1${path}`, {
      method,
      signal,
      credentials: 'same-origin',
      headers: body === undefined || isForm ? undefined : { 'Content-Type': 'application/json' },
      body: body === undefined ? undefined : isForm ? body : JSON.stringify(body),
    })
  } catch (e) {
    if (e instanceof DOMException && e.name === 'AbortError') throw e
    throw new ApiError(0, 'network', '無法連線到伺服器，請檢查網路後再試')
  }
  if (res.status === 204) return undefined as T

  const data = await res.json().catch(() => null)
  if (!res.ok) {
    const err = data?.error
    let message: string = err?.message ?? `伺服器錯誤（${res.status}）`
    if (err?.request_id) message += `（代碼 ${err.request_id}）`
    throw new ApiError(res.status, err?.code ?? 'unknown', message, err?.request_id)
  }
  return data as T
}
