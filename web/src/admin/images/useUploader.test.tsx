import type { ReactNode } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { act, renderHook } from '@testing-library/react'
import { expect, it, vi } from 'vitest'
import { ApiError } from '@/shared/api'
import { useUploader } from './useUploader'

const uploadImage = vi.hoisted(() => vi.fn())
vi.mock('../api', async (importOriginal) => ({ ...(await importOriginal<object>()), uploadImage }))

const wrapper = ({ children }: { children: ReactNode }) => (
  <QueryClientProvider client={new QueryClient()}>{children}</QueryClientProvider>
)

const file = (name: string, type: string, size = 10) => new File([new Uint8Array(size)], name, { type })

it('依序上傳、回報進度，並在前端先擋下不合格的檔案', async () => {
  uploadImage.mockImplementation(async (f: File, onProgress: (r: number) => void) => {
    if (f.name === 'server-reject.png') throw new ApiError(400, 'unsupported_type', '不是支援的圖片格式')
    onProgress(0.5)
    return { name: `saved-${f.name}`, url: `/php/picture/saved-${f.name}`, size: f.size, mod_time: '' }
  })
  const { result } = renderHook(() => useUploader(), { wrapper })

  let done: Awaited<ReturnType<typeof result.current.upload>> = []
  await act(async () => {
    done = await result.current.upload([
      file('a.png', 'image/png'),
      file('big.jpg', 'image/jpeg', 5 * 1024 * 1024 + 1),
      file('doc.pdf', 'application/pdf'),
      file('server-reject.png', 'image/png'),
    ])
  })

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
