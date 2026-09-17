import { useCallback, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import type { ImageItem } from '@/shared/types'
import { adminKeys, uploadImage } from '../api'

export type UploadItem = { key: string; name: string; progress: number; error?: string; done?: boolean }

const MAX_BYTES = 5 * 1024 * 1024
const TYPES = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
let seq = 0

export function useUploader() {
  const qc = useQueryClient()
  const [items, setItems] = useState<UploadItem[]>([])

  const upload = useCallback(
    async (files: File[]) => {
      const patch = (key: string, p: Partial<UploadItem>) =>
        setItems((list) => list.map((item) => (item.key === key ? { ...item, ...p } : item)))
      const queue = files.map((file) => ({ file, key: String(++seq) }))
      setItems((list) => [...list, ...queue.map(({ file, key }) => ({ key, name: file.name, progress: 0 }))])

      const saved: ImageItem[] = []
      // 依序上傳：進度清楚，也避免同一秒大量請求
      for (const { file, key } of queue) {
        if (!TYPES.includes(file.type)) {
          patch(key, { error: '只能上傳 jpg、png、gif、webp 圖片' })
          continue
        }
        if (file.size > MAX_BYTES) {
          patch(key, { error: '檔案超過 5MB' })
          continue
        }
        try {
          const image = await uploadImage(file, (progress) => patch(key, { progress }))
          patch(key, { progress: 1, done: true })
          saved.push(image)
        } catch (e) {
          patch(key, { error: (e as Error).message })
        }
      }
      if (saved.length) await qc.invalidateQueries({ queryKey: adminKeys.images })
      return saved
    },
    [qc],
  )

  const clear = useCallback(() => setItems([]), [])
  return { items, upload, clear }
}
