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
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { Textarea } from '@/components/ui/textarea'
import type { Marquee } from '@/shared/types'
import { useDocumentTitle } from '@/site/useDocumentTitle'
import { adminApi, adminKeys, type MarqueeInput } from '../api'
import ConfirmDialog from '../ConfirmDialog'
import { useAction } from '../useAction'

export default function MarqueesPage() {
  useDocumentTitle('跑馬燈｜後台')
  const list = useQuery({ queryKey: adminKeys.marquees, queryFn: adminApi.marquees })
  const action = useAction(adminKeys.marquees)
  const [editing, setEditing] = useState<Marquee | 'new' | null>(null)
  const [deleting, setDeleting] = useState<Marquee | null>(null)

  const save = (input: MarqueeInput) => {
    const target = editing
    action.mutate(() => (target === 'new' || !target ? adminApi.createMarquee(input) : adminApi.updateMarquee(target.id, input)), {
      onSuccess: () => {
        setEditing(null)
        toast.success('已儲存跑馬燈')
      },
    })
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="font-serif text-2xl font-bold">跑馬燈</h1>
        <Button onClick={() => setEditing('new')}>新增跑馬燈</Button>
      </div>
      {list.isError && <p role="alert" className="text-destructive">{list.error.message}</p>}
      {list.isPending && <Skeleton className="h-32" />}
      {list.data?.length === 0 && <p className="text-muted-foreground">目前沒有跑馬燈，前台不會顯示公告列。</p>}
      <ul className="divide-y rounded-md border bg-background">
        {list.data?.map((m) => (
          <li key={m.id} className="flex items-center gap-3 p-3">
            <span aria-hidden="true" className="h-5 w-5 shrink-0 rounded border" style={{ backgroundColor: m.color }} />
            <span className="min-w-0 flex-1 truncate">{m.text}</span>
            <Button size="sm" variant="ghost" aria-label={`編輯 ${m.text}`} onClick={() => setEditing(m)}>編輯</Button>
            <Button size="sm" variant="ghost" aria-label={`刪除 ${m.text}`} onClick={() => setDeleting(m)}>刪除</Button>
          </li>
        ))}
      </ul>

      {editing && (
        <MarqueeDialog
          initial={editing === 'new' ? { text: '', color: '#ffffff' } : editing}
          pending={action.isPending}
          onClose={() => setEditing(null)}
          onSave={save}
        />
      )}
      <ConfirmDialog
        open={deleting !== null}
        onOpenChange={(open) => !open && setDeleting(null)}
        title="刪除跑馬燈"
        description={`確定要刪除「${deleting?.text}」嗎？`}
        confirmLabel="刪除"
        destructive
        onConfirm={() => {
          const id = deleting!.id
          setDeleting(null)
          action.mutate(() => adminApi.deleteMarquee(id), { onSuccess: () => toast.success('已刪除跑馬燈') })
        }}
      />
    </div>
  )
}

type DialogProps = { initial: MarqueeInput; pending: boolean; onClose: () => void; onSave: (v: MarqueeInput) => void }

function MarqueeDialog({ initial, pending, onClose, onSave }: DialogProps) {
  const [text, setText] = useState(initial.text)
  const [color, setColor] = useState(initial.color)
  return (
    <Dialog open onOpenChange={(open) => !open && onClose()}>
      <DialogContent>
        <form
          className="space-y-4"
          onSubmit={(e) => {
            e.preventDefault()
            onSave({ text: text.trim(), color })
          }}
        >
          <DialogHeader>
            <DialogTitle>{initial.text ? '編輯跑馬燈' : '新增跑馬燈'}</DialogTitle>
            <DialogDescription>公告會在前台輪播圖下方持續捲動。</DialogDescription>
          </DialogHeader>
          <div className="space-y-2">
            <Label htmlFor="marquee-text">文字</Label>
            <Textarea id="marquee-text" value={text} maxLength={500} required onChange={(e) => setText(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="marquee-color">顏色</Label>
            <div className="flex items-center gap-2">
              <input id="marquee-color" type="color" value={color} onChange={(e) => setColor(e.target.value)} className="h-9 w-14 cursor-pointer rounded border" />
              <span className="font-mono text-sm text-muted-foreground">{color.toUpperCase()}</span>
            </div>
          </div>
          <div className="rounded-md bg-brand-dark px-4 py-2" aria-label="預覽">
            <span data-testid="marquee-preview" style={{ color }}>
              {text || '公告文字會顯示在這裡'}
            </span>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={onClose}>取消</Button>
            <Button type="submit" disabled={pending}>儲存</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
