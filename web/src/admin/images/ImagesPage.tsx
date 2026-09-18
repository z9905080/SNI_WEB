import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import type { ImageItem, ImageUsages } from '@/shared/types'
import { useDocumentTitle } from '@/site/useDocumentTitle'
import { adminApi, adminKeys } from '../api'
import ConfirmDialog from '../ConfirmDialog'
import { useAction } from '../useAction'
import { useUploader } from './useUploader'
import { DropZone, UploadList } from './UploadWidgets'

const formatSize = (bytes: number) =>
  bytes >= 1024 * 1024 ? `${(bytes / 1024 / 1024).toFixed(1)} MB` : `${Math.ceil(bytes / 1024)} KB`

export default function ImagesPage() {
  useDocumentTitle('圖庫｜後台')
  const images = useQuery({ queryKey: adminKeys.images, queryFn: adminApi.images })
  const uploader = useUploader()
  const action = useAction(adminKeys.images)
  const [target, setTarget] = useState<{ image: ImageItem; usages: ImageUsages } | null>(null)

  const askDelete = async (image: ImageItem) => {
    try {
      setTarget({ image, usages: await adminApi.imageUsages(image.name) })
    } catch (e) {
      toast.error((e as Error).message)
    }
  }

  const copyUrl = async (image: ImageItem) => {
    try {
      await navigator.clipboard.writeText(new URL(image.url, window.location.origin).href)
      toast.success('已複製網址')
    } catch {
      toast.error('無法存取剪貼簿，請手動複製網址')
    }
  }

  const used = target && (target.usages.pages.length > 0 || target.usages.carousels.length > 0)

  return (
    <div className="space-y-6">
      <h1 className="font-serif text-2xl font-bold">圖庫</h1>
      <DropZone onFiles={uploader.upload} />
      <UploadList items={uploader.items} onClear={uploader.clear} />

      {images.isError && <p role="alert" className="text-destructive">{images.error.message}</p>}
      {images.isPending && <Skeleton className="h-48" />}
      {images.data?.length === 0 && <p className="text-muted-foreground">圖庫還沒有圖片，請從上方上傳。</p>}
      <ul className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
        {images.data?.map((image) => (
          <li key={image.name} className="overflow-hidden rounded-md border bg-background">
            <img src={image.url} alt={image.name} loading="lazy" className="aspect-square w-full bg-muted object-cover" />
            <div className="space-y-2 p-2 text-xs">
              <p className="truncate" title={image.name}>{image.name}</p>
              <p className="text-muted-foreground">{formatSize(image.size)}</p>
              <div className="flex gap-1">
                <Button size="sm" variant="outline" onClick={() => copyUrl(image)}>
                  複製網址
                </Button>
                <Button size="sm" variant="ghost" aria-label={`刪除 ${image.name}`} onClick={() => askDelete(image)}>
                  刪除
                </Button>
              </div>
            </div>
          </li>
        ))}
      </ul>

      <ConfirmDialog
        open={target !== null}
        onOpenChange={(open) => !open && setTarget(null)}
        title={used ? '這張圖片仍在使用中' : '刪除圖片'}
        description={
          used ? (
            <div className="space-y-2">
              <p>刪除後，下列位置將無法顯示這張圖片：</p>
              <ul className="list-disc pl-5">
                {target!.usages.pages.map((p) => (
                  <li key={`p${p.id}`}>
                    頁面：<Link to={`/admin/pages/${p.id}`} className="underline">{p.name}</Link>
                  </li>
                ))}
                {target!.usages.carousels.map((c) => (
                  <li key={`c${c.id}`}>輪播圖 #{c.id}</li>
                ))}
              </ul>
            </div>
          ) : (
            `確定要刪除「${target?.image.name}」嗎？此動作無法復原。`
          )
        }
        confirmLabel={used ? '仍要刪除' : '刪除'}
        destructive
        onConfirm={() => {
          const name = target!.image.name
          setTarget(null)
          action.mutate(() => adminApi.deleteImage(name), { onSuccess: () => toast.success('已刪除圖片') })
        }}
      />
    </div>
  )
}
