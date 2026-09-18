import { useEffect, useRef, useState } from 'react'
import { html as htmlLanguage } from '@codemirror/lang-html'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { EditorContent, useEditor } from '@tiptap/react'
import CodeMirror from '@uiw/react-codemirror'
import { Link, useBlocker, useNavigate, useParams, useSearchParams } from 'react-router'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import type { Group, Page } from '@/shared/types'
import { useDocumentTitle } from '@/site/useDocumentTitle'
import { adminApi, adminKeys } from '../api'
import ConfirmDialog from '../ConfirmDialog'
import ImagePickerDialog from '../images/ImagePickerDialog'
import { createExtensions, roundTrip } from './extensions'
import { detectLoss } from './lossDetection'
import Toolbar from './Toolbar'
import VideoDialog from './VideoDialog'

export default function PageEditorPage() {
  const { id } = useParams()
  const pageId = id ? Number(id) : undefined
  const [params] = useSearchParams()
  const groups = useQuery({ queryKey: adminKeys.groups, queryFn: adminApi.groups })
  const page = useQuery({
    queryKey: adminKeys.page(pageId ?? 0),
    queryFn: () => adminApi.page(pageId!),
    enabled: pageId !== undefined,
  })

  const error = groups.error ?? page.error
  if (error) return <p role="alert" className="text-destructive">{error.message}</p>
  if (!groups.data || (pageId !== undefined && !page.data)) return <Skeleton className="h-96" />

  const initial: Page = page.data ?? {
    id: 0,
    name: '',
    group_id: Number(params.get('group')) || groups.data[0]?.id || 0,
    html: '',
  }
  return <PageForm key={pageId ?? 'new'} initial={initial} groups={groups.data} />
}

function PageForm({ initial, groups }: { initial: Page; groups: Group[] }) {
  const isNew = initial.id === 0
  useDocumentTitle(`${isNew ? '新增頁面' : `編輯：${initial.name}`}｜後台`)
  const qc = useQueryClient()
  const navigate = useNavigate()

  const [name, setName] = useState(initial.name)
  const [groupId, setGroupId] = useState(initial.group_id)
  // source 永遠是要送出的內容：視覺模式的每次修改都會同步回來
  const [source, setSource] = useState(initial.html)
  const [lost, setLost] = useState(() => detectLoss(initial.html, roundTrip(initial.html)))
  const [mode, setMode] = useState<'visual' | 'source'>(lost.length ? 'source' : 'visual')
  const [switchWarning, setSwitchWarning] = useState<string[] | null>(null)
  const [picker, setPicker] = useState(false)
  const [video, setVideo] = useState(false)

  // 導覽攔截在 render 之外判斷，需讀取即時值
  const dirtyRef = useRef(false)
  const [dirty, setDirtyState] = useState(false)
  const setDirty = (value: boolean) => {
    dirtyRef.current = value
    setDirtyState(value)
  }

  const editor = useEditor({
    extensions: createExtensions(),
    content: initial.html,
    onUpdate: ({ editor: e }) => {
      setSource(e.getHTML())
      setDirty(true)
    },
    editorProps: {
      attributes: { class: 'legacy-content min-h-[50vh] p-4 focus:outline-none', 'aria-label': '頁面內容' },
    },
  })

  const blocker = useBlocker(
    ({ currentLocation, nextLocation }) => dirtyRef.current && currentLocation.pathname !== nextLocation.pathname,
  )

  useEffect(() => {
    if (!dirty) return
    const warn = (e: BeforeUnloadEvent) => e.preventDefault()
    window.addEventListener('beforeunload', warn)
    return () => window.removeEventListener('beforeunload', warn)
  }, [dirty])

  const save = useMutation({
    mutationFn: () => {
      const body = { name: name.trim(), group_id: groupId, html: source }
      return isNew ? adminApi.createPage(body) : adminApi.updatePage(initial.id, body)
    },
    onSuccess: (saved) => {
      setDirty(false)
      qc.setQueryData(adminKeys.page(saved.id), saved)
      qc.invalidateQueries({ queryKey: adminKeys.groups })
      toast.success('已儲存')
      if (isNew) {
        navigate(`/admin/pages/${saved.id}`, { replace: true })
        return
      }
      // 伺服器會過濾 HTML，以回傳內容為準
      setSource(saved.html)
      if (mode === 'visual') editor?.commands.setContent(saved.html, { emitUpdate: false })
    },
    onError: (e) => toast.error(e.message),
  })

  if (!editor) return null

  const toVisual = (force: boolean) => {
    const loss = detectLoss(source, roundTrip(source))
    if (loss.length && !force) {
      setSwitchWarning(loss)
      return
    }
    editor.commands.setContent(source, { emitUpdate: false })
    setSwitchWarning(null)
    setLost([])
    setMode('visual')
  }

  const markDirty = () => setDirty(true)

  return (
    <form
      className="space-y-4"
      onSubmit={(e) => {
        e.preventDefault()
        save.mutate()
      }}
    >
      <div className="flex flex-wrap items-center justify-between gap-2">
        <h1 className="font-serif text-2xl font-bold">{isNew ? '新增頁面' : '編輯頁面'}</h1>
        <div className="flex flex-wrap gap-2">
          <Button type="button" variant="outline" asChild>
            <Link to={`/admin/groups?group=${groupId}`}>返回列表</Link>
          </Button>
          {!isNew && (
            <Button type="button" variant="outline" asChild>
              <a href={`/page/${initial.id}`} target="_blank" rel="noreferrer">
                預覽
              </a>
            </Button>
          )}
          <Button type="submit" disabled={save.isPending}>
            {save.isPending ? '儲存中…' : '儲存'}
          </Button>
        </div>
      </div>

      <div className="grid gap-4 sm:grid-cols-[1fr_240px]">
        <div className="space-y-2">
          <Label htmlFor="page-name">頁面名稱</Label>
          <Input
            id="page-name"
            value={name}
            maxLength={50}
            required
            onChange={(e) => {
              setName(e.target.value)
              markDirty()
            }}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="page-group">所屬頁籤</Label>
          <select
            id="page-group"
            value={groupId}
            onChange={(e) => {
              setGroupId(Number(e.target.value))
              markDirty()
            }}
            className="h-9 w-full rounded-md border bg-background px-3 text-sm"
          >
            {groups.map((g) => (
              <option key={g.id} value={g.id}>
                {g.name}
              </option>
            ))}
          </select>
        </div>
      </div>

      {mode === 'source' && lost.length > 0 && (
        <div role="status" className="rounded-md border border-amber-300 bg-amber-50 p-3 text-sm text-amber-900">
          這個頁面有視覺編輯器無法完整保留的格式（{lost.slice(0, 5).join('、')}
          {lost.length > 5 ? ' 等' : ''}），因此以 HTML 原始碼模式開啟。切換到視覺編輯後再儲存，這些格式會遺失。
        </div>
      )}

      <div className="overflow-hidden rounded-md border bg-background">
        <div className="flex items-center justify-between border-b px-2 py-1">
          <span className="text-sm text-muted-foreground">內容</span>
          <div role="group" aria-label="編輯模式" className="flex gap-1">
            <Button
              type="button"
              size="sm"
              variant={mode === 'visual' ? 'secondary' : 'ghost'}
              aria-pressed={mode === 'visual'}
              onClick={() => mode === 'source' && toVisual(false)}
            >
              視覺編輯
            </Button>
            <Button
              type="button"
              size="sm"
              variant={mode === 'source' ? 'secondary' : 'ghost'}
              aria-pressed={mode === 'source'}
              onClick={() => setMode('source')}
            >
              HTML 原始碼
            </Button>
          </div>
        </div>
        {/* 編輯器保持掛載，切換模式時只隱藏，避免重建 */}
        <div hidden={mode !== 'visual'}>
          <Toolbar editor={editor} onInsertImage={() => setPicker(true)} onInsertVideo={() => setVideo(true)} />
          <EditorContent editor={editor} />
        </div>
        {mode === 'source' && (
          <CodeMirror
            value={source}
            height="60vh"
            extensions={[htmlLanguage()]}
            basicSetup={{ foldGutter: false }}
            onChange={(value) => {
              setSource(value)
              markDirty()
            }}
          />
        )}
      </div>

      <ImagePickerDialog
        open={picker}
        onOpenChange={setPicker}
        onSelect={(url) => {
          editor.chain().focus().setImage({ src: url }).run()
          setPicker(false)
        }}
      />
      <VideoDialog
        open={video}
        onOpenChange={setVideo}
        onInsert={(src) => {
          editor.chain().focus().setIframe({ src }).run()
          setVideo(false)
        }}
      />
      <ConfirmDialog
        open={switchWarning !== null}
        onOpenChange={(open) => !open && setSwitchWarning(null)}
        title="切換到視覺編輯？"
        description={`下列格式無法在視覺編輯器中保留：${switchWarning?.slice(0, 5).join('、')}。切換後若修改並儲存，這些格式會遺失。`}
        confirmLabel="切換到視覺編輯"
        onConfirm={() => toVisual(true)}
      />
      <ConfirmDialog
        open={blocker.state === 'blocked'}
        onOpenChange={(open) => !open && blocker.reset?.()}
        title="還有變更沒有儲存"
        description="離開這個頁面會捨棄尚未儲存的變更。"
        confirmLabel="捨棄變更並離開"
        destructive
        onConfirm={() => blocker.proceed?.()}
      />
    </form>
  )
}
