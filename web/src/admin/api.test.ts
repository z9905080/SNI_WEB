import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '@/shared/api'
import { adminApi, uploadImage } from './api'

afterEach(() => vi.unstubAllGlobals())

describe('adminApi', () => {
  it('解開回應外層並送出正確的請求', async () => {
    const fetchMock = vi.fn(async () =>
      new Response(JSON.stringify({ page: { id: 5, group_id: 2, name: 'x', html: '' } }), { status: 201 }),
    )
    vi.stubGlobal('fetch', fetchMock)
    const page = await adminApi.createPage({ name: 'x', group_id: 2, html: '' })
    expect(page.id).toBe(5)
    expect(fetchMock).toHaveBeenCalledWith('/api/v1/admin/pages', expect.objectContaining({ method: 'POST' }))
  })

  it('圖片名稱會編碼', async () => {
    // 型別上明確標出 url 參數，讓 fetchMock.mock.calls[0][0] 能推得出正確型別
    const fetchMock = vi.fn(async (_url: string) => new Response(null, { status: 204 }))
    vi.stubGlobal('fetch', fetchMock)
    await adminApi.deleteImage('a b.jpg')
    expect(fetchMock.mock.calls[0][0]).toBe('/api/v1/admin/images/a%20b.jpg')
  })
})

class FakeXHR {
  static last: FakeXHR
  upload = { onprogress: null as ((e: { lengthComputable: boolean; loaded: number; total: number }) => void) | null }
  onload: (() => void) | null = null
  onerror: (() => void) | null = null
  status = 0
  responseText = ''
  method = ''
  url = ''
  body: FormData | null = null
  constructor() {
    FakeXHR.last = this
  }
  open(method: string, url: string) {
    this.method = method
    this.url = url
  }
  send(body: FormData) {
    this.body = body
  }
}

describe('uploadImage', () => {
  it('回報進度並回傳上傳結果', async () => {
    vi.stubGlobal('XMLHttpRequest', FakeXHR)
    const progress: number[] = []
    const file = new File(['x'], 'a.png', { type: 'image/png' })
    const promise = uploadImage(file, (r) => progress.push(r))
    const xhr = FakeXHR.last
    expect(xhr.method).toBe('POST')
    expect(xhr.url).toBe('/api/v1/admin/images')
    expect(xhr.body?.get('files')).toBe(file)
    xhr.upload.onprogress?.({ lengthComputable: true, loaded: 50, total: 100 })
    xhr.status = 201
    xhr.responseText = JSON.stringify({ images: [{ name: 'n.png', url: '/php/picture/n.png', size: 1, mod_time: '' }] })
    xhr.onload?.()
    await expect(promise).resolves.toMatchObject({ name: 'n.png' })
    expect(progress).toEqual([0.5, 1])
  })

  it('錯誤回應轉成 ApiError', async () => {
    vi.stubGlobal('XMLHttpRequest', FakeXHR)
    const promise = uploadImage(new File(['x'], 'a.svg'), () => {})
    const xhr = FakeXHR.last
    xhr.status = 400
    xhr.responseText = JSON.stringify({ error: { code: 'unsupported_type', message: '不是支援的圖片格式' } })
    xhr.onload?.()
    await expect(promise).rejects.toMatchObject({ status: 400, code: 'unsupported_type', message: '不是支援的圖片格式' })
    await expect(promise).rejects.toBeInstanceOf(ApiError)
  })

  it('網路錯誤', async () => {
    vi.stubGlobal('XMLHttpRequest', FakeXHR)
    const promise = uploadImage(new File(['x'], 'a.png'), () => {})
    FakeXHR.last.onerror?.()
    await expect(promise).rejects.toMatchObject({ status: 0, code: 'network' })
  })
})
