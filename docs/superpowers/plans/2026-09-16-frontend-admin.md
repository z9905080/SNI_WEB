# 生長之家網站重寫 — 後台（React）Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 完成 `/admin/*` 後台：登入、頁籤與頁面管理（拖曳排序）、Tiptap 頁面編輯器（含 HTML 原始碼模式與遺失偵測）、圖庫、輪播圖、跑馬燈、網站設定。

**Architecture:** 後台是 `web/src/admin/`，由 `App.tsx` 以 `React.lazy` 載入，前台 bundle 不含後台與 Tiptap。資料存取用 TanStack Query（query key 一律以 `['admin', ...]` 開頭），任何 401 由全域監聽導向登入頁。改用 React Router data router（`createBrowserRouter`）以便編輯頁使用 `useBlocker` 提示未存檔。

**Tech Stack:** React 19、React Router 7、TanStack Query 5、shadcn/ui（Radix）、sonner、dnd-kit、Tiptap 3、CodeMirror 6（`@uiw/react-codemirror`）、Vitest、Playwright。

**Spec:** `docs/superpowers/specs/2026-09-16-sni-web-rewrite-design.md`

**前置條件：** `2026-09-16-backend.md` 與 `2026-09-16-frontend-site.md` 已完成。

## Global Constraints

- 後台路由：`/admin/login`；其餘頁面需登入。收到 401 時導向 `/admin/login?redirect=<原路徑>`；`redirect` 只接受 `/admin/` 開頭的站內路徑。
- 左側選單：頁籤與頁面、輪播圖、跑馬燈、圖庫、網站設定、登出。
- 所有操作結果以 toast（sonner）顯示；錯誤直接顯示 API 回傳的 `message`。
- 拖曳排序放開即儲存，樂觀更新，失敗時回復並提示。
- 頁面編輯：欄位為名稱、所屬群組、內容；原始碼模式存檔時直接送出原始碼；載入時若 Tiptap 會遺失標籤、屬性或文字，預設進入原始碼模式並顯示警告。
- iframe 網域白名單與後端 `content.IframeHosts` 一致：`www.youtube.com`、`youtube.com`、`www.youtube-nocookie.com`、`players.brightcove.net`、`drive.google.com`、`www.facebook.com`。
- 圖片上傳：jpg/png/gif/webp、單檔 5MB（前端先檢查，後端再檢查）；每個檔案一個請求以顯示個別進度。
- 前台 bundle 不可包含 Tiptap / ProseMirror / dnd-kit（E2E 檢查）。
- 所有程式先寫失敗的測試再實作（純版面除外）；拖曳與編輯器整體流程由 E2E 覆蓋。
- Commit 訊息沿用 repo 慣例：`[feat] ...`、`[test] ...`、`[chore] ...`。

## 設計方向

沿用前台的設計方向（`2026-09-16-frontend-site.md`），後台以工具性為主、保持安靜：
- shadcn/ui neutral 主題，`--primary` 換成寶藍 `#2230B8`；頁面 H1 使用明體（`font-serif`），與前台一致。
- 不使用卡片陰影堆疊、全大寫標籤、以「·」串接的中繼資訊；檔名與大小等資訊分行呈現。
- 按鈕文字與結果訊息使用同一個動詞：「儲存」→「已儲存」、「新增頁籤」→「已新增頁籤」、「刪除」→「已刪除圖片」。
- 錯誤訊息說明發生什麼事與怎麼處理，直接使用 API 的 `message`。

## 本 plan 的實作層決定

- 群組下拉選單用原生 `<select>`（不用 Radix Select，測試與無障礙都較單純）。
- 連結網址輸入用 `window.prompt`（後台內部工具，不另做對話框）。
- 共用確認對話框 `ConfirmDialog`（shadcn AlertDialog）。
- 預設 query 重試規則移到 `shared/api.ts` 的 `retryUnlessClientError`，並設為全域預設（取代前台 `usePage` 內的個別設定）。

## 檔案結構

```
web/src/
  main.tsx                         改用 createBrowserRouter；全域 retry 規則
  shared/api.ts                    新增 retryUnlessClientError
  site/queries.ts                  移除個別 retry 設定
  components/ui/*                  shadcn 產生（勿手改，主題色除外）
  lib/utils.ts                     shadcn 產生
  admin/
    AdminApp.tsx                   路由、Toaster、401 監聽
    api.ts(+test)                  後台 API 與 uploadImage（XHR）
    auth.tsx(+test)                useMe、RequireAuth、useRedirectOn401、safeRedirect
    LoginPage.tsx
    AdminLayout.tsx                側欄與登出
    ConfirmDialog.tsx
    groups/
      order.ts(+test)              applyOrder
      SortableList.tsx             dnd-kit 通用排序清單
      GroupsPage.tsx(+test)
      NameDialog.tsx
    images/
      useUploader.ts(+test)
      ImagesPage.tsx
      ImagePickerDialog.tsx
    editor/
      iframeHosts.ts(+test)        白名單、toEmbedUrl
      extensions.ts(+test)         Tiptap 擴充設定（含 iframe、style/class 保留）
      lossDetection.ts(+test)
      Toolbar.tsx
      VideoDialog.tsx
      PageEditorPage.tsx
    carousels/CarouselsPage.tsx
    marquees/MarqueesPage.tsx(+test)
    settings/SettingsPage.tsx(+test)
web/e2e/admin.spec.ts
```

---

### Task 1: 後台基礎（shadcn/ui、data router、API、登入與版面）

**Files:**
- Modify: `web/src/main.tsx`、`web/src/shared/api.ts`、`web/src/site/queries.ts`、`web/src/index.css`（shadcn 產生後調整主色）
- Create: `web/components.json`、`web/src/components/ui/*`、`web/src/lib/utils.ts`（shadcn 產生）
- Create: `web/src/admin/api.ts`、`auth.tsx`、`LoginPage.tsx`、`AdminLayout.tsx`、`ConfirmDialog.tsx`
- Modify: `web/src/admin/AdminApp.tsx`（整個取代）
- Test: `web/src/admin/api.test.ts`、`web/src/admin/auth.test.tsx`

**Interfaces:**
- Consumes: `api`、`ApiError`、型別（前台 plan Task 2）
- Produces:
  ```ts
  // shared/api.ts
  export function retryUnlessClientError(failureCount: number, error: unknown): boolean
  // admin/api.ts
  export const adminKeys = { all: ['admin'], me: [...], groups: [...], page: (id: number) => [...], carousels, marquees, settings, images }
  export const adminApi = {
    me(): Promise<User>; login(account: string, password: string): Promise<User>; logout(): Promise<void>
    groups(): Promise<Group[]>; createGroup(name: string): Promise<Group>; renameGroup(id: number, name: string): Promise<void>
    deleteGroup(id: number): Promise<void>; setGroupOrder(ids: number[]): Promise<void>; setPageOrder(groupId: number, ids: number[]): Promise<void>
    page(id: number): Promise<Page>; createPage(p: PageInput): Promise<Page>; updatePage(id: number, p: Partial<PageInput>): Promise<Page>; deletePage(id: number): Promise<void>
    carousels(): Promise<Carousel[]>; createCarousel(c: CarouselInput): Promise<Carousel>; updateCarousel(id: number, c: CarouselInput): Promise<Carousel>; deleteCarousel(id: number): Promise<void>
    marquees(): Promise<Marquee[]>; createMarquee(m: MarqueeInput): Promise<Marquee>; updateMarquee(id: number, m: MarqueeInput): Promise<Marquee>; deleteMarquee(id: number): Promise<void>
    settings(): Promise<Settings>; saveSettings(s: Settings): Promise<Settings>
    images(): Promise<ImageItem[]>; imageUsages(name: string): Promise<ImageUsages>; deleteImage(name: string): Promise<void>
  }
  export type PageInput = { name: string; group_id: number; html: string }
  export type CarouselInput = { image: string; url: string }
  export type MarqueeInput = { text: string; color: string }
  export function uploadImage(file: File, onProgress: (ratio: number) => void): Promise<ImageItem>
  // admin/auth.tsx
  export function useMe(): UseQueryResult<User>
  export function RequireAuth(): JSX.Element     // 登入後渲染 <AdminLayout />
  export function useRedirectOn401(): void
  export function safeRedirect(value: string | null): string
  // admin/ConfirmDialog.tsx
  export default function ConfirmDialog(props: { open: boolean; onOpenChange(open: boolean): void; title: string; description: ReactNode; confirmLabel: string; destructive?: boolean; onConfirm(): void })
  ```

- [ ] **Step 1: 安裝 shadcn/ui 與套件**

```bash
cd web
pnpm dlx shadcn@latest init --base-color neutral --yes
pnpm dlx shadcn@latest add button input label textarea dialog alert-dialog dropdown-menu sonner card skeleton --yes
pnpm add @dnd-kit/core @dnd-kit/sortable @dnd-kit/utilities @dnd-kit/modifiers \
  @tiptap/react @tiptap/pm @tiptap/core @tiptap/starter-kit @tiptap/extension-text-style \
  @tiptap/extension-highlight @tiptap/extension-text-align @tiptap/extension-image @tiptap/extension-table \
  @uiw/react-codemirror @codemirror/lang-html
cd ..
```

> `init` 若仍詢問問題，選 Vite、`src/index.css`、別名 `@/components`、`@/lib/utils`。完成後確認 `web/src/index.css` 仍保留原本的 `@theme`、`body`、`.legacy-content` 等規則。

在 `web/src/index.css` 的 `:root { ... }`（shadcn 產生）中，把 `--primary` 與 `--primary-foreground` 改為：

```css
  --primary: #2230b8;
  --primary-foreground: #ffffff;
  --ring: #2230b8;
```

Radix 與 sonner 會用到 jsdom 沒有的 API，在 `web/src/test/setup.ts` 最後加入：

```ts
// jsdom 沒有的瀏覽器 API（sonner 讀取 prefers-color-scheme、Radix Popper 量測尺寸）
if (!window.matchMedia) {
  window.matchMedia = (query: string) =>
    ({
      matches: false,
      media: query,
      onchange: null,
      addEventListener: () => {},
      removeEventListener: () => {},
      addListener: () => {},
      removeListener: () => {},
      dispatchEvent: () => false,
    }) as MediaQueryList
}
window.ResizeObserver ??= class {
  observe() {}
  unobserve() {}
  disconnect() {}
}
Element.prototype.scrollIntoView ??= () => {}
```

Run: `pnpm --dir web test && pnpm --dir web build`
Expected: 既有測試 PASS；建置成功。

- [ ] **Step 2: 寫失敗的測試**

`web/src/admin/api.test.ts`：

```ts
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '@/shared/api'
import { adminApi, uploadImage } from './api'

afterEach(() => vi.unstubAllGlobals())

describe('adminApi', () => {
  it('解開回應外層並送出正確的請求', async () => {
    const fetchMock = vi.fn(async () =>
      new Response(JSON.stringify({ page: { id: 5, group_id: 2, name: 'x', html: '' } }), { status: 201 }),
    )
    vi.stubGlobal('fetch', fetchMock)
    const page = await adminApi.createPage({ name: 'x', group_id: 2, html: '' })
    expect(page.id).toBe(5)
    expect(fetchMock).toHaveBeenCalledWith('/api/v1/admin/pages', expect.objectContaining({ method: 'POST' }))
  })

  it('圖片名稱會編碼', async () => {
    const fetchMock = vi.fn(async () => new Response(null, { status: 204 }))
    vi.stubGlobal('fetch', fetchMock)
    await adminApi.deleteImage('a b.jpg')
    expect(fetchMock.mock.calls[0][0]).toBe('/api/v1/admin/images/a%20b.jpg')
  })
})

class FakeXHR {
  static last: FakeXHR
  upload = { onprogress: null as ((e: { lengthComputable: boolean; loaded: number; total: number }) => void) | null }
  onload: (() => void) | null = null
  onerror: (() => void) | null = null
  status = 0
  responseText = ''
  method = ''
  url = ''
  body: FormData | null = null
  constructor() {
    FakeXHR.last = this
  }
  open(method: string, url: string) {
    this.method = method
    this.url = url
  }
  send(body: FormData) {
    this.body = body
  }
}

describe('uploadImage', () => {
  it('回報進度並回傳上傳結果', async () => {
    vi.stubGlobal('XMLHttpRequest', FakeXHR)
    const progress: number[] = []
    const file = new File(['x'], 'a.png', { type: 'image/png' })
    const promise = uploadImage(file, (r) => progress.push(r))
    const xhr = FakeXHR.last
    expect(xhr.method).toBe('POST')
    expect(xhr.url).toBe('/api/v1/admin/images')
    expect(xhr.body?.get('files')).toBe(file)
    xhr.upload.onprogress?.({ lengthComputable: true, loaded: 50, total: 100 })
    xhr.status = 201
    xhr.responseText = JSON.stringify({ images: [{ name: 'n.png', url: '/php/picture/n.png', size: 1, mod_time: '' }] })
    xhr.onload?.()
    await expect(promise).resolves.toMatchObject({ name: 'n.png' })
    expect(progress).toEqual([0.5, 1])
  })

  it('錯誤回應轉成 ApiError', async () => {
    vi.stubGlobal('XMLHttpRequest', FakeXHR)
    const promise = uploadImage(new File(['x'], 'a.svg'), () => {})
    const xhr = FakeXHR.last
    xhr.status = 400
    xhr.responseText = JSON.stringify({ error: { code: 'unsupported_type', message: '不是支援的圖片格式' } })
    xhr.onload?.()
    await expect(promise).rejects.toMatchObject({ status: 400, code: 'unsupported_type', message: '不是支援的圖片格式' })
    await expect(promise).rejects.toBeInstanceOf(ApiError)
  })

  it('網路錯誤', async () => {
    vi.stubGlobal('XMLHttpRequest', FakeXHR)
    const promise = uploadImage(new File(['x'], 'a.png'), () => {})
    FakeXHR.last.onerror?.()
    await expect(promise).rejects.toMatchObject({ status: 0, code: 'network' })
  })
})
```

`web/src/admin/auth.test.tsx`：

```tsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, RouterProvider, useLocation } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { retryUnlessClientError } from '@/shared/api'
import AdminApp from './AdminApp'
import { safeRedirect } from './auth'

describe('safeRedirect', () => {
  it.each([
    [null, '/admin/groups'],
    ['/admin/pages/3?x=1', '/admin/pages/3?x=1'],
    ['/admin', '/admin/groups'],
    ['//evil.example/admin/', '/admin/groups'],
    ['https://evil.example/admin/x', '/admin/groups'],
    ['/page/3', '/admin/groups'],
  ])('%s → %s', (input, want) => {
    expect(safeRedirect(input)).toBe(want)
  })
})

function Where() {
  const l = useLocation()
  return <p data-testid="where">{l.pathname + l.search}</p>
}

function renderAdmin(path: string) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: retryUnlessClientError } } })
  const router = createMemoryRouter(
    [{ path: '/admin/*', element: <><AdminApp /><Where /></> }],
    { initialEntries: [path] },
  )
  render(
    <QueryClientProvider client={qc}>
      <RouterProvider router={router} />
    </QueryClientProvider>,
  )
}

const json = (status: number, body: unknown) => new Response(JSON.stringify(body), { status })

afterEach(() => vi.unstubAllGlobals())

describe('登入流程', () => {
  it('未登入時導向登入頁並帶 redirect，登入後回到原頁', async () => {
    let loggedIn = false
    vi.stubGlobal('fetch', vi.fn(async (url: string, init?: RequestInit) => {
      if (url.endsWith('/auth/me')) {
        return loggedIn
          ? json(200, { user: { id: 1, account: 'admin', name: '管理者' } })
          : json(401, { error: { code: 'unauthorized', message: '請先登入' } })
      }
      if (url.endsWith('/auth/login')) {
        const body = JSON.parse(String(init?.body))
        if (body.password !== 'password1') return json(401, { error: { code: 'invalid_credentials', message: '帳號或密碼錯誤' } })
        loggedIn = true
        return json(200, { user: { id: 1, account: 'admin', name: '管理者' } })
      }
      if (url.endsWith('/admin/settings')) {
        return json(200, { settings: { web_title: '生長之家', web_sub_title: '', facebook_url: '' } })
      }
      return json(404, { error: { code: 'not_found', message: '找不到' } })
    }))

    renderAdmin('/admin/settings')
    expect(await screen.findByTestId('where')).toHaveTextContent('/admin/login?redirect=%2Fadmin%2Fsettings')

    const user = userEvent.setup()
    await user.type(screen.getByLabelText('帳號'), 'admin')
    await user.type(screen.getByLabelText('密碼'), 'wrong')
    await user.click(screen.getByRole('button', { name: '登入' }))
    expect(await screen.findByRole('alert')).toHaveTextContent('帳號或密碼錯誤')

    await user.clear(screen.getByLabelText('密碼'))
    await user.type(screen.getByLabelText('密碼'), 'password1')
    await user.click(screen.getByRole('button', { name: '登入' }))
    // 手機與桌機版面各顯示一次使用者名稱（jsdom 不套用 CSS）
    expect(await screen.findAllByText('管理者')).not.toHaveLength(0)
    expect(screen.getByTestId('where')).toHaveTextContent('/admin/settings')
  })
})
```

> 此測試在 Task 6 建立 `SettingsPage` 前，`/admin/settings` 會由暫時的佔位頁面顯示；只要 `AdminLayout` 顯示使用者名稱即通過。

Run: `pnpm --dir web test`
Expected: FAIL（找不到模組）

- [ ] **Step 3: 全域 retry 規則與 data router**

`web/src/shared/api.ts` 最後加入：

```ts
// 4xx 不重試（例如 401、404），其他錯誤重試一次
export function retryUnlessClientError(failureCount: number, error: unknown) {
  if (error instanceof ApiError && error.status >= 400 && error.status < 500) return false
  return failureCount < 1
}
```

`web/src/site/queries.ts` 的 `usePage` 移除 `retry` 設定與 `ApiError` import：

```ts
export const usePage = (id: string) =>
  useQuery({ queryKey: ['page', id], queryFn: ({ signal }) => api<PageData>(`/pages/${id}`, { signal }) })
```

`web/src/main.tsx`（整個取代）：

```tsx
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { retryUnlessClientError } from '@/shared/api'
import { legacyHashTarget } from '@/site/legacyLinks'
import App from './App'
import './index.css'

// 舊網址 /#/5 → /page/5；在 Router 啟動前改寫，避免多一次導覽
const legacyTarget = legacyHashTarget(window.location.hash)
if (legacyTarget) window.history.replaceState(null, '', legacyTarget)

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 60_000, retry: retryUnlessClientError } },
})

// data router 才能使用 useBlocker（編輯頁未存檔提示）
const router = createBrowserRouter([{ path: '*', element: <App /> }])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
)
```

- [ ] **Step 4: 後台 API**

`web/src/admin/api.ts`：

```ts
import { api, ApiError } from '@/shared/api'
import type { Carousel, Group, ImageItem, ImageUsages, Marquee, Page, Settings, User } from '@/shared/types'

export type PageInput = { name: string; group_id: number; html: string }
export type CarouselInput = { image: string; url: string }
export type MarqueeInput = { text: string; color: string }

export const adminKeys = {
  all: ['admin'] as const,
  me: ['admin', 'me'] as const,
  groups: ['admin', 'groups'] as const,
  page: (id: number) => ['admin', 'page', id] as const,
  carousels: ['admin', 'carousels'] as const,
  marquees: ['admin', 'marquees'] as const,
  settings: ['admin', 'settings'] as const,
  images: ['admin', 'images'] as const,
}

const A = '/admin'
const post = (body: unknown) => ({ method: 'POST', body })
const patch = (body: unknown) => ({ method: 'PATCH', body })
const put = (body: unknown) => ({ method: 'PUT', body })
const del = { method: 'DELETE' }

export const adminApi = {
  me: () => api<{ user: User }>(`${A}/auth/me`).then((r) => r.user),
  login: (account: string, password: string) =>
    api<{ user: User }>(`${A}/auth/login`, post({ account, password })).then((r) => r.user),
  logout: () => api<void>(`${A}/auth/logout`, { method: 'POST' }),

  groups: () => api<{ groups: Group[] }>(`${A}/groups`).then((r) => r.groups),
  createGroup: (name: string) => api<{ group: Group }>(`${A}/groups`, post({ name })).then((r) => r.group),
  renameGroup: (id: number, name: string) => api<unknown>(`${A}/groups/${id}`, patch({ name })).then(() => {}),
  deleteGroup: (id: number) => api<void>(`${A}/groups/${id}`, del),
  setGroupOrder: (ids: number[]) => api<void>(`${A}/groups/order`, put({ ids })),
  setPageOrder: (groupId: number, ids: number[]) => api<void>(`${A}/groups/${groupId}/pages/order`, put({ ids })),

  page: (id: number) => api<{ page: Page }>(`${A}/pages/${id}`).then((r) => r.page),
  createPage: (p: PageInput) => api<{ page: Page }>(`${A}/pages`, post(p)).then((r) => r.page),
  updatePage: (id: number, p: Partial<PageInput>) => api<{ page: Page }>(`${A}/pages/${id}`, patch(p)).then((r) => r.page),
  deletePage: (id: number) => api<void>(`${A}/pages/${id}`, del),

  carousels: () => api<{ carousels: Carousel[] }>(`${A}/carousels`).then((r) => r.carousels),
  createCarousel: (c: CarouselInput) => api<{ carousel: Carousel }>(`${A}/carousels`, post(c)).then((r) => r.carousel),
  updateCarousel: (id: number, c: CarouselInput) =>
    api<{ carousel: Carousel }>(`${A}/carousels/${id}`, patch(c)).then((r) => r.carousel),
  deleteCarousel: (id: number) => api<void>(`${A}/carousels/${id}`, del),

  marquees: () => api<{ marquees: Marquee[] }>(`${A}/marquees`).then((r) => r.marquees),
  createMarquee: (m: MarqueeInput) => api<{ marquee: Marquee }>(`${A}/marquees`, post(m)).then((r) => r.marquee),
  updateMarquee: (id: number, m: MarqueeInput) =>
    api<{ marquee: Marquee }>(`${A}/marquees/${id}`, patch(m)).then((r) => r.marquee),
  deleteMarquee: (id: number) => api<void>(`${A}/marquees/${id}`, del),

  settings: () => api<{ settings: Settings }>(`${A}/settings`).then((r) => r.settings),
  saveSettings: (s: Settings) => api<{ settings: Settings }>(`${A}/settings`, put(s)).then((r) => r.settings),

  images: () => api<{ images: ImageItem[] }>(`${A}/images`).then((r) => r.images),
  imageUsages: (name: string) =>
    api<{ usages: ImageUsages }>(`${A}/images/${encodeURIComponent(name)}/usages`).then((r) => r.usages),
  deleteImage: (name: string) => api<void>(`${A}/images/${encodeURIComponent(name)}`, del),
}

// fetch 無法回報上傳進度，因此上傳改用 XMLHttpRequest
export function uploadImage(file: File, onProgress: (ratio: number) => void): Promise<ImageItem> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    const body = new FormData()
    body.append('files', file)
    xhr.open('POST', `/api/v1${A}/images`)
    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable) onProgress(e.loaded / e.total)
    }
    xhr.onload = () => {
      let data: { images?: ImageItem[]; error?: { code: string; message: string; request_id?: string } } | null = null
      try {
        data = JSON.parse(xhr.responseText)
      } catch {
        // 非 JSON 回應，下面統一處理
      }
      if (xhr.status >= 200 && xhr.status < 300 && data?.images?.[0]) {
        onProgress(1)
        resolve(data.images[0])
        return
      }
      const err = data?.error
      reject(new ApiError(xhr.status, err?.code ?? 'unknown', err?.message ?? `上傳失敗（${xhr.status}）`, err?.request_id))
    }
    xhr.onerror = () => reject(new ApiError(0, 'network', '無法連線到伺服器，請檢查網路後再試'))
    xhr.send(body)
  })
}
```

- [ ] **Step 5: 認證、登入頁、版面、確認對話框**

`web/src/admin/auth.tsx`：

```tsx
import { useEffect, useRef } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { useLocation, useNavigate } from 'react-router'
import { ApiError } from '@/shared/api'
import AdminLayout from './AdminLayout'
import { adminApi, adminKeys } from './api'

export const useMe = () => useQuery({ queryKey: adminKeys.me, queryFn: adminApi.me, staleTime: Infinity })

export function safeRedirect(value: string | null) {
  return value && value.startsWith('/admin/') ? value : '/admin/groups'
}

// 任何後台 query / mutation 收到 401 時，清掉後台快取並導向登入頁
export function useRedirectOn401() {
  const qc = useQueryClient()
  const navigate = useNavigate()
  const location = useLocation()
  const here = useRef(location)
  useEffect(() => {
    here.current = location
  })

  useEffect(() => {
    const handle = (error: unknown) => {
      if (!(error instanceof ApiError) || error.status !== 401) return
      const { pathname, search } = here.current
      if (pathname === '/admin/login') return
      qc.removeQueries({ queryKey: adminKeys.all })
      navigate(`/admin/login?redirect=${encodeURIComponent(pathname + search)}`, { replace: true })
    }
    const unsubQueries = qc.getQueryCache().subscribe((e) => {
      if (e.type === 'updated' && e.action.type === 'error') handle(e.action.error)
    })
    const unsubMutations = qc.getMutationCache().subscribe((e) => {
      if (e.type === 'updated' && e.action.type === 'error') handle(e.action.error)
    })
    return () => {
      unsubQueries()
      unsubMutations()
    }
  }, [qc, navigate])
}

export function RequireAuth() {
  const me = useMe()
  if (me.data) return <AdminLayout user={me.data} />
  if (me.isError && !(me.error instanceof ApiError && me.error.status === 401)) {
    return <p role="alert" className="p-8 text-destructive">{me.error.message}</p>
  }
  return <p className="p-8 text-muted-foreground">載入中…</p>
}
```

`web/src/admin/LoginPage.tsx`：

```tsx
import { type FormEvent } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate, useSearchParams } from 'react-router'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useDocumentTitle } from '@/site/useDocumentTitle'
import { adminApi, adminKeys } from './api'
import { safeRedirect } from './auth'

export default function LoginPage() {
  useDocumentTitle('後台登入')
  const [params] = useSearchParams()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const login = useMutation({
    mutationFn: (v: { account: string; password: string }) => adminApi.login(v.account, v.password),
    onSuccess: (user) => {
      qc.setQueryData(adminKeys.me, user)
      navigate(safeRedirect(params.get('redirect')), { replace: true })
    },
  })

  const onSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = new FormData(e.currentTarget)
    login.mutate({ account: String(form.get('account')), password: String(form.get('password')) })
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-muted p-4">
      <Card className="w-full max-w-sm">
        <CardHeader className="items-center text-center">
          <img src="/logo.png" alt="" className="mx-auto h-14" />
          <CardTitle>生長之家後台</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="account">帳號</Label>
              <Input id="account" name="account" autoComplete="username" required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">密碼</Label>
              <Input id="password" name="password" type="password" autoComplete="current-password" required />
            </div>
            {login.isError && (
              <p role="alert" className="text-sm text-destructive">
                {login.error.message}
              </p>
            )}
            <Button type="submit" className="w-full" disabled={login.isPending}>
              {login.isPending ? '登入中…' : '登入'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
```

`web/src/admin/AdminLayout.tsx`：

```tsx
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { NavLink, Outlet, useNavigate } from 'react-router'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import type { User } from '@/shared/types'
import { adminApi, adminKeys } from './api'

const links = [
  { to: '/admin/groups', label: '頁籤與頁面' },
  { to: '/admin/carousels', label: '輪播圖' },
  { to: '/admin/marquees', label: '跑馬燈' },
  { to: '/admin/images', label: '圖庫' },
  { to: '/admin/settings', label: '網站設定' },
]

export default function AdminLayout({ user }: { user: User }) {
  const qc = useQueryClient()
  const navigate = useNavigate()
  const logout = useMutation({
    mutationFn: adminApi.logout,
    onSuccess: () => {
      qc.removeQueries({ queryKey: adminKeys.all })
      navigate('/admin/login', { replace: true })
    },
    onError: (e) => toast.error(e.message),
  })

  return (
    <div className="min-h-screen bg-muted/40 md:flex">
      <aside className="border-b bg-background md:sticky md:top-0 md:h-screen md:w-56 md:shrink-0 md:border-b-0 md:border-r">
        <div className="flex items-center gap-2 p-4">
          <img src="/logo.png" alt="" className="h-8" />
          <span className="font-bold">後台管理</span>
        </div>
        <nav aria-label="後台選單" className="flex gap-1 overflow-x-auto px-2 pb-2 md:flex-col md:pb-0">
          {links.map((l) => (
            <NavLink
              key={l.to}
              to={l.to}
              className={({ isActive }) =>
                `whitespace-nowrap rounded-md px-3 py-2 text-sm ${isActive ? 'bg-primary text-primary-foreground' : 'hover:bg-muted'}`
              }
            >
              {l.label}
            </NavLink>
          ))}
        </nav>
        <div className="hidden space-y-2 border-t p-4 text-sm md:absolute md:bottom-0 md:block md:w-full">
          <p className="text-muted-foreground">{user.name}</p>
          <div className="flex gap-2">
            <Button variant="outline" size="sm" asChild>
              <a href="/" target="_blank" rel="noreferrer">查看網站</a>
            </Button>
            <Button variant="ghost" size="sm" onClick={() => logout.mutate()} disabled={logout.isPending}>
              登出
            </Button>
          </div>
        </div>
      </aside>
      <main className="min-w-0 flex-1 p-4 md:p-8">
        <div className="mb-4 flex items-center justify-end gap-2 text-sm md:hidden">
          <span className="text-muted-foreground">{user.name}</span>
          <Button variant="ghost" size="sm" onClick={() => logout.mutate()}>
            登出
          </Button>
        </div>
        <Outlet />
      </main>
    </div>
  )
}
```

`web/src/admin/ConfirmDialog.tsx`：

```tsx
import type { ReactNode } from 'react'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { buttonVariants } from '@/components/ui/button'

type Props = {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  description: ReactNode
  confirmLabel: string
  destructive?: boolean
  onConfirm: () => void
}

export default function ConfirmDialog({ open, onOpenChange, title, description, confirmLabel, destructive, onConfirm }: Props) {
  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{title}</AlertDialogTitle>
          <AlertDialogDescription asChild>
            <div>{description}</div>
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>取消</AlertDialogCancel>
          <AlertDialogAction
            className={destructive ? buttonVariants({ variant: 'destructive' }) : undefined}
            onClick={onConfirm}
          >
            {confirmLabel}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
```

`web/src/admin/AdminApp.tsx`（整個取代；各頁面在後續 task 取代佔位）：

```tsx
import { Navigate, Route, Routes } from 'react-router'
import { Toaster } from '@/components/ui/sonner'
import { RequireAuth, useRedirectOn401 } from './auth'
import LoginPage from './LoginPage'

const Todo = ({ name }: { name: string }) => <p>{name}（建置中）</p>

export default function AdminApp() {
  useRedirectOn401()
  return (
    <>
      <Routes>
        <Route path="login" element={<LoginPage />} />
        <Route element={<RequireAuth />}>
          <Route index element={<Navigate to="groups" replace />} />
          <Route path="groups" element={<Todo name="頁籤與頁面" />} />
          <Route path="pages/new" element={<Todo name="新增頁面" />} />
          <Route path="pages/:id" element={<Todo name="編輯頁面" />} />
          <Route path="carousels" element={<Todo name="輪播圖" />} />
          <Route path="marquees" element={<Todo name="跑馬燈" />} />
          <Route path="images" element={<Todo name="圖庫" />} />
          <Route path="settings" element={<Todo name="網站設定" />} />
          <Route path="*" element={<Navigate to="groups" replace />} />
        </Route>
      </Routes>
      <Toaster richColors position="top-center" />
    </>
  )
}
```

- [ ] **Step 6: 確認測試通過**

Run: `pnpm --dir web test && pnpm --dir web typecheck`
Expected: PASS

目視：`make dev` 後開 `http://localhost:5173/admin`，應導到登入頁；以 `admin` / `admin1234` 登入後看到側欄；登出回到登入頁。

- [ ] **Step 7: Commit**

```bash
git add web
git commit -m "[feat] 後台基礎：shadcn/ui、登入、401 導向與版面"
```

---

### Task 2: 頁籤與頁面管理（拖曳排序）

**Files:**
- Create: `web/src/admin/useAction.ts`、`web/src/admin/groups/order.ts`、`SortableList.tsx`、`NameDialog.tsx`、`GroupsPage.tsx`
- Modify: `web/src/admin/AdminApp.tsx`（`groups` 路由改用 `GroupsPage`）
- Test: `web/src/admin/groups/order.test.ts`、`web/src/admin/groups/GroupsPage.test.tsx`

**Interfaces:**
- Consumes: `adminApi`、`adminKeys`、`ConfirmDialog`（Task 1）
- Produces:
  ```ts
  // useAction.ts：執行任意 API 呼叫，成功後 invalidate 指定 key，失敗時 toast 錯誤
  export function useAction(invalidate: QueryKey): UseMutationResult<unknown, Error, () => Promise<unknown>>
  // order.ts
  export function applyOrder<T>(items: T[], ids: number[], getId: (item: T) => number): T[]
  // SortableList.tsx
  export default function SortableList<T>(props: { label: string; items: T[]; getId(item: T): number; onReorder(ids: number[]): void; renderItem(item: T, handle: ReactNode): ReactNode })
  ```

- [ ] **Step 1: 寫失敗的測試**

`web/src/admin/groups/order.test.ts`：

```ts
import { expect, it } from 'vitest'
import { applyOrder } from './order'

it('依 id 陣列重新排列，忽略不存在的 id', () => {
  const items = [{ id: 1 }, { id: 2 }, { id: 3 }]
  expect(applyOrder(items, [3, 9, 1, 2], (i) => i.id)).toEqual([{ id: 3 }, { id: 1 }, { id: 2 }])
})
```

`web/src/admin/groups/GroupsPage.test.tsx`：

```tsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, RouterProvider } from 'react-router'
import { afterEach, beforeEach, expect, it, vi } from 'vitest'
import { Toaster } from '@/components/ui/sonner'
import GroupsPage from './GroupsPage'

const json = (status: number, body: unknown) => new Response(body === null ? null : JSON.stringify(body), { status })

let groups: { id: number; name: string; pages: { id: number; name: string }[] }[]
let fetchMock: ReturnType<typeof vi.fn>

beforeEach(() => {
  groups = [
    { id: 1, name: '首頁', pages: [{ id: 1, name: '最新消息' }] },
    { id: 2, name: '主要行事活動', pages: [{ id: 3, name: '練成會' }, { id: 4, name: '誌友會' }] },
  ]
  fetchMock = vi.fn(async (url: string, init?: RequestInit) => {
    const method = init?.method ?? 'GET'
    if (url === '/api/v1/admin/groups' && method === 'GET') return json(200, { groups })
    if (url === '/api/v1/admin/groups' && method === 'POST') {
      const g = { id: 9, name: JSON.parse(String(init?.body)).name, pages: [] }
      groups = [...groups, g]
      return json(201, { group: g })
    }
    if (url === '/api/v1/admin/groups/2' && method === 'DELETE') {
      return json(409, { error: { code: 'group_not_empty', message: '此頁籤仍有頁面，請先移除或刪除頁面' } })
    }
    return json(404, { error: { code: 'not_found', message: '找不到' } })
  })
  vi.stubGlobal('fetch', fetchMock)
})

afterEach(() => vi.unstubAllGlobals())

function renderPage() {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  const router = createMemoryRouter([{ path: '/admin/groups', element: <GroupsPage /> }], {
    initialEntries: ['/admin/groups'],
  })
  render(
    <QueryClientProvider client={qc}>
      <RouterProvider router={router} />
      <Toaster />
    </QueryClientProvider>,
  )
  return userEvent.setup()
}

it('列出頁籤，選擇後顯示該頁籤的頁面', async () => {
  const user = renderPage()
  const groupList = await screen.findByRole('list', { name: '頁籤' })
  expect(within(groupList).getAllByRole('listitem')).toHaveLength(2)
  expect(screen.getByRole('link', { name: '最新消息' })).toHaveAttribute('href', '/admin/pages/1')

  await user.click(within(groupList).getByRole('button', { name: /^主要行事活動/ }))
  const pageList = screen.getByRole('list', { name: '頁面' })
  expect(within(pageList).getByRole('link', { name: '練成會' })).toHaveAttribute('href', '/admin/pages/3')
  expect(within(pageList).getByRole('link', { name: '誌友會' })).toHaveAttribute('href', '/admin/pages/4')
  expect(screen.getByRole('link', { name: '新增頁面' })).toHaveAttribute('href', '/admin/pages/new?group=2')
})

it('新增頁籤後自動選取', async () => {
  const user = renderPage()
  await user.click(await screen.findByRole('button', { name: '新增頁籤' }))
  await user.type(screen.getByLabelText('名稱'), '關於我們')
  await user.click(screen.getByRole('button', { name: '儲存' }))
  expect(await screen.findByText('已新增頁籤')).toBeInTheDocument()
  expect(await screen.findByText('此頁籤還沒有頁面。')).toBeInTheDocument()
  expect(fetchMock).toHaveBeenCalledWith('/api/v1/admin/groups', expect.objectContaining({ method: 'POST', body: '{"name":"關於我們"}' }))
})

it('刪除失敗時顯示 API 訊息；首頁頁籤不可刪除', async () => {
  const user = renderPage()
  expect(await screen.findByRole('button', { name: '刪除 首頁' })).toBeDisabled()
  await user.click(screen.getByRole('button', { name: '刪除 主要行事活動' }))
  await user.click(within(screen.getByRole('alertdialog')).getByRole('button', { name: '刪除' }))
  expect(await screen.findByText('此頁籤仍有頁面，請先移除或刪除頁面')).toBeInTheDocument()
})
```

Run: `pnpm --dir web test`
Expected: FAIL（找不到模組）

- [ ] **Step 2: 實作共用工具**

`web/src/admin/useAction.ts`：

```ts
import { useMutation, useQueryClient, type QueryKey } from '@tanstack/react-query'
import { toast } from 'sonner'

// 後台多數操作的共同模式：呼叫 API → 成功後重新載入清單 → 失敗時顯示錯誤
export function useAction(invalidate: QueryKey) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (fn: () => Promise<unknown>) => fn(),
    onSuccess: () => qc.invalidateQueries({ queryKey: invalidate }),
    onError: (e) => toast.error(e.message),
  })
}
```

`web/src/admin/groups/order.ts`：

```ts
export function applyOrder<T>(items: T[], ids: number[], getId: (item: T) => number): T[] {
  const byId = new Map(items.map((item) => [getId(item), item]))
  return ids.flatMap((id) => byId.get(id) ?? [])
}
```

`web/src/admin/groups/SortableList.tsx`：

```tsx
import type { ReactNode } from 'react'
import {
  closestCenter,
  DndContext,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
} from '@dnd-kit/core'
import { restrictToVerticalAxis } from '@dnd-kit/modifiers'
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'

type Props<T> = {
  label: string
  items: T[]
  getId: (item: T) => number
  onReorder: (ids: number[]) => void
  renderItem: (item: T, handle: ReactNode) => ReactNode
}

const announcements = {
  onDragStart: () => '已拿起項目，使用方向鍵移動，空白鍵放下',
  onDragOver: () => '',
  onDragEnd: () => '已放下項目',
  onDragCancel: () => '已取消移動',
}

export default function SortableList<T>({ label, items, getId, onReorder, renderItem }: Props<T>) {
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 4 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  )
  const ids = items.map(getId)

  const onDragEnd = ({ active, over }: DragEndEvent) => {
    if (!over || active.id === over.id) return
    onReorder(arrayMove(ids, ids.indexOf(Number(active.id)), ids.indexOf(Number(over.id))))
  }

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      modifiers={[restrictToVerticalAxis]}
      onDragEnd={onDragEnd}
      accessibility={{ announcements, screenReaderInstructions: { draggable: '按空白鍵拿起，方向鍵移動，再按空白鍵放下' } }}
    >
      <SortableContext items={ids} strategy={verticalListSortingStrategy}>
        <ul aria-label={label} className="space-y-2">
          {items.map((item) => (
            <SortableRow key={getId(item)} id={getId(item)}>
              {(handle) => renderItem(item, handle)}
            </SortableRow>
          ))}
        </ul>
      </SortableContext>
    </DndContext>
  )
}

function SortableRow({ id, children }: { id: number; children: (handle: ReactNode) => ReactNode }) {
  const { attributes, listeners, setNodeRef, setActivatorNodeRef, transform, transition, isDragging } = useSortable({ id })
  const handle = (
    <button
      type="button"
      ref={setActivatorNodeRef}
      {...attributes}
      {...listeners}
      aria-label="拖曳排序"
      className="cursor-grab touch-none rounded px-1 text-muted-foreground hover:bg-muted active:cursor-grabbing"
    >
      ⋮⋮
    </button>
  )
  return (
    <li
      ref={setNodeRef}
      style={{ transform: CSS.Transform.toString(transform), transition }}
      className={isDragging ? 'relative z-10 opacity-80 shadow-lg' : undefined}
    >
      {children(handle)}
    </li>
  )
}
```

`web/src/admin/groups/NameDialog.tsx`：

```tsx
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
```

- [ ] **Step 3: 實作 GroupsPage**

`web/src/admin/groups/GroupsPage.tsx`：

```tsx
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
```

`web/src/admin/AdminApp.tsx`：`groups` 路由改為 `<Route path="groups" element={<GroupsPage />} />`，並 `import GroupsPage from './groups/GroupsPage'`。

- [ ] **Step 4: 確認測試通過**

Run: `pnpm --dir web test && pnpm --dir web typecheck`
Expected: PASS

目視：`make dev`，在 `/admin/groups` 拖曳頁籤與頁面（滑鼠與鍵盤：Tab 到「⋮⋮」、空白鍵、方向鍵、空白鍵），重新整理後順序保留，前台選單同步。

- [ ] **Step 5: Commit**

```bash
git add web/src/admin
git commit -m "[feat] 後台頁籤與頁面管理（拖曳排序、新增、改名、刪除）"
```

---

### Task 3: 編輯器基礎（Tiptap 擴充、iframe 白名單、遺失偵測）

**Files:**
- Create: `web/src/admin/editor/iframeHosts.ts`、`extensions.ts`、`lossDetection.ts`
- Test: `web/src/admin/editor/iframeHosts.test.ts`、`extensions.test.ts`、`lossDetection.test.ts`

**Interfaces:**
- Produces:
  ```ts
  // iframeHosts.ts
  export const IFRAME_HOSTS: string[]
  export function isAllowedIframe(src: string): boolean
  export function toEmbedUrl(input: string): string          // YouTube 網址 / 整段 <iframe> → 可嵌入的 src
  // extensions.ts
  export function createExtensions(): AnyExtension[]
  export function roundTrip(html: string): string            // 經 Tiptap 解析再輸出
  export function removeStyleProps(style: string | null, props: string[]): string | null
  //   editor.commands.setIframe({ src })（型別擴充）
  // lossDetection.ts
  export function detectLoss(original: string, roundTripped: string): string[]  // 遺失的標籤、屬性、樣式；文字不同時含「文字內容」
  ```

- [ ] **Step 1: 寫失敗的測試**

`web/src/admin/editor/iframeHosts.test.ts`：

```ts
import { describe, expect, it } from 'vitest'
import { isAllowedIframe, toEmbedUrl } from './iframeHosts'

describe('isAllowedIframe', () => {
  it.each([
    ['https://www.youtube.com/embed/abc', true],
    ['//www.youtube.com/embed/abc', true],
    ['https://players.brightcove.net/1/x/index.html?videoId=2', true],
    ['https://drive.google.com/file/d/x/preview', true],
    ['https://www.youtube.com.evil.com/embed', false],
    ['https://evil.example/embed', false],
    ['/embed/abc', false],
    ['javascript:alert(1)', false],
    ['', false],
  ])('%s → %s', (src, want) => {
    expect(isAllowedIframe(src)).toBe(want)
  })
})

describe('toEmbedUrl', () => {
  it.each([
    ['https://www.youtube.com/watch?v=V4LhnGIPZt4&t=10s', 'https://www.youtube.com/embed/V4LhnGIPZt4'],
    ['https://youtu.be/V4LhnGIPZt4', 'https://www.youtube.com/embed/V4LhnGIPZt4'],
    ['https://www.youtube.com/shorts/abc123', 'https://www.youtube.com/embed/abc123'],
    ['  https://www.youtube.com/embed/x  ', 'https://www.youtube.com/embed/x'],
    ['<iframe width="560" src="https://www.youtube.com/embed/x?si=1" frameborder="0"></iframe>', 'https://www.youtube.com/embed/x?si=1'],
    ['https://drive.google.com/file/d/x/preview', 'https://drive.google.com/file/d/x/preview'],
  ])('%s', (input, want) => {
    expect(toEmbedUrl(input)).toBe(want)
  })
})
```

`web/src/admin/editor/lossDetection.test.ts`：

```ts
import { describe, expect, it } from 'vitest'
import { detectLoss } from './lossDetection'

describe('detectLoss', () => {
  it('等價的 HTML 不算遺失', () => {
    expect(detectLoss('<p><b>粗</b>&nbsp;<i>斜</i></p>', '<p><strong>粗</strong> <em>斜</em></p>')).toEqual([])
    expect(
      detectLoss(
        '<p style="color: #e03e2d; font-size: 14pt">x</p>',
        '<p style="font-size:14pt;color:rgb(224, 62, 45)">x</p>',
      ),
    ).toEqual([])
    expect(detectLoss('<table><tr><td>1</td></tr></table>', '<table><colgroup><col></colgroup><tbody><tr><td><p>1</p></td></tr></tbody></table>')).toEqual([])
    expect(detectLoss('<p class="a b">x</p>', '<p class="b a">x</p>')).toEqual([])
  })

  it('找出遺失的標籤、屬性、樣式與文字', () => {
    expect(detectLoss('<div class="box"><p>x</p></div>', '<p>x</p>')).toEqual(['<div>', 'class="box"'])
    expect(detectLoss('<table border="1"><tr><td>x</td></tr></table>', '<table><tr><td>x</td></tr></table>')).toEqual(['border="1"'])
    expect(detectLoss('<p style="font-weight: bold">x</p>', '<p>x</p>')).toEqual(['font-weight: bold'])
    expect(detectLoss('<p>一段文字</p>', '<p>一段</p>')).toEqual(['文字內容'])
    expect(detectLoss('<p><span>a</span><span>b</span></p>', '<p><span>ab</span></p>')).toEqual(['<span>'])
  })
})
```

`web/src/admin/editor/extensions.test.ts`：

```ts
import { describe, expect, it } from 'vitest'
import { removeStyleProps, roundTrip } from './extensions'
import { detectLoss } from './lossDetection'

// 取自舊站實際內容的寫法
const legacy = `<p><a title="傳道協會FB" href="https://www.facebook.com/seichonoie.tw"><img src="/php/picture/2019-12-21_18-41-17.jpg"></a>&nbsp;</p>
<p style="text-align: left;"><span style="font-size: 14pt; color: #e03e2d; font-family: arial, helvetica, sans-serif;">●官方網站重新改版</span></p>
<p><span style="background-color: #fbeeb8;">每週日上午舉行。</span></p>
<h2 style="text-align: center;">標題</h2>
<table style="border-collapse: collapse; width: 100%;" border="1"><tbody>
<tr style="height: 34px;"><td style="width: 100%; background-color: #18a085; text-align: center;"><strong><span style="color: #ffffff;">國際本部發行的影片</span></strong></td></tr>
<tr><td><iframe src="//www.youtube.com/embed/V4LhnGIPZt4" width="100%" height="225" allowfullscreen="allowfullscreen"></iframe></td></tr>
</tbody></table>
<ul><li>項目一</li><li><u>項目二</u></li></ul>`

describe('roundTrip', () => {
  it('舊站常見格式不會遺失', () => {
    expect(detectLoss(legacy, roundTrip(legacy))).toEqual([])
  })

  it('不在白名單的 iframe 會被移除', () => {
    const out = roundTrip('<iframe src="https://evil.example/x"></iframe><p>後面</p>')
    expect(out).not.toContain('iframe')
    expect(detectLoss('<iframe src="https://evil.example/x"></iframe>', out)).toContain('<iframe>')
  })

  it('無法保留的結構會被偵測到', () => {
    const html = '<div class="box"><p>框內</p></div><p><font color="red">舊字型</font></p>'
    const lost = detectLoss(html, roundTrip(html))
    expect(lost).toEqual(expect.arrayContaining(['<div>', '<font>']))
  })

  it('連結不會被加上 target 與 rel', () => {
    expect(roundTrip('<p><a href="/page/3">練成會</a></p>')).toBe('<p><a href="/page/3">練成會</a></p>')
  })
})

describe('removeStyleProps', () => {
  it.each([
    [null, ['color'], null],
    ['color: red; font-weight: bold;', ['color'], 'font-weight: bold;'],
    ['text-align: left', ['text-align'], null],
    ['COLOR:red;;  width:1px', ['color'], 'width:1px;'],
    ['background: url(a;b)', [], 'background: url(a;b);'],
  ])('%s', (style, props, want) => {
    expect(removeStyleProps(style, props)).toBe(want)
  })
})
```

> `background: url(a;b)` 這種值內含分號的情況以括號深度判斷，不能單純 `split(';')`。

Run: `pnpm --dir web test`
Expected: FAIL（找不到模組）

- [ ] **Step 2: 實作 iframeHosts 與 lossDetection**

`web/src/admin/editor/iframeHosts.ts`：

```ts
// 需與後端 backend/internal/content/sanitize.go 的 IframeHosts 一致
export const IFRAME_HOSTS = [
  'www.youtube.com',
  'youtube.com',
  'www.youtube-nocookie.com',
  'players.brightcove.net',
  'drive.google.com',
  'www.facebook.com',
]

export function isAllowedIframe(src: string): boolean {
  try {
    // 相對網址會解析到 invalid 主機，自然不在白名單內
    const url = new URL(src, 'https://relative.invalid')
    return /^https?:$/.test(url.protocol) && IFRAME_HOSTS.includes(url.hostname)
  } catch {
    return false
  }
}

export function toEmbedUrl(input: string): string {
  const trimmed = input.trim()
  const raw = /<iframe[^>]*\ssrc=["']([^"']+)["']/i.exec(trimmed)?.[1] ?? trimmed
  let url: URL
  try {
    url = new URL(raw)
  } catch {
    return raw
  }
  const host = url.hostname.replace(/^(www|m)\./, '')
  if (host === 'youtube.com' && url.pathname === '/watch' && url.searchParams.get('v')) {
    return `https://www.youtube.com/embed/${url.searchParams.get('v')}`
  }
  if (host === 'youtube.com' && url.pathname.startsWith('/shorts/')) {
    return `https://www.youtube.com/embed/${url.pathname.slice('/shorts/'.length)}`
  }
  if (host === 'youtu.be') return `https://www.youtube.com/embed${url.pathname}`
  return raw
}
```

`web/src/admin/editor/lossDetection.ts`：

```ts
// 視為相同的標籤（Tiptap 會把 b 輸出成 strong 等）
const aliases: Record<string, string> = { b: 'strong', i: 'em', strike: 's', del: 's' }
// 解析或 Tiptap 會自動補上／調整的結構
const ignoredTags = new Set(['tbody', 'colgroup', 'col'])

// 交給瀏覽器正規化 CSS（#e03e2d 與 rgb(224, 62, 45) 視為相同）
function styleDeclarations(style: string): string[] {
  const el = document.createElement('div')
  el.style.cssText = style
  return Array.from(el.style, (prop) => `${prop}: ${el.style.getPropertyValue(prop)}`)
}

function signature(html: string) {
  const doc = new DOMParser().parseFromString(`<body>${html}</body>`, 'text/html')
  const counts = new Map<string, number>()
  const add = (key: string) => counts.set(key, (counts.get(key) ?? 0) + 1)
  for (const el of doc.body.querySelectorAll('*')) {
    const tag = aliases[el.localName] ?? el.localName
    if (ignoredTags.has(tag)) continue
    add(`<${tag}>`)
    for (const attr of el.attributes) {
      if (attr.name === 'style') styleDeclarations(attr.value).forEach(add)
      else if (attr.name === 'class') attr.value.split(/\s+/).filter(Boolean).forEach((c) => add(`class="${c}"`))
      else add(`${attr.name}="${attr.value}"`)
    }
  }
  const text = (doc.body.textContent ?? '').replace(/\s+/g, ' ').trim()
  return { counts, text }
}

export function detectLoss(original: string, roundTripped: string): string[] {
  const before = signature(original)
  const after = signature(roundTripped)
  const lost: string[] = []
  for (const [key, n] of before.counts) {
    if ((after.counts.get(key) ?? 0) < n) lost.push(key)
  }
  if (before.text !== after.text) lost.push('文字內容')
  return lost
}
```

- [ ] **Step 3: 實作 extensions**

`web/src/admin/editor/extensions.ts`：

```ts
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
```

- [ ] **Step 4: 確認測試通過**

Run: `pnpm --dir web test -- editor`
Expected: PASS。

若 `舊站常見格式不會遺失` 失敗，依輸出的遺失清單調整 `PreserveAttributes` 或擴充設定（例如 Tiptap 版本的屬性名稱不同），**不要修改測試的 legacy 內容**。已知可能的差異與處理方式：
- `<a>` 被加上 `rel`、`target`：確認 `link.HTMLAttributes` 設為 `null` 有生效。
- 表格多了 `min-width` 樣式：屬於新增，不算遺失，不需處理。
- 儲存格內文字被包進 `<p>`：屬於新增，不需處理。

- [ ] **Step 5: Commit**

```bash
git add web/src/admin/editor
git commit -m "[feat] Tiptap 擴充設定、iframe 白名單與格式遺失偵測"
```

---

### Task 4: 圖庫與圖片選擇器

**Files:**
- Create: `web/src/admin/images/useUploader.ts`、`UploadWidgets.tsx`、`ImagesPage.tsx`、`ImagePickerDialog.tsx`
- Modify: `web/src/admin/AdminApp.tsx`（`images` 路由改用 `ImagesPage`）
- Test: `web/src/admin/images/useUploader.test.tsx`

**Interfaces:**
- Consumes: `uploadImage`、`adminApi`、`adminKeys`（Task 1）、`useAction`、`ConfirmDialog`
- Produces:
  ```ts
  export type UploadItem = { key: string; name: string; progress: number; error?: string; done?: boolean }
  export function useUploader(): { items: UploadItem[]; upload(files: File[]): Promise<ImageItem[]>; clear(): void }
  export function DropZone(props: { onFiles(files: File[]): void }): JSX.Element
  export function UploadList(props: { items: UploadItem[]; onClear(): void }): JSX.Element | null
  export default function ImagePickerDialog(props: { open: boolean; onOpenChange(open: boolean): void; onSelect(url: string): void })
  ```

- [ ] **Step 1: 寫失敗的測試**

`web/src/admin/images/useUploader.test.tsx`：

```tsx
import type { ReactNode } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { act, renderHook } from '@testing-library/react'
import { expect, it, vi } from 'vitest'
import { ApiError } from '@/shared/api'
import { useUploader } from './useUploader'

const uploadImage = vi.hoisted(() => vi.fn())
vi.mock('../api', async (importOriginal) => ({ ...(await importOriginal<object>()), uploadImage }))

const wrapper = ({ children }: { children: ReactNode }) => (
  <QueryClientProvider client={new QueryClient()}>{children}</QueryClientProvider>
)

const file = (name: string, type: string, size = 10) => new File([new Uint8Array(size)], name, { type })

it('依序上傳、回報進度，並在前端先擋下不合格的檔案', async () => {
  uploadImage.mockImplementation(async (f: File, onProgress: (r: number) => void) => {
    if (f.name === 'server-reject.png') throw new ApiError(400, 'unsupported_type', '不是支援的圖片格式')
    onProgress(0.5)
    return { name: `saved-${f.name}`, url: `/php/picture/saved-${f.name}`, size: f.size, mod_time: '' }
  })
  const { result } = renderHook(() => useUploader(), { wrapper })

  let done: Awaited<ReturnType<typeof result.current.upload>> = []
  await act(async () => {
    done = await result.current.upload([
      file('a.png', 'image/png'),
      file('big.jpg', 'image/jpeg', 5 * 1024 * 1024 + 1),
      file('doc.pdf', 'application/pdf'),
      file('server-reject.png', 'image/png'),
    ])
  })

  expect(done.map((d) => d.name)).toEqual(['saved-a.png'])
  expect(uploadImage).toHaveBeenCalledTimes(2)
  expect(result.current.items).toEqual([
    expect.objectContaining({ name: 'a.png', progress: 1, done: true }),
    expect.objectContaining({ name: 'big.jpg', error: '檔案超過 5MB' }),
    expect.objectContaining({ name: 'doc.pdf', error: '只能上傳 jpg、png、gif、webp 圖片' }),
    expect.objectContaining({ name: 'server-reject.png', error: '不是支援的圖片格式' }),
  ])

  act(() => result.current.clear())
  expect(result.current.items).toEqual([])
})
```

Run: `pnpm --dir web test -- useUploader`
Expected: FAIL（找不到模組）

- [ ] **Step 2: 實作上傳 hook 與元件**

`web/src/admin/images/useUploader.ts`：

```ts
import { useCallback, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import type { ImageItem } from '@/shared/types'
import { adminKeys, uploadImage } from '../api'

export type UploadItem = { key: string; name: string; progress: number; error?: string; done?: boolean }

const MAX_BYTES = 5 * 1024 * 1024
const TYPES = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
let seq = 0

export function useUploader() {
  const qc = useQueryClient()
  const [items, setItems] = useState<UploadItem[]>([])

  const upload = useCallback(
    async (files: File[]) => {
      const patch = (key: string, p: Partial<UploadItem>) =>
        setItems((list) => list.map((item) => (item.key === key ? { ...item, ...p } : item)))
      const queue = files.map((file) => ({ file, key: String(++seq) }))
      setItems((list) => [...list, ...queue.map(({ file, key }) => ({ key, name: file.name, progress: 0 }))])

      const saved: ImageItem[] = []
      // 依序上傳：進度清楚，也避免同一秒大量請求
      for (const { file, key } of queue) {
        if (!TYPES.includes(file.type)) {
          patch(key, { error: '只能上傳 jpg、png、gif、webp 圖片' })
          continue
        }
        if (file.size > MAX_BYTES) {
          patch(key, { error: '檔案超過 5MB' })
          continue
        }
        try {
          const image = await uploadImage(file, (progress) => patch(key, { progress }))
          patch(key, { progress: 1, done: true })
          saved.push(image)
        } catch (e) {
          patch(key, { error: (e as Error).message })
        }
      }
      if (saved.length) await qc.invalidateQueries({ queryKey: adminKeys.images })
      return saved
    },
    [qc],
  )

  const clear = useCallback(() => setItems([]), [])
  return { items, upload, clear }
}
```

`web/src/admin/images/UploadWidgets.tsx`：

```tsx
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import type { UploadItem } from './useUploader'

export function DropZone({ onFiles }: { onFiles: (files: File[]) => void }) {
  const [over, setOver] = useState(false)
  return (
    <label
      onDragOver={(e) => {
        e.preventDefault()
        setOver(true)
      }}
      onDragLeave={() => setOver(false)}
      onDrop={(e) => {
        e.preventDefault()
        setOver(false)
        onFiles([...e.dataTransfer.files])
      }}
      className={`flex cursor-pointer flex-col items-center gap-1 rounded-md border-2 border-dashed p-6 text-center text-sm transition-colors focus-within:ring-2 focus-within:ring-ring ${over ? 'border-primary bg-primary/5' : 'border-border'}`}
    >
      <span className="font-medium">拖曳圖片到這裡，或點擊選擇檔案</span>
      <span className="text-muted-foreground">jpg、png、gif、webp，每個檔案 5MB 以內</span>
      <input
        type="file"
        multiple
        accept="image/jpeg,image/png,image/gif,image/webp"
        aria-label="選擇圖片"
        className="sr-only"
        onChange={(e) => {
          onFiles([...(e.target.files ?? [])])
          e.target.value = ''
        }}
      />
    </label>
  )
}

export function UploadList({ items, onClear }: { items: UploadItem[]; onClear: () => void }) {
  if (items.length === 0) return null
  const busy = items.some((i) => !i.done && !i.error)
  return (
    <div className="space-y-2 rounded-md border bg-background p-3 text-sm">
      <ul className="space-y-2">
        {items.map((item) => (
          <li key={item.key}>
            <div className="flex justify-between gap-2">
              <span className="truncate">{item.name}</span>
              <span className={item.error ? 'text-destructive' : 'text-muted-foreground'}>
                {item.error ?? (item.done ? '已完成' : `${Math.round(item.progress * 100)}%`)}
              </span>
            </div>
            {!item.error && !item.done && (
              <div className="mt-1 h-1.5 overflow-hidden rounded bg-muted" role="progressbar" aria-label={item.name} aria-valuenow={Math.round(item.progress * 100)}>
                <div className="h-full bg-primary transition-[width]" style={{ width: `${item.progress * 100}%` }} />
              </div>
            )}
          </li>
        ))}
      </ul>
      {!busy && (
        <Button type="button" size="sm" variant="ghost" onClick={onClear}>
          清除上傳紀錄
        </Button>
      )}
    </div>
  )
}
```

- [ ] **Step 3: 圖庫頁與選擇器**

`web/src/admin/images/ImagesPage.tsx`：

```tsx
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
```

`web/src/admin/images/ImagePickerDialog.tsx`：

```tsx
import { useQuery } from '@tanstack/react-query'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { adminApi, adminKeys } from '../api'
import { useUploader } from './useUploader'
import { DropZone, UploadList } from './UploadWidgets'

type Props = { open: boolean; onOpenChange: (open: boolean) => void; onSelect: (url: string) => void }

export default function ImagePickerDialog({ open, onOpenChange, onSelect }: Props) {
  const images = useQuery({ queryKey: adminKeys.images, queryFn: adminApi.images, enabled: open })
  const uploader = useUploader()

  const uploadAndSelect = async (files: File[]) => {
    const saved = await uploader.upload(files)
    if (saved.length === 1) onSelect(saved[0].url)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-3xl">
        <DialogHeader>
          <DialogTitle>選擇圖片</DialogTitle>
          <DialogDescription>點選圖庫中的圖片，或上傳一張新圖片直接使用。</DialogDescription>
        </DialogHeader>
        <DropZone onFiles={uploadAndSelect} />
        <UploadList items={uploader.items} onClear={uploader.clear} />
        {images.isError && <p role="alert" className="text-sm text-destructive">{images.error.message}</p>}
        <ul className="grid max-h-[50vh] grid-cols-3 gap-2 overflow-y-auto sm:grid-cols-4">
          {images.data?.map((image) => (
            <li key={image.name}>
              <button
                type="button"
                onClick={() => onSelect(image.url)}
                className="block w-full overflow-hidden rounded border hover:ring-2 hover:ring-primary focus-visible:ring-2 focus-visible:ring-primary"
              >
                <img src={image.url} alt={image.name} loading="lazy" className="aspect-square w-full bg-muted object-cover" />
              </button>
            </li>
          ))}
        </ul>
      </DialogContent>
    </Dialog>
  )
}
```

`web/src/admin/AdminApp.tsx`：`images` 路由改為 `<Route path="images" element={<ImagesPage />} />`，並 `import ImagesPage from './images/ImagesPage'`。

- [ ] **Step 4: 確認測試通過**

Run: `pnpm --dir web test && pnpm --dir web typecheck`
Expected: PASS

目視：`/admin/images` 拖曳多個檔案，逐一顯示進度；超過 5MB 或 PDF 直接顯示錯誤；刪除 `seed` 中被輪播圖使用的 `sample-1.jpg` 時，對話框列出「輪播圖 #1」（按取消，不要真的刪除）。

- [ ] **Step 5: Commit**

```bash
git add web/src/admin
git commit -m "[feat] 後台圖庫：多檔上傳進度、引用檢查與圖片選擇器"
```

---

### Task 5: 頁面編輯器

**Files:**
- Create: `web/src/admin/editor/Toolbar.tsx`、`VideoDialog.tsx`、`PageEditorPage.tsx`
- Modify: `web/src/admin/AdminApp.tsx`（`pages/new`、`pages/:id` 改用 `PageEditorPage`）

**Interfaces:**
- Consumes: `createExtensions`、`roundTrip`、`detectLoss`、`toEmbedUrl`、`isAllowedIframe`、`IFRAME_HOSTS`（Task 3）、`ImagePickerDialog`（Task 4）、`ConfirmDialog`、`adminApi`
- Produces: `/admin/pages/new?group=:gid` 與 `/admin/pages/:id` 頁面。可存取名稱（E2E 使用）：
  - 欄位：`頁面名稱`、`所屬頁籤`；編輯區 `頁面內容`（contenteditable）；原始碼區為 CodeMirror（`.cm-content`）
  - 按鈕：`儲存`、`預覽`、`返回列表`、`視覺編輯`、`HTML 原始碼`，工具列：`粗體`、`斜體`、`底線`、`刪除線`、`連結`、`插入圖片`、`插入影片`、`表格`、`復原`、`重做`
  - 遺失警告：`role="status"`，內容含「HTML 原始碼模式」

整體流程（載入、切換模式、儲存、離開提示）由 Task 7 的 E2E 驗證；個別規則已在 Task 3 單元測試。

- [ ] **Step 1: 工具列**

`web/src/admin/editor/Toolbar.tsx`：

```tsx
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
```

`web/src/admin/editor/VideoDialog.tsx`：

```tsx
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
```

- [ ] **Step 2: 編輯頁**

`web/src/admin/editor/PageEditorPage.tsx`：

```tsx
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
```

> `setDirty(false)` 會先更新 `dirtyRef`，因此新增頁面後緊接著的 `navigate` 不會被 `useBlocker` 攔下。

`web/src/admin/AdminApp.tsx`：`pages/new` 與 `pages/:id` 路由都改為 `element={<PageEditorPage />}`，並 `import PageEditorPage from './editor/PageEditorPage'`。

- [ ] **Step 3: 型別檢查與目視確認**

Run: `pnpm --dir web test && pnpm --dir web typecheck && pnpm --dir web build`
Expected: PASS；建置成功，且 `web/dist/assets/` 中 Tiptap 相關程式碼只出現在後台 chunk：

```bash
grep -l ProseMirror web/dist/assets/*.js
grep -o 'src="/assets/[^"]*\.js"' web/dist/index.html
```

第一個指令列出的檔案不可是第二個指令列出的入口檔。

目視（`make dev`）：
1. `/admin/pages/3`（練成會）以視覺模式開啟，表格的框線與寬度保留；工具列各按鈕可用；插入 YouTube 影片與圖庫圖片。
2. 以原始碼模式貼上 `<div class="box"><p>x</p></div>`，按「視覺編輯」會出現確認對話框。
3. 修改後點側欄「輪播圖」會出現「還有變更沒有儲存」；關閉分頁時瀏覽器會提示。
4. 新增頁面儲存後網址變成 `/admin/pages/<id>`，不會出現離開提示。

- [ ] **Step 4: Commit**

```bash
git add web/src/admin
git commit -m "[feat] 後台頁面編輯器：Tiptap 工具列、HTML 原始碼模式與未存檔提示"
```

---

### Task 6: 輪播圖、跑馬燈與網站設定

**Files:**
- Create: `web/src/admin/carousels/CarouselsPage.tsx`、`web/src/admin/marquees/MarqueesPage.tsx`、`web/src/admin/settings/SettingsPage.tsx`
- Modify: `web/src/admin/AdminApp.tsx`（三個路由改用新頁面，移除 `Todo`）
- Test: `web/src/admin/marquees/MarqueesPage.test.tsx`、`web/src/admin/settings/SettingsPage.test.tsx`

**Interfaces:**
- Consumes: `adminApi`、`adminKeys`、`useAction`、`ConfirmDialog`、`ImagePickerDialog`
- Produces: `/admin/carousels`、`/admin/marquees`、`/admin/settings`

- [ ] **Step 1: 寫失敗的測試**

兩個測試共用的渲染方式直接寫在各自檔案中。

`web/src/admin/marquees/MarqueesPage.test.tsx`：

```tsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { fireEvent, render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router'
import { afterEach, expect, it, vi } from 'vitest'
import { Toaster } from '@/components/ui/sonner'
import MarqueesPage from './MarqueesPage'

const json = (status: number, body: unknown) => new Response(JSON.stringify(body), { status })

afterEach(() => vi.unstubAllGlobals())

it('新增跑馬燈時即時預覽文字與顏色，儲存後顯示結果', async () => {
  let marquees = [{ id: 1, text: '合掌感謝！', color: '#1EFF00' }]
  const fetchMock = vi.fn(async (url: string, init?: RequestInit) => {
    if (init?.method === 'POST') {
      const m = { id: 2, ...JSON.parse(String(init.body)) }
      marquees = [...marquees, m]
      return json(201, { marquee: m })
    }
    return json(200, { marquees })
  })
  vi.stubGlobal('fetch', fetchMock)

  render(
    <QueryClientProvider client={new QueryClient()}>
      <MemoryRouter>
        <MarqueesPage />
      </MemoryRouter>
      <Toaster />
    </QueryClientProvider>,
  )
  const user = userEvent.setup()
  expect(await screen.findByText('合掌感謝！')).toBeInTheDocument()

  await user.click(screen.getByRole('button', { name: '新增跑馬燈' }))
  const dialog = screen.getByRole('dialog')
  await user.type(within(dialog).getByLabelText('文字'), '歡迎參加練成會')
  fireEvent.input(within(dialog).getByLabelText('顏色'), { target: { value: '#ffcc00' } })
  const preview = within(dialog).getByTestId('marquee-preview')
  expect(preview).toHaveTextContent('歡迎參加練成會')
  expect(preview).toHaveStyle({ color: '#ffcc00' })

  await user.click(within(dialog).getByRole('button', { name: '儲存' }))
  expect(await screen.findByText('已儲存跑馬燈')).toBeInTheDocument()
  expect(fetchMock).toHaveBeenCalledWith(
    '/api/v1/admin/marquees',
    expect.objectContaining({ method: 'POST', body: JSON.stringify({ text: '歡迎參加練成會', color: '#ffcc00' }) }),
  )
  expect(await screen.findByText('歡迎參加練成會')).toBeInTheDocument()
})
```

`web/src/admin/settings/SettingsPage.test.tsx`：

```tsx
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, expect, it, vi } from 'vitest'
import { Toaster } from '@/components/ui/sonner'
import SettingsPage from './SettingsPage'

const json = (status: number, body: unknown) => new Response(JSON.stringify(body), { status })

afterEach(() => vi.unstubAllGlobals())

function setup(putResponse: Response) {
  const fetchMock = vi.fn(async (_url: string, init?: RequestInit) =>
    init?.method === 'PUT'
      ? putResponse
      : json(200, { settings: { web_title: '生長之家', web_sub_title: '一句簡單的感謝', facebook_url: '' } }),
  )
  vi.stubGlobal('fetch', fetchMock)
  render(
    <QueryClientProvider client={new QueryClient()}>
      <SettingsPage />
      <Toaster />
    </QueryClientProvider>,
  )
  return { fetchMock, user: userEvent.setup() }
}

it('載入並儲存網站設定', async () => {
  const saved = { web_title: '生長之家台灣', web_sub_title: '一句簡單的感謝', facebook_url: 'https://www.facebook.com/seichonoie.tw' }
  const { fetchMock, user } = setup(json(200, { settings: saved }))

  const title = await screen.findByLabelText('網站名稱')
  expect(title).toHaveValue('生長之家')
  await user.clear(title)
  await user.type(title, '生長之家台灣')
  await user.type(screen.getByLabelText('Facebook 粉絲專頁網址'), 'https://www.facebook.com/seichonoie.tw')
  await user.click(screen.getByRole('button', { name: '儲存設定' }))

  expect(await screen.findByText('已儲存設定')).toBeInTheDocument()
  expect(fetchMock).toHaveBeenLastCalledWith('/api/v1/admin/settings', expect.objectContaining({ method: 'PUT', body: JSON.stringify(saved) }))
})

it('顯示伺服器的驗證訊息', async () => {
  const { user } = setup(json(400, { error: { code: 'validation', message: 'Facebook 網址需為 http(s):// 開頭' } }))
  await user.click(await screen.findByRole('button', { name: '儲存設定' }))
  expect(await screen.findByText('Facebook 網址需為 http(s):// 開頭')).toBeInTheDocument()
})
```

Run: `pnpm --dir web test -- Marquees Settings`
Expected: FAIL（找不到模組）

- [ ] **Step 2: 跑馬燈頁**

`web/src/admin/marquees/MarqueesPage.tsx`：

```tsx
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
```

- [ ] **Step 3: 輪播圖頁**

`web/src/admin/carousels/CarouselsPage.tsx`：

```tsx
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
```

- [ ] **Step 4: 網站設定頁**

`web/src/admin/settings/SettingsPage.tsx`：

```tsx
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
```

- [ ] **Step 5: 更新路由**

`web/src/admin/AdminApp.tsx`（整個取代）：

```tsx
import { Navigate, Route, Routes } from 'react-router'
import { Toaster } from '@/components/ui/sonner'
import { RequireAuth, useRedirectOn401 } from './auth'
import CarouselsPage from './carousels/CarouselsPage'
import PageEditorPage from './editor/PageEditorPage'
import GroupsPage from './groups/GroupsPage'
import ImagesPage from './images/ImagesPage'
import LoginPage from './LoginPage'
import MarqueesPage from './marquees/MarqueesPage'
import SettingsPage from './settings/SettingsPage'

export default function AdminApp() {
  useRedirectOn401()
  return (
    <>
      <Routes>
        <Route path="login" element={<LoginPage />} />
        <Route element={<RequireAuth />}>
          <Route index element={<Navigate to="groups" replace />} />
          <Route path="groups" element={<GroupsPage />} />
          <Route path="pages/new" element={<PageEditorPage />} />
          <Route path="pages/:id" element={<PageEditorPage />} />
          <Route path="carousels" element={<CarouselsPage />} />
          <Route path="marquees" element={<MarqueesPage />} />
          <Route path="images" element={<ImagesPage />} />
          <Route path="settings" element={<SettingsPage />} />
          <Route path="*" element={<Navigate to="groups" replace />} />
        </Route>
      </Routes>
      <Toaster richColors position="top-center" />
    </>
  )
}
```

`auth.test.tsx` 中對 `/admin/settings` 的 mock 已回傳設定資料，此時會渲染真正的 `SettingsPage`，測試仍應通過。

- [ ] **Step 6: 確認測試通過**

Run: `pnpm --dir web test && pnpm --dir web typecheck`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add web/src/admin
git commit -m "[feat] 後台輪播圖、跑馬燈與網站設定"
```

---

### Task 7: 後台 E2E 與整體驗證

**Files:**
- Create: `web/e2e/admin.spec.ts`

**Interfaces:**
- Consumes: 前台 plan 的 `web/playwright.config.ts` 與 `make e2e`；`seed.sql`（`admin` / `admin1234`；頁籤順序 首頁、關於我們、主要行事活動）

- [ ] **Step 1: 寫 E2E 測試**

`web/e2e/admin.spec.ts`：

```ts
import path from 'node:path'
import { expect, test, type Page } from '@playwright/test'

const origin = { Origin: 'http://localhost:8080' }

async function login(page: Page, password = 'admin1234') {
  await page.goto('/admin/groups')
  await expect(page).toHaveURL(/\/admin\/login\?redirect=%2Fadmin%2Fgroups/)
  await page.getByLabel('帳號').fill('admin')
  await page.getByLabel('密碼').fill(password)
  await page.getByRole('button', { name: '登入' }).click()
}

test('密碼錯誤時顯示訊息', async ({ page }) => {
  await login(page, 'wrong-password')
  await expect(page.getByRole('alert')).toHaveText('帳號或密碼錯誤')
})

test('新增頁面並在前台看到', async ({ page }) => {
  await login(page)
  await expect(page).toHaveURL('/admin/groups')
  const name = `E2E 頁面 ${Date.now()}`

  await page.getByRole('list', { name: '頁籤' }).getByRole('button', { name: /^關於我們/ }).click()
  await page.getByRole('link', { name: '新增頁面' }).click()
  await page.getByLabel('頁面名稱').fill(name)
  await expect(page.getByLabel('所屬頁籤')).toHaveValue('3')

  await page.getByLabel('頁面內容').click()
  await page.keyboard.type('這是 E2E 建立的內容，')
  await page.getByRole('button', { name: '粗體' }).click()
  await page.keyboard.type('粗體字')
  await page.getByRole('button', { name: '儲存', exact: true }).click()

  await expect(page.getByText('已儲存')).toBeVisible()
  await expect(page).toHaveURL(/\/admin\/pages\/\d+$/)
  const id = page.url().split('/').pop()

  await page.goto(`/page/${id}`)
  await expect(page.getByRole('heading', { name })).toBeVisible()
  await expect(page.locator('.legacy-content strong')).toHaveText('粗體字')
})

test('無法保留格式的頁面以原始碼模式開啟並可直接儲存', async ({ page }) => {
  await login(page)
  const res = await page.request.post('/api/v1/admin/pages', {
    headers: origin,
    data: { name: `E2E 原始碼 ${Date.now()}`, group_id: 3, html: '<div class="box"><p>框內文字</p></div>' },
  })
  expect(res.ok()).toBeTruthy()
  const { page: created } = await res.json()

  await page.goto(`/admin/pages/${created.id}`)
  await expect(page.getByRole('status')).toContainText('HTML 原始碼模式')
  await expect(page.getByRole('button', { name: 'HTML 原始碼' })).toHaveAttribute('aria-pressed', 'true')

  await page.locator('.cm-content').click()
  await page.keyboard.press('ControlOrMeta+End')
  await page.keyboard.insertText('<p>新增段落</p>')
  await page.getByRole('button', { name: '儲存', exact: true }).click()
  await expect(page.getByText('已儲存')).toBeVisible()

  const saved = await (await page.request.get(`/api/v1/pages/${created.id}`)).json()
  expect(saved.page.html).toContain('<div class="box">')
  expect(saved.page.html).toContain('<p>新增段落</p>')

  await page.getByRole('button', { name: '視覺編輯' }).click()
  await expect(page.getByRole('alertdialog')).toContainText('無法在視覺編輯器中保留')
  await page.getByRole('button', { name: '取消' }).click()
  await expect(page.getByRole('button', { name: 'HTML 原始碼' })).toHaveAttribute('aria-pressed', 'true')
})

test('未儲存時離開會提示', async ({ page }) => {
  await login(page)
  await page.goto('/admin/pages/4')
  await page.getByLabel('頁面名稱').fill('誌友會（修改中）')
  await page.getByRole('navigation', { name: '後台選單' }).getByRole('link', { name: '輪播圖' }).click()
  await expect(page.getByRole('alertdialog')).toContainText('還有變更沒有儲存')
  await page.getByRole('button', { name: '取消' }).click()
  await expect(page).toHaveURL('/admin/pages/4')
})

test('以鍵盤排序頁籤並反映到前台', async ({ page }) => {
  await login(page)
  const items = () => page.getByRole('list', { name: '頁籤' }).getByRole('listitem')
  await expect(items()).toHaveText([/首頁/, /關於我們/, /主要行事活動/])

  await items().filter({ hasText: '關於我們' }).getByRole('button', { name: '拖曳排序' }).focus()
  await page.keyboard.press('Space')
  await page.keyboard.press('ArrowDown')
  await page.keyboard.press('Space')
  await expect(items()).toHaveText([/首頁/, /主要行事活動/, /關於我們/])

  await expect
    .poll(async () => (await (await page.request.get('/api/v1/site')).json()).menu.map((g: { name: string }) => g.name))
    .toEqual(['首頁', '主要行事活動', '關於我們'])
  await page.reload()
  await expect(items()).toHaveText([/首頁/, /主要行事活動/, /關於我們/])
})

test('上傳圖片後出現在圖庫', async ({ page }) => {
  await login(page)
  await page.getByRole('navigation', { name: '後台選單' }).getByRole('link', { name: '圖庫' }).click()
  await page.getByLabel('選擇圖片').setInputFiles(path.resolve(import.meta.dirname, '../public/logo.png'))
  await expect(page.getByText('已完成')).toBeVisible()
  await expect(page.getByRole('img', { name: /^\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}.*\.png$/ }).first()).toBeVisible()
})

test('登出後需要重新登入', async ({ page }) => {
  await login(page)
  await page.getByRole('button', { name: '登出' }).first().click()
  await expect(page).toHaveURL('/admin/login')
  await page.goto('/admin/settings')
  await expect(page).toHaveURL(/\/admin\/login\?redirect=/)
})

test('前台入口 bundle 不含後台程式', async ({ request }) => {
  const html = await (await request.get('/')).text()
  const entry = /<script type="module" crossorigin src="([^"]+)"/.exec(html)?.[1]
  expect(entry).toBeTruthy()
  const js = await (await request.get(entry!)).text()
  for (const marker of ['ProseMirror', 'DndDescribedBy', 'cm-content']) {
    expect(js, marker).not.toContain(marker)
  }
})
```

> 前台 E2E（`site.spec.ts`）不依賴頁籤順序，因此即使後台測試先執行也不受影響；`make e2e` 每次都會重設資料庫。

- [ ] **Step 2: 執行 E2E**

Run: `make up && make e2e`
Expected: `admin.spec.ts` 與 `site.spec.ts` 全部 PASS。

- [ ] **Step 3: 整體驗證**

Run:

```bash
make lint
make test
make build
docker build -t sniweb:dev .
```

Expected: 全部成功。

最後以正式建置做一次手動檢查：

```bash
set -a; . ./.env; set +a
PUBLIC_BASE_URL=http://localhost:8080 DEV_ORIGINS= ./bin/sniweb serve
```

開 `http://localhost:8080/`、`/#/3`、`/page/3`、`/admin`，確認前台與後台都正常；在瀏覽器開發者工具的 Console 確認沒有 CSP 錯誤（Google Fonts、YouTube iframe、GA）。

- [ ] **Step 4: Commit**

```bash
git add web/e2e/admin.spec.ts
git commit -m "[test] 後台 E2E：登入、編輯頁面、原始碼模式、排序、圖庫與 bundle 檢查"
```

- [ ] **Step 5: 收尾**

使用 superpowers:finishing-a-development-branch 決定合併或開 PR。

---

## 自我檢查紀錄

- spec 第 8 節：登入與 401 導向（Task 1）、左側選單（Task 1）、頁籤與頁面雙欄拖曳排序與 dialog（Task 2）、頁面編輯的欄位／Tiptap 擴充／工具列／原始碼模式／遺失偵測／預覽／未存檔提示（Task 3、5）、輪播圖／跑馬燈／圖庫／網站設定（Task 4、6）、toast（全部）。
- spec 第 13 節：Tiptap 遺失偵測單元測試（Task 3）；後台 E2E 的登入 → 編輯頁面 → 前台看到更新、拖曳排序（Task 7）。
- 型別一致性：`adminApi`、`adminKeys`、`useAction`、`ConfirmDialog`、`ImagePickerDialog`、`createExtensions` / `roundTrip` / `detectLoss` 的名稱與簽章在各 task 中一致。
