import { useState } from 'react'
import { Button } from '@/components/ui/button'
import type { UploadItem } from './useUploader'

export function DropZone({ onFiles }: { onFiles: (files: File[]) => void }) {
  const [over, setOver] = useState(false)
  return (
    <label
      onDragOver={(e) => {
        e.preventDefault()
        setOver(true)
      }}
      onDragLeave={() => setOver(false)}
      onDrop={(e) => {
        e.preventDefault()
        setOver(false)
        onFiles([...e.dataTransfer.files])
      }}
      className={`flex cursor-pointer flex-col items-center gap-1 rounded-md border-2 border-dashed p-6 text-center text-sm transition-colors focus-within:ring-2 focus-within:ring-ring ${over ? 'border-primary bg-primary/5' : 'border-border'}`}
    >
      <span className="font-medium">拖曳圖片到這裡，或點擊選擇檔案</span>
      <span className="text-muted-foreground">jpg、png、gif、webp，每個檔案 5MB 以內</span>
      <input
        type="file"
        multiple
        accept="image/jpeg,image/png,image/gif,image/webp"
        aria-label="選擇圖片"
        className="sr-only"
        onChange={(e) => {
          onFiles([...(e.target.files ?? [])])
          e.target.value = ''
        }}
      />
    </label>
  )
}

export function UploadList({ items, onClear }: { items: UploadItem[]; onClear: () => void }) {
  if (items.length === 0) return null
  const busy = items.some((i) => !i.done && !i.error)
  return (
    <div className="space-y-2 rounded-md border bg-background p-3 text-sm">
      <ul className="space-y-2">
        {items.map((item) => (
          <li key={item.key}>
            <div className="flex justify-between gap-2">
              <span className="truncate">{item.name}</span>
              <span className={item.error ? 'text-destructive' : 'text-muted-foreground'}>
                {item.error ?? (item.done ? '已完成' : `${Math.round(item.progress * 100)}%`)}
              </span>
            </div>
            {!item.error && !item.done && (
              <div className="mt-1 h-1.5 overflow-hidden rounded bg-muted" role="progressbar" aria-label={item.name} aria-valuenow={Math.round(item.progress * 100)}>
                <div className="h-full bg-primary transition-[width]" style={{ width: `${item.progress * 100}%` }} />
              </div>
            )}
          </li>
        ))}
      </ul>
      {!busy && (
        <Button type="button" size="sm" variant="ghost" onClick={onClear}>
          清除上傳紀錄
        </Button>
      )}
    </div>
  )
}
