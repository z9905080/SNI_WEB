import { afterEach, describe, expect, it, vi } from 'vitest'
import { api, ApiError } from './api'

function mockFetch(impl: (...args: Parameters<typeof fetch>) => Promise<Response>) {
  const fn = vi.fn(impl)
  vi.stubGlobal('fetch', fn)
  return fn
}

const json = (status: number, body: unknown) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } })

afterEach(() => vi.unstubAllGlobals())

describe('api', () => {
  it('送出 JSON 並解析回應', async () => {
    const fetchMock = mockFetch(async () => json(200, { ok: true }))
    await expect(api('/admin/groups', { method: 'POST', body: { name: 'x' } })).resolves.toEqual({ ok: true })
    const [url, init] = fetchMock.mock.calls[0]
    expect(url).toBe('/api/v1/admin/groups')
    expect(init?.method).toBe('POST')
    expect(init?.body).toBe('{"name":"x"}')
    expect(init?.credentials).toBe('same-origin')
    expect(new Headers(init?.headers).get('Content-Type')).toBe('application/json')
  })

  it('FormData 原樣送出且不設定 Content-Type', async () => {
    const fetchMock = mockFetch(async () => json(201, {}))
    const form = new FormData()
    await api('/admin/images', { method: 'POST', body: form })
    const init = fetchMock.mock.calls[0][1]
    expect(init?.body).toBe(form)
    expect(new Headers(init?.headers).has('Content-Type')).toBe(false)
  })

  it('204 回傳 undefined', async () => {
    mockFetch(async () => new Response(null, { status: 204 }))
    await expect(api('/admin/pages/1', { method: 'DELETE' })).resolves.toBeUndefined()
  })

  it('錯誤回應轉成 ApiError', async () => {
    mockFetch(async () => json(409, { error: { code: 'group_not_empty', message: '此頁籤仍有頁面' } }))
    const err = await api('/x').catch((e) => e)
    expect(err).toBeInstanceOf(ApiError)
    expect(err).toMatchObject({ status: 409, code: 'group_not_empty', message: '此頁籤仍有頁面' })
  })

  it('500 訊息附上錯誤代碼', async () => {
    mockFetch(async () => json(500, { error: { code: 'internal', message: '系統發生錯誤', request_id: 'abc123' } }))
    const err = (await api('/x').catch((e) => e)) as ApiError
    expect(err.message).toBe('系統發生錯誤（代碼 abc123）')
    expect(err.requestId).toBe('abc123')
  })

  it('非 JSON 錯誤回應也有可讀訊息', async () => {
    mockFetch(async () => new Response('bad gateway', { status: 502 }))
    const err = await api('/x').catch((e) => e)
    expect(err).toMatchObject({ status: 502, code: 'unknown', message: '伺服器錯誤（502）' })
  })

  it('網路錯誤', async () => {
    mockFetch(async () => {
      throw new TypeError('Failed to fetch')
    })
    const err = await api('/x').catch((e) => e)
    expect(err).toMatchObject({ status: 0, code: 'network' })
  })

  it('取消請求時保留 AbortError', async () => {
    mockFetch(async () => {
      throw new DOMException('aborted', 'AbortError')
    })
    const err = (await api('/x').catch((e) => e)) as DOMException
    expect(err.name).toBe('AbortError')
  })
})
