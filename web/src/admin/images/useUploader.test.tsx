import type { ReactNode } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { act, renderHook, screen } from '@testing-library/react'
import { expect, it, vi } from 'vitest'
import { Toaster } from '@/components/ui/sonner'
import { ApiError } from '@/shared/api'
import { useUploader } from './useUploader'

const uploadImage = vi.hoisted(() => vi.fn())
vi.mock('../api', async (importOriginal) => ({ ...(await importOriginal<object>()), uploadImage }))

const wrapper = ({ children }: { children: ReactNode }) => (
  <QueryClientProvider client={new QueryClient()}>{children}</QueryClientProvider>
)

const withToaster = ({ children }: { children: ReactNode }) => (
  <QueryClientProvider client={new QueryClient()}>
    {children}
    <Toaster />
  </QueryClientProvider>
)

const file = (name: string, type: string, size = 10) => new File([new Uint8Array(size)], name, { type })

it('依序上傳、回報進度，並在前端先擋下不合格的檔案', async () => {
  const order: string[] = []
  let resolveA!: () => void
  const aGate = new Promise<void>((resolve) => {
    resolveA = resolve
  })

  uploadImage.mockImplementation(async (f: File, onProgress: (r: number) => void) => {
    order.push(`start:${f.name}`)
    if (f.name === 'a.png') await aGate
    if (f.name === 'server-reject.png') throw new ApiError(400, 'unsupported_type', '不是支援的圖片格式')
    onProgress(0.5)
    order.push(`end:${f.name}`)
    return { name: `saved-${f.name}`, url: `/php/picture/saved-${f.name}`, size: f.size, mod_time: '' }
  })
  const { result } = renderHook(() => useUploader(), { wrapper })

  let done: Awaited<ReturnType<typeof result.current.upload>> = []
  let uploadPromise!: ReturnType<typeof result.current.upload>
  act(() => {
    uploadPromise = result.current.upload([
      file('a.png', 'image/png'),
      file('big.jpg', 'image/jpeg', 5 * 1024 * 1024 + 1),
      file('doc.pdf', 'application/pdf'),
      file('server-reject.png', 'image/png'),
    ])
  })

  // a.png 的上傳被 aGate 卡住：此時只有第一個檔案開始呼叫 uploadImage，
  // 若實作改成 Promise.all（平行上傳），server-reject.png 這時就會已經開始，斷言會失敗。
  await act(async () => {
    await Promise.resolve()
    await Promise.resolve()
  })
  expect(order).toEqual(['start:a.png'])
  expect(uploadImage).toHaveBeenCalledTimes(1)

  await act(async () => {
    resolveA()
    done = await uploadPromise
  })

  // a.png 結束後才輪到 server-reject.png 開始，證明是依序（非平行）上傳。
  expect(order).toEqual(['start:a.png', 'end:a.png', 'start:server-reject.png'])

  expect(done.map((d) => d.name)).toEqual(['saved-a.png'])
  expect(uploadImage).toHaveBeenCalledTimes(2)
  expect(result.current.items).toEqual([
    expect.objectContaining({ name: 'a.png', progress: 1, done: true }),
    expect.objectContaining({ name: 'big.jpg', error: '檔案超過 5MB' }),
    expect.objectContaining({ name: 'doc.pdf', error: '只能上傳 jpg、png、gif、webp 圖片' }),
    expect.objectContaining({ name: 'server-reject.png', error: '不是支援的圖片格式' }),
  ])

  act(() => result.current.clear())
  expect(result.current.items).toEqual([])
})

it('全部上傳成功時顯示一則彙總成功 toast', async () => {
  uploadImage.mockImplementation(async (f: File, onProgress: (r: number) => void) => {
    onProgress(1)
    return { name: `saved-${f.name}`, url: `/php/picture/saved-${f.name}`, size: f.size, mod_time: '' }
  })
  const { result } = renderHook(() => useUploader(), { wrapper: withToaster })

  await act(async () => {
    await result.current.upload([file('a.png', 'image/png'), file('b.png', 'image/png')])
  })

  expect(await screen.findByText('已上傳 2 張圖片')).toBeInTheDocument()
})

it('有檔案上傳失敗時顯示一則彙總失敗 toast', async () => {
  uploadImage.mockImplementation(async (f: File, onProgress: (r: number) => void) => {
    onProgress(1)
    return { name: `saved-${f.name}`, url: `/php/picture/saved-${f.name}`, size: f.size, mod_time: '' }
  })
  const { result } = renderHook(() => useUploader(), { wrapper: withToaster })

  await act(async () => {
    await result.current.upload([file('a.png', 'image/png'), file('doc.pdf', 'application/pdf')])
  })

  expect(await screen.findByText('1 張圖片上傳失敗')).toBeInTheDocument()
})
