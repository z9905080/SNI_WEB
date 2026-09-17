import type { ReactNode } from 'react'
import { useEditorState, type Editor } from '@tiptap/react'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

type Props = { editor: Editor; onInsertImage: () => void; onInsertVideo: () => void }

const aligns = [
  ['left', '靠左對齊', '左'],
  ['center', '置中對齊', '中'],
  ['right', '靠右對齊', '右'],
  ['justify', '左右對齊', '齊'],
] as const

export default function Toolbar({ editor, onInsertImage, onInsertVideo }: Props) {
  const s = useEditorState({
    editor,
    selector: ({ editor: e }) => ({
      heading: ([1, 2, 3, 4] as const).find((level) => e.isActive('heading', { level })) ?? 0,
      bold: e.isActive('bold'),
      italic: e.isActive('italic'),
      underline: e.isActive('underline'),
      strike: e.isActive('strike'),
      align: aligns.find(([a]) => e.isActive({ textAlign: a }))?.[0] ?? 'left',
      bullet: e.isActive('bulletList'),
      ordered: e.isActive('orderedList'),
      link: e.isActive('link'),
      table: e.isActive('table'),
      canUndo: e.can().undo(),
      canRedo: e.can().redo(),
    }),
  })
  const chain = () => editor.chain().focus()

  const editLink = () => {
    const current = (editor.getAttributes('link').href as string | undefined) ?? ''
    const url = window.prompt('連結網址（清空代表移除連結）', current)
    if (url === null) return
    if (url.trim() === '') chain().extendMarkRange('link').unsetLink().run()
    else chain().extendMarkRange('link').setLink({ href: url.trim() }).run()
  }

  return (
    <div role="toolbar" aria-label="編輯工具" className="flex flex-wrap items-center gap-1 border-b bg-muted/40 p-2">
      <select
        aria-label="段落格式"
        value={s.heading}
        onChange={(e) => {
          const level = Number(e.target.value) as 0 | 1 | 2 | 3 | 4
          if (level) chain().setHeading({ level }).run()
          else chain().setParagraph().run()
        }}
        className="h-8 rounded-md border bg-background px-2 text-sm"
      >
        <option value={0}>內文</option>
        <option value={1}>標題 1</option>
        <option value={2}>標題 2</option>
        <option value={3}>標題 3</option>
        <option value={4}>標題 4</option>
      </select>
      <Divider />
      <Tool label="粗體" active={s.bold} onClick={() => chain().toggleBold().run()}><b>B</b></Tool>
      <Tool label="斜體" active={s.italic} onClick={() => chain().toggleItalic().run()}><i>I</i></Tool>
      <Tool label="底線" active={s.underline} onClick={() => chain().toggleUnderline().run()}><u>U</u></Tool>
      <Tool label="刪除線" active={s.strike} onClick={() => chain().toggleStrike().run()}><s>S</s></Tool>
      <ColorTool label="文字顏色" onPick={(c) => chain().setColor(c).run()} onClear={() => chain().unsetColor().run()}>
        <span className="border-b-2 border-current">A</span>
      </ColorTool>
      <ColorTool label="螢光筆" onPick={(c) => chain().setHighlight({ color: c }).run()} onClear={() => chain().unsetHighlight().run()}>
        <span className="bg-yellow-200 px-0.5">A</span>
      </ColorTool>
      <Divider />
      {aligns.map(([value, label, text]) => (
        <Tool key={value} label={label} active={s.align === value} onClick={() => chain().setTextAlign(value).run()}>
          {text}
        </Tool>
      ))}
      <Divider />
      <Tool label="項目清單" active={s.bullet} onClick={() => chain().toggleBulletList().run()}>•</Tool>
      <Tool label="編號清單" active={s.ordered} onClick={() => chain().toggleOrderedList().run()}>1.</Tool>
      <Divider />
      <Tool label="連結" active={s.link} onClick={editLink}>連結</Tool>
      <Tool label="插入圖片" onClick={onInsertImage}>圖片</Tool>
      <Tool label="插入影片" onClick={onInsertVideo}>影片</Tool>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button type="button" size="sm" variant="ghost" className="h-8">表格</Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem onSelect={() => chain().insertTable({ rows: 3, cols: 3, withHeaderRow: false }).run()}>
            插入 3×3 表格
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem disabled={!s.table} onSelect={() => chain().addRowAfter().run()}>在下方插入列</DropdownMenuItem>
          <DropdownMenuItem disabled={!s.table} onSelect={() => chain().addColumnAfter().run()}>在右側插入欄</DropdownMenuItem>
          <DropdownMenuItem disabled={!s.table} onSelect={() => chain().deleteRow().run()}>刪除這一列</DropdownMenuItem>
          <DropdownMenuItem disabled={!s.table} onSelect={() => chain().deleteColumn().run()}>刪除這一欄</DropdownMenuItem>
          <DropdownMenuItem disabled={!s.table} onSelect={() => chain().deleteTable().run()}>刪除表格</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <Divider />
      <Tool label="復原" disabled={!s.canUndo} onClick={() => chain().undo().run()}>↶</Tool>
      <Tool label="重做" disabled={!s.canRedo} onClick={() => chain().redo().run()}>↷</Tool>
    </div>
  )
}

function Divider() {
  return <span aria-hidden="true" className="mx-1 h-5 w-px bg-border" />
}

type ToolProps = { label: string; active?: boolean; disabled?: boolean; onClick: () => void; children: ReactNode }

function Tool({ label, active, disabled, onClick, children }: ToolProps) {
  return (
    <Button
      type="button"
      size="sm"
      variant={active ? 'secondary' : 'ghost'}
      aria-label={label}
      aria-pressed={active}
      title={label}
      disabled={disabled}
      onClick={onClick}
      className="h-8 min-w-8 px-2"
    >
      {children}
    </Button>
  )
}

type ColorToolProps = { label: string; onPick: (color: string) => void; onClear: () => void; children: ReactNode }

function ColorTool({ label, onPick, onClear, children }: ColorToolProps) {
  return (
    <span className="inline-flex items-center">
      <label title={label} className="relative inline-flex h-8 min-w-8 cursor-pointer items-center justify-center rounded-md px-2 hover:bg-muted focus-within:ring-2 focus-within:ring-ring">
        {children}
        <input type="color" aria-label={label} className="absolute inset-0 cursor-pointer opacity-0" onChange={(e) => onPick(e.target.value)} />
      </label>
      <button type="button" aria-label={`清除${label}`} title={`清除${label}`} onClick={onClear} className="px-1 text-xs text-muted-foreground hover:text-foreground">
        ×
      </button>
    </span>
  )
}
