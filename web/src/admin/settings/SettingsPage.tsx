import type { FormEvent } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import type { Settings } from '@/shared/types'
import { useDocumentTitle } from '@/site/useDocumentTitle'
import { adminApi, adminKeys } from '../api'

export default function SettingsPage() {
  useDocumentTitle('網站設定｜後台')
  const query = useQuery({ queryKey: adminKeys.settings, queryFn: adminApi.settings })
  if (query.isError) return <p role="alert" className="text-destructive">{query.error.message}</p>
  if (!query.data) return <Skeleton className="h-64" />
  return <SettingsForm initial={query.data} />
}

function SettingsForm({ initial }: { initial: Settings }) {
  const qc = useQueryClient()
  const save = useMutation({
    mutationFn: adminApi.saveSettings,
    onSuccess: (saved) => {
      qc.setQueryData(adminKeys.settings, saved)
      toast.success('已儲存設定')
    },
    onError: (e) => toast.error(e.message),
  })

  const onSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = new FormData(e.currentTarget)
    save.mutate({
      web_title: String(form.get('web_title')).trim(),
      web_sub_title: String(form.get('web_sub_title')).trim(),
      facebook_url: String(form.get('facebook_url')).trim(),
    })
  }

  return (
    <form onSubmit={onSubmit} className="max-w-xl space-y-6">
      <h1 className="font-serif text-2xl font-bold">網站設定</h1>
      <div className="space-y-2">
        <Label htmlFor="web_title">網站名稱</Label>
        <Input id="web_title" name="web_title" defaultValue={initial.web_title} maxLength={100} required />
        <p className="text-sm text-muted-foreground">顯示在頁首、瀏覽器分頁與搜尋結果。</p>
      </div>
      <div className="space-y-2">
        <Label htmlFor="web_sub_title">副標題</Label>
        <Input id="web_sub_title" name="web_sub_title" defaultValue={initial.web_sub_title} maxLength={200} />
        <p className="text-sm text-muted-foreground">首頁會以大字顯示這一句話。</p>
      </div>
      <div className="space-y-2">
        <Label htmlFor="facebook_url">Facebook 粉絲專頁網址</Label>
        <Input id="facebook_url" name="facebook_url" type="url" defaultValue={initial.facebook_url} placeholder="https://www.facebook.com/…" />
        <p className="text-sm text-muted-foreground">留空則頁尾不顯示 Facebook 連結。</p>
      </div>
      <Button type="submit" disabled={save.isPending}>
        {save.isPending ? '儲存中…' : '儲存設定'}
      </Button>
    </form>
  )
}
