import { Editor, Extension, mergeAttributes, Node, type AnyExtension } from '@tiptap/core'
import Highlight from '@tiptap/extension-highlight'
import Image from '@tiptap/extension-image'
import { TableKit } from '@tiptap/extension-table'
import TextAlign from '@tiptap/extension-text-align'
import { TextStyleKit } from '@tiptap/extension-text-style'
import StarterKit from '@tiptap/starter-kit'
import { isAllowedIframe } from './iframeHosts'

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    iframe: { setIframe: (attrs: { src: string }) => ReturnType }
  }
}

export const Iframe = Node.create({
  name: 'iframe',
  group: 'block',
  atom: true,
  draggable: true,
  addAttributes() {
    return {
      src: { default: null },
      width: { default: '100%' },
      height: { default: null },
      frameborder: { default: null },
      allow: { default: null },
      allowfullscreen: { default: 'allowfullscreen' },
    }
  },
  parseHTML() {
    return [{ tag: 'iframe', getAttrs: (el) => (isAllowedIframe(el.getAttribute('src') ?? '') ? null : false) }]
  },
  renderHTML({ HTMLAttributes }) {
    return ['iframe', mergeAttributes(HTMLAttributes)]
  },
  addCommands() {
    return {
      setIframe:
        (attrs) =>
        ({ commands }) =>
          isAllowedIframe(attrs.src) && commands.insertContent({ type: this.name, attrs }),
    }
  },
})

// 以括號深度切開 CSS 宣告，並移除其他擴充已負責輸出的屬性
export function removeStyleProps(style: string | null, props: string[]): string | null {
  if (!style) return null
  const decls: string[] = []
  let depth = 0
  let current = ''
  for (const ch of style) {
    if (ch === '(') depth++
    if (ch === ')') depth--
    if (ch === ';' && depth === 0) {
      decls.push(current)
      current = ''
    } else {
      current += ch
    }
  }
  decls.push(current)
  const kept = decls
    .map((d) => d.trim())
    .filter((d) => d.includes(':') && !props.includes(d.slice(0, d.indexOf(':')).trim().toLowerCase()))
  return kept.length ? `${kept.join('; ')};` : null
}

const keep = (name: string) => ({
  default: null,
  parseHTML: (el: HTMLElement) => el.getAttribute(name),
  renderHTML: (attrs: Record<string, unknown>) => (attrs[name] ? { [name]: attrs[name] } : {}),
})

const keepStyle = (handledBy: string[]) => ({
  default: null,
  parseHTML: (el: HTMLElement) => removeStyleProps(el.getAttribute('style'), handledBy),
  renderHTML: (attrs: Record<string, unknown>) => (attrs.style ? { style: attrs.style } : {}),
})

// 舊內容大量使用行內樣式與表格屬性，這裡讓 Tiptap 原樣保留
const PreserveAttributes = Extension.create({
  name: 'preserveAttributes',
  addGlobalAttributes() {
    return [
      { types: ['paragraph', 'heading'], attributes: { style: keepStyle(['text-align']), class: keep('class') } },
      {
        types: ['textStyle'],
        attributes: {
          style: keepStyle(['color', 'background-color', 'font-family', 'font-size', 'line-height']),
          class: keep('class'),
        },
      },
      {
        types: ['table'],
        attributes: {
          style: keepStyle([]),
          class: keep('class'),
          border: keep('border'),
          cellpadding: keep('cellpadding'),
          cellspacing: keep('cellspacing'),
          width: keep('width'),
          align: keep('align'),
        },
      },
      {
        types: ['tableRow', 'tableCell', 'tableHeader'],
        attributes: {
          style: keepStyle([]),
          class: keep('class'),
          width: keep('width'),
          height: keep('height'),
          align: keep('align'),
          valign: keep('valign'),
          bgcolor: keep('bgcolor'),
        },
      },
      {
        types: ['image', 'blockquote', 'bulletList', 'orderedList', 'listItem', 'iframe'],
        attributes: { style: keepStyle([]), class: keep('class') },
      },
      { types: ['link'], attributes: { title: keep('title') } },
    ]
  },
})

export function createExtensions(): AnyExtension[] {
  return [
    StarterKit.configure({
      heading: { levels: [1, 2, 3, 4] },
      link: { openOnClick: false, HTMLAttributes: { target: null, rel: null } },
    }),
    TextStyleKit,
    Highlight.configure({ multicolor: true }),
    TextAlign.configure({ types: ['heading', 'paragraph'] }),
    Image.configure({ inline: true }),
    TableKit.configure({ table: { resizable: false } }),
    Iframe,
    PreserveAttributes,
  ]
}

export function roundTrip(html: string): string {
  const editor = new Editor({ extensions: createExtensions(), content: html })
  try {
    return editor.getHTML()
  } finally {
    editor.destroy()
  }
}
