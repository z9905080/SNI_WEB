import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import type { Carousel } from '@/shared/types'
import { useDocumentTitle } from '@/site/useDocumentTitle'
import { adminApi, adminKeys, type CarouselInput } from '../api'
import ConfirmDialog from '../ConfirmDialog'
import ImagePickerDialog from '../images/ImagePickerDialog'
import { useAction } from '../useAction'

export default function CarouselsPage() {
  useDocumentTitle('輪播圖｜後台')
  const list = useQuery({ queryKey: adminKeys.carousels, queryFn: adminApi.carousels })
  const action = useAction(adminKeys.carousels)
  const [editing, setEditing] = useState<Carousel | 'new' | null>(null)
  const [deleting, setDeleting] = useState<Carousel | null>(null)

  const save = (input: CarouselInput) => {
    const target = editing
    action.mutate(() => (target === 'new' || !target ? adminApi.createCarousel(input) : adminApi.updateCarousel(target.id, input)), {
      onSuccess: () => {
        setEditing(null)
        toast.success('已儲存輪播圖')
      },
    })
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="font-serif text-2xl font-bold">輪播圖</h1>
        <Button onClick={() => setEditing('new')}>新增輪播圖</Button>
      </div>
      <p className="text-sm text-muted-foreground">依建立順序播放。建議使用寬高比 16:7 的圖片，重要文字請放在畫面中央。</p>
      {list.isError && <p role="alert" className="text-destructive">{list.error.message}</p>}
      {list.isPending && <Skeleton className="h-48" />}
      {list.data?.length === 0 && <p className="text-muted-foreground">目前沒有輪播圖，前台不會顯示輪播區。</p>}
      <ul className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        {list.data?.map((c) => (
          <li key={c.id} className="overflow-hidden rounded-md border bg-background">
            <img src={c.image} alt="" className="aspect-[16/7] w-full bg-muted object-cover" />
            <div className="flex items-center gap-2 p-3 text-sm">
              <span className="min-w-0 flex-1 truncate text-muted-foreground" title={c.url}>
                {c.url || '沒有連結'}
              </span>
              <Button size="sm" variant="ghost" aria-label={`編輯輪播圖 #${c.id}`} onClick={() => setEditing(c)}>編輯</Button>
              <Button size="sm" variant="ghost" aria-label={`刪除輪播圖 #${c.id}`} onClick={() => setDeleting(c)}>刪除</Button>
            </div>
          </li>
        ))}
      </ul>

      {editing && (
        <CarouselDialog
          initial={editing === 'new' ? { image: '', url: '' } : editing}
          isNew={editing === 'new'}
          pending={action.isPending}
          onClose={() => setEditing(null)}
          onSave={save}
        />
      )}
      <ConfirmDialog
        open={deleting !== null}
        onOpenChange={(open) => !open && setDeleting(null)}
        title="刪除輪播圖"
        description="圖片檔案仍會留在圖庫，只會從輪播中移除。"
        confirmLabel="刪除"
        destructive
        onConfirm={() => {
          const id = deleting!.id
          setDeleting(null)
          action.mutate(() => adminApi.deleteCarousel(id), { onSuccess: () => toast.success('已刪除輪播圖') })
        }}
      />
    </div>
  )
}

type DialogProps = {
  initial: CarouselInput
  isNew: boolean
  pending: boolean
  onClose: () => void
  onSave: (v: CarouselInput) => void
}

function CarouselDialog({ initial, isNew, pending, onClose, onSave }: DialogProps) {
  const [image, setImage] = useState(initial.image)
  const [url, setUrl] = useState(initial.url)
  const [picker, setPicker] = useState(false)
  return (
    <Dialog open onOpenChange={(open) => !open && onClose()}>
      <DialogContent>
        <form
          className="space-y-4"
          onSubmit={(e) => {
            e.preventDefault()
            onSave({ image, url: url.trim() })
          }}
        >
          <DialogHeader>
            <DialogTitle>{isNew ? '新增輪播圖' : '編輯輪播圖'}</DialogTitle>
            <DialogDescription>從圖庫選擇圖片，並可設定點擊後前往的網址。</DialogDescription>
          </DialogHeader>
          <div className="space-y-2">
            {image ? (
              <img src={image} alt="已選擇的圖片" className="aspect-[16/7] w-full rounded-md bg-muted object-cover" />
            ) : (
              <div className="flex aspect-[16/7] items-center justify-center rounded-md border-2 border-dashed text-sm text-muted-foreground">
                尚未選擇圖片
              </div>
            )}
            <Button type="button" variant="outline" onClick={() => setPicker(true)}>
              {image ? '更換圖片' : '選擇圖片'}
            </Button>
          </div>
          <div className="space-y-2">
            <Label htmlFor="carousel-url">連結網址</Label>
            <Input id="carousel-url" value={url} placeholder="https://… 或 /page/3，可留空" onChange={(e) => setUrl(e.target.value)} />
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={onClose}>取消</Button>
            <Button type="submit" disabled={pending || !image}>儲存</Button>
          </DialogFooter>
        </form>
        <ImagePickerDialog
          open={picker}
          onOpenChange={setPicker}
          onSelect={(selected) => {
            setImage(selected)
            setPicker(false)
          }}
        />
      </DialogContent>
    </Dialog>
  )
}
