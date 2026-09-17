import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link, useSearchParams } from 'react-router'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import type { Group, PageRef } from '@/shared/types'
import { useDocumentTitle } from '@/site/useDocumentTitle'
import { adminApi, adminKeys } from '../api'
import ConfirmDialog from '../ConfirmDialog'
import { useAction } from '../useAction'
import NameDialog from './NameDialog'
import { applyOrder } from './order'
import SortableList from './SortableList'

const HOME_GROUP_ID = 1

// 樂觀更新排序：先改快取，失敗時還原
function useReorder<V>(save: (v: V) => Promise<void>, apply: (groups: Group[], v: V) => Group[]) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: save,
    onMutate: async (v) => {
      await qc.cancelQueries({ queryKey: adminKeys.groups })
      const previous = qc.getQueryData<Group[]>(adminKeys.groups)
      if (previous) qc.setQueryData(adminKeys.groups, apply(previous, v))
      return { previous }
    },
    onError: (e, _v, ctx) => {
      if (ctx?.previous) qc.setQueryData(adminKeys.groups, ctx.previous)
      toast.error(`排序未儲存：${e.message}`)
    },
    onSettled: () => qc.invalidateQueries({ queryKey: adminKeys.groups }),
  })
}

type Editing = { mode: 'create' } | { mode: 'rename'; group: Group }
type Deleting = { kind: 'group'; group: Group } | { kind: 'page'; page: PageRef }

export default function GroupsPage() {
  useDocumentTitle('頁籤與頁面｜後台')
  const query = useQuery({ queryKey: adminKeys.groups, queryFn: adminApi.groups })
  const [params, setParams] = useSearchParams()
  const [editing, setEditing] = useState<Editing | null>(null)
  const [deleting, setDeleting] = useState<Deleting | null>(null)
  const action = useAction(adminKeys.groups)

  const groupOrder = useReorder(adminApi.setGroupOrder, (groups, ids: number[]) => applyOrder(groups, ids, (g) => g.id))
  const pageOrder = useReorder(
    (v: { groupId: number; ids: number[] }) => adminApi.setPageOrder(v.groupId, v.ids),
    (groups, v) => groups.map((g) => (g.id === v.groupId ? { ...g, pages: applyOrder(g.pages, v.ids, (p) => p.id) } : g)),
  )

  if (query.isError) return <p role="alert" className="text-destructive">{query.error.message}</p>
  if (!query.data) return <Skeleton className="h-64" />

  const groups = query.data
  const selected = groups.find((g) => g.id === Number(params.get('group'))) ?? groups[0]
  const select = (id: number) => setParams({ group: String(id) }, { replace: true })

  const saveName = (name: string) => {
    const e = editing!
    action.mutate(() => (e.mode === 'rename' ? adminApi.renameGroup(e.group.id, name) : adminApi.createGroup(name)), {
      onSuccess: (result) => {
        setEditing(null)
        toast.success(e.mode === 'rename' ? '已更新頁籤' : '已新增頁籤')
        if (e.mode === 'create') select((result as Group).id)
      },
    })
  }

  const confirmDelete = () => {
    const d = deleting!
    setDeleting(null)
    action.mutate(() => (d.kind === 'group' ? adminApi.deleteGroup(d.group.id) : adminApi.deletePage(d.page.id)), {
      onSuccess: () => toast.success(d.kind === 'group' ? '已刪除頁籤' : '已刪除頁面'),
    })
  }

  return (
    <div className="space-y-6">
      <h1 className="font-serif text-2xl font-bold">頁籤與頁面</h1>
      <div className="grid gap-6 lg:grid-cols-[340px_1fr]">
        <section className="rounded-lg border bg-background p-4">
          <div className="mb-3 flex items-center justify-between">
            <h2 className="font-semibold">頁籤</h2>
            <Button size="sm" onClick={() => setEditing({ mode: 'create' })}>
              新增頁籤
            </Button>
          </div>
          <SortableList
            label="頁籤"
            items={groups}
            getId={(g) => g.id}
            onReorder={(ids) => groupOrder.mutate(ids)}
            renderItem={(g, handle) => (
              <div
                className={`flex items-center gap-1 rounded-md border px-2 py-1 ${g.id === selected?.id ? 'border-primary bg-primary/5' : 'bg-background'}`}
              >
                {handle}
                <button
                  type="button"
                  aria-current={g.id === selected?.id}
                  onClick={() => select(g.id)}
                  className="min-w-0 flex-1 truncate py-1 text-left"
                >
                  {g.name} <span className="text-xs text-muted-foreground">({g.pages.length})</span>
                </button>
                <Button size="sm" variant="ghost" aria-label={`重新命名 ${g.name}`} onClick={() => setEditing({ mode: 'rename', group: g })}>
                  改名
                </Button>
                <Button
                  size="sm"
                  variant="ghost"
                  aria-label={`刪除 ${g.name}`}
                  disabled={g.id === HOME_GROUP_ID}
                  title={g.id === HOME_GROUP_ID ? '首頁頁籤不可刪除' : undefined}
                  onClick={() => setDeleting({ kind: 'group', group: g })}
                >
                  刪除
                </Button>
              </div>
            )}
          />
        </section>

        <section className="rounded-lg border bg-background p-4">
          {selected ? (
            <>
              <div className="mb-3 flex items-center justify-between gap-2">
                <h2 className="font-semibold">{selected.name}的頁面</h2>
                <Button size="sm" asChild>
                  <Link to={`/admin/pages/new?group=${selected.id}`}>新增頁面</Link>
                </Button>
              </div>
              {selected.pages.length === 0 ? (
                <p className="text-sm text-muted-foreground">此頁籤還沒有頁面。</p>
              ) : (
                <SortableList
                  label="頁面"
                  items={selected.pages}
                  getId={(p) => p.id}
                  onReorder={(ids) => pageOrder.mutate({ groupId: selected.id, ids })}
                  renderItem={(p, handle) => (
                    <div className="flex items-center gap-1 rounded-md border bg-background px-2 py-1">
                      {handle}
                      <Link to={`/admin/pages/${p.id}`} className="min-w-0 flex-1 truncate py-1 hover:underline">
                        {p.name}
                      </Link>
                      <Button size="sm" variant="ghost" asChild>
                        <a href={`/page/${p.id}`} target="_blank" rel="noreferrer" aria-label={`預覽 ${p.name}`}>
                          預覽
                        </a>
                      </Button>
                      <Button size="sm" variant="ghost" aria-label={`刪除 ${p.name}`} onClick={() => setDeleting({ kind: 'page', page: p })}>
                        刪除
                      </Button>
                    </div>
                  )}
                />
              )}
            </>
          ) : (
            <p className="text-sm text-muted-foreground">請先新增頁籤。</p>
          )}
        </section>
      </div>

      <NameDialog
        key={editing?.mode === 'rename' ? editing.group.id : 'new'}
        open={editing !== null}
        title={editing?.mode === 'rename' ? '重新命名頁籤' : '新增頁籤'}
        initial={editing?.mode === 'rename' ? editing.group.name : ''}
        pending={action.isPending}
        onOpenChange={(open) => !open && setEditing(null)}
        onSubmit={saveName}
      />
      <ConfirmDialog
        open={deleting !== null}
        onOpenChange={(open) => !open && setDeleting(null)}
        title={deleting?.kind === 'group' ? '刪除頁籤' : '刪除頁面'}
        description={`確定要刪除「${deleting?.kind === 'group' ? deleting.group.name : deleting?.page.name}」嗎？此動作無法復原。`}
        confirmLabel="刪除"
        destructive
        onConfirm={confirmDelete}
      />
    </div>
  )
}
