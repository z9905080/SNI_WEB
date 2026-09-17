import { useQuery } from '@tanstack/react-query'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { adminApi, adminKeys } from '../api'
import { useUploader } from './useUploader'
import { DropZone, UploadList } from './UploadWidgets'

type Props = { open: boolean; onOpenChange: (open: boolean) => void; onSelect: (url: string) => void }

export default function ImagePickerDialog({ open, onOpenChange, onSelect }: Props) {
  const images = useQuery({ queryKey: adminKeys.images, queryFn: adminApi.images, enabled: open })
  const uploader = useUploader()

  const uploadAndSelect = async (files: File[]) => {
    const saved = await uploader.upload(files)
    if (saved.length === 1) onSelect(saved[0].url)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-3xl">
        <DialogHeader>
          <DialogTitle>選擇圖片</DialogTitle>
          <DialogDescription>點選圖庫中的圖片，或上傳一張新圖片直接使用。</DialogDescription>
        </DialogHeader>
        <DropZone onFiles={uploadAndSelect} />
        <UploadList items={uploader.items} onClear={uploader.clear} />
        {images.isError && <p role="alert" className="text-sm text-destructive">{images.error.message}</p>}
        <ul className="grid max-h-[50vh] grid-cols-3 gap-2 overflow-y-auto sm:grid-cols-4">
          {images.data?.map((image) => (
            <li key={image.name}>
              <button
                type="button"
                onClick={() => onSelect(image.url)}
                className="block w-full overflow-hidden rounded border hover:ring-2 hover:ring-primary focus-visible:ring-2 focus-visible:ring-primary"
              >
                <img src={image.url} alt={image.name} loading="lazy" className="aspect-square w-full bg-muted object-cover" />
              </button>
            </li>
          ))}
        </ul>
      </DialogContent>
    </Dialog>
  )
}
