import { useState } from 'react'
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
import { IFRAME_HOSTS, isAllowedIframe, toEmbedUrl } from './iframeHosts'

type Props = { open: boolean; onOpenChange: (open: boolean) => void; onInsert: (src: string) => void }

export default function VideoDialog({ open, onOpenChange, onInsert }: Props) {
  const [value, setValue] = useState('')
  const src = toEmbedUrl(value)
  const allowed = isAllowedIframe(src)

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <form
          className="space-y-4"
          onSubmit={(e) => {
            e.preventDefault()
            if (!allowed) return
            onInsert(src)
            setValue('')
          }}
        >
          <DialogHeader>
            <DialogTitle>插入影片</DialogTitle>
            <DialogDescription>貼上 YouTube 網址，或允許網域的嵌入網址、整段嵌入碼。</DialogDescription>
          </DialogHeader>
          <Input aria-label="影片網址" value={value} onChange={(e) => setValue(e.target.value)} placeholder="https://www.youtube.com/watch?v=…" autoFocus />
          {value.trim() !== '' && !allowed && (
            <p role="alert" className="text-sm text-destructive">
              無法嵌入這個網址。可嵌入的網域：{IFRAME_HOSTS.join('、')}
            </p>
          )}
          <DialogFooter>
            <Button type="submit" disabled={!allowed}>插入影片</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
