import type { FormEvent } from 'react'
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

type Props = {
  open: boolean
  title: string
  initial: string
  pending: boolean
  onOpenChange: (open: boolean) => void
  onSubmit: (name: string) => void
}

export default function NameDialog({ open, title, initial, pending, onOpenChange, onSubmit }: Props) {
  const submit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    onSubmit(String(new FormData(e.currentTarget).get('name')).trim())
  }
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <form onSubmit={submit} className="space-y-4">
          <DialogHeader>
            <DialogTitle>{title}</DialogTitle>
            <DialogDescription>名稱會顯示在前台選單，最多 50 個字。</DialogDescription>
          </DialogHeader>
          <div className="space-y-2">
            <Label htmlFor="name-dialog-input">名稱</Label>
            <Input id="name-dialog-input" name="name" defaultValue={initial} maxLength={50} required autoFocus />
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              取消
            </Button>
            <Button type="submit" disabled={pending}>
              儲存
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
