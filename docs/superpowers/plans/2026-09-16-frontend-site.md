# 生長之家網站重寫 — 前台（React）Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 建立 `web/` 前端專案與前台網站（首頁、內頁、選單、輪播、跑馬燈、footer、GA4），並完成 Docker 建置。

**Architecture:** Vite + React 19 SPA，由 Go 服務以 embed 提供。前台以 TanStack Query 讀取公開 API；`/admin/*` 以 `React.lazy` 載入（後台由第三份 plan 實作，本 plan 只放佔位元件）。舊站 HTML 內容原樣輸出在 `.legacy-content` 容器內，渲染後包裝表格並攔截站內連結。

**Tech Stack:** Vite、React 19、TypeScript、React Router 7、TanStack Query 5、Tailwind CSS 4、Embla Carousel、Vitest + Testing Library、Playwright、pnpm。

**Spec:** `docs/superpowers/specs/2026-09-16-sni-web-rewrite-design.md`

**前置條件：** `docs/superpowers/plans/2026-09-16-backend.md` 已全部完成（API、`<!--sni:head-->` 注入、`#sni-config`、`backend/internal/site/webdist/dist/`、`docker-compose.yml`、`Makefile`）。

## Global Constraints

- Node.js 24+、pnpm 10（`web/package.json` 的 `packageManager` 固定版本）。
- 所有 API 呼叫走 `/api/v1`，同網域，`credentials: 'same-origin'`；錯誤格式 `{"error":{"code","message","request_id?"}}`，`message` 可直接顯示。
- 前台 bundle 不可包含後台程式碼與 Tiptap：`/admin/*` 一律 `React.lazy`。
- `web/index.html` 的 `<head>` 必須含 `<!--sni:head-->`，且不可自帶 `<title>`。
- 前端公開設定從 `<script id="sni-config" type="application/json">` 讀取（`gaMeasurementId`、`counterScriptUrl`），沒有時視為空字串。
- 路由：`/` 首頁、`/page/:id` 內頁、其他顯示 404；舊網址 `#/(\d+)` → `/page/$1`、`#/` → `/`（`replace`）。
- 內文連結：同主機名稱（不分 http/https）且為舊 hash 路由或 `/page/:id` 時改用前端路由；其他連結維持原行為。
- iframe 以 16:9 自適應；`img { max-width: 100%; height: auto }`；每個 `<table>` 包在 `overflow-x: auto` 容器內。
- 視覺依下方「設計方向」；手機寬度不可出現橫向捲動。
- 跑馬燈每則用自己的 `color`，hover 暫停；`prefers-reduced-motion` 時改為靜態輪替。
- 所有程式先寫失敗的測試再實作（純 UI 樣式除外）。
- Commit 訊息沿用 repo 慣例：`[feat] ...`、`[chore] ...`、`[test] ...`、`[docs] ...`。

## 設計方向（依 frontend-design skill 訂定）

- **對象與任務：** 多為年長信眾、常用手機瀏覽；網站的工作是讓人找到活動資訊、安靜地讀完文章。
- **色彩：** 寶藍 `#2230B8`（logo，連結與互動）、夜藍 `#10185C`（跑馬燈、footer）、晨光紙色 `#F7F8FC`（頁面底色）、墨色 `#1C2233`（文字）、淡藍 `#E9ECFA`（hover 與骨架）、橄欖綠 `#6E8B3D`（取自鴿子銜的橄欖枝，只用於非文字標記，例如目前頁籤的底線）。
- **字體：** 標題 Noto Serif TC（明體，帶教義閱讀感）；內文 Noto Sans TC 18px、行高 1.85；內文欄寬 `42rem`（約 40 個中文字）。
- **版面：** 單欄、靠左對齊的閱讀版面；header 白底，logo＋明體站名，桌機在站名旁顯示副標題；輪播在手機滿版、桌機小圓角且無陰影。
- **唯一的大膽之處：** 首頁在輪播與跑馬燈下方，以大號明體排出網站副標題（預設「一句簡單的感謝」）作為問候語。其他元素保持安靜。
- **不做：** 漸層底色、卡片陰影、膠囊按鈕、全大寫標籤、「大數字＋小標」的 404。動態只保留使用者預期的輪播、跑馬燈與選單展開。
- **品質底線：** 所有互動元素有清楚的 `:focus-visible` 外框；尊重 `prefers-reduced-motion`；文字對比至少 4.5:1。

## 本 plan 的實作層決定

- 只用 `tsc --noEmit` 當前端 lint（不加 ESLint），`make lint` 會一併執行。
- 不引入 UI 元件庫到前台；後台才使用 shadcn/ui（第三份 plan）。
- 本機與 E2E 範例圖片：`make sample-images` 把 `legacy/img/bg1.jpg`、`bg3.jpg` 複製成 `tmp/picture/sample-1.jpg`、`sample-2.jpg`（對應 `seed.sql` 的輪播圖）。

## 檔案結構

```
web/
  package.json / pnpm-lock.yaml
  index.html
  vite.config.ts
  tsconfig.json
  playwright.config.ts
  public/favicon.png / logo.png / logo-text.png
  src/
    main.tsx                 進入點：舊 hash 轉址、QueryClient、Router
    App.tsx                  /admin/* lazy、其他交給 SiteApp
    index.css                Tailwind、主題色、.legacy-content
    test/setup.ts
    shared/
      api.ts(+test)          fetch 包裝與 ApiError
      config.ts(+test)       讀取 #sni-config
      types.ts               API 型別
    site/
      legacyLinks.ts(+test)  舊 hash 與內文連結判斷
      wrapTables.ts(+test)
      queries.ts             useSite / useHome / usePage
      analytics.ts(+test)    GA4
      useDocumentTitle.ts
      SiteApp.tsx            版面與路由
      Header.tsx(+test)      桌機下拉選單、手機抽屜
      Banner.tsx
      Marquee.tsx(+test)
      LegacyContent.tsx(+test)
      PageView.tsx           首頁／內頁共用：banner + 跑馬燈 + 內容
      Footer.tsx(+test)
      NotFound.tsx
    admin/
      AdminApp.tsx           佔位（第三份 plan 取代）
  e2e/
    site.spec.ts
Dockerfile
Makefile                     補上 web、dev、sample-images、db-reset、e2e
```

---

### Task 1: web 專案骨架、建置串接與 Dockerfile

**Files:**
- Create: `web/package.json`、`web/index.html`、`web/vite.config.ts`、`web/tsconfig.json`、`web/src/main.tsx`、`web/src/App.tsx`、`web/src/index.css`、`web/src/test/setup.ts`、`web/src/test/smoke.test.ts`、`web/src/admin/AdminApp.tsx`、`web/public/*`
- Create: `Dockerfile`
- Modify: `Makefile`

**Interfaces:**
- Produces: `pnpm --dir web build` 產出 `web/dist/`；`make web` 會把產物複製到 `backend/internal/site/webdist/dist/`；path alias `@/` → `web/src/`；`web/src/admin/AdminApp.tsx` 預設匯出一個元件（第三份 plan 取代內容）。

- [ ] **Step 1: 建立 package.json 並安裝套件**

`web/package.json`：

```json
{
  "name": "sniweb-web",
  "private": true,
  "type": "module",
  "packageManager": "pnpm@10.30.3",
  "scripts": {
    "dev": "vite",
    "build": "tsc --noEmit && vite build",
    "typecheck": "tsc --noEmit",
    "test": "vitest run",
    "e2e": "playwright test"
  }
}
```

```bash
cd web
pnpm add react react-dom react-router @tanstack/react-query embla-carousel-react embla-carousel-autoplay
pnpm add -D vite @vitejs/plugin-react typescript @types/react @types/react-dom @types/node \
  tailwindcss @tailwindcss/vite vitest jsdom @testing-library/react @testing-library/user-event \
  @testing-library/jest-dom @playwright/test
cd ..
```

- [ ] **Step 2: 設定檔**

`web/vite.config.ts`：

```ts
import { fileURLToPath } from 'node:url'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { defineConfig } from 'vitest/config'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
  },
  server: {
    // 保留瀏覽器原本的 Origin（http://localhost:5173），後端以 DEV_ORIGINS 允許
    proxy: {
      '/api': 'http://localhost:8080',
      '/php/picture': 'http://localhost:8080',
    },
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    include: ['src/**/*.test.{ts,tsx}'],
  },
})
```

`web/tsconfig.json`：

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "lib": ["ES2023", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "moduleResolution": "bundler",
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noEmit": true,
    "skipLibCheck": true,
    "isolatedModules": true,
    "verbatimModuleSyntax": true,
    "types": ["vite/client", "node"],
    "baseUrl": ".",
    "paths": { "@/*": ["./src/*"] }
  },
  "include": ["src", "e2e", "vite.config.ts", "playwright.config.ts"]
}
```

`web/index.html`：

```html
<!doctype html>
<html lang="zh-Hant-TW">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <link rel="icon" type="image/png" href="/favicon.png" />
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
    <link href="https://fonts.googleapis.com/css2?family=Noto+Sans+TC:wght@400;500;700&family=Noto+Serif+TC:wght@600;700&display=swap" rel="stylesheet" />
    <!--sni:head-->
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

`web/src/test/setup.ts`：

```ts
import '@testing-library/jest-dom/vitest'
```

- [ ] **Step 3: 素材**

```bash
mkdir -p web/public
cp legacy/favicon.png web/public/favicon.png
cp legacy/front/public/img/title_logo.png web/public/logo.png
cp legacy/front/public/img/title_text.png web/public/logo-text.png
```

- [ ] **Step 4: 最小可執行的程式與冒煙測試**

`web/src/test/smoke.test.ts`：

```ts
import { readFileSync } from 'node:fs'
import { expect, test } from 'vitest'

test('index.html 保留 meta 佔位符且沒有自帶 title', () => {
  const html = readFileSync(new URL('../../index.html', import.meta.url), 'utf8')
  expect(html).toContain('<!--sni:head-->')
  expect(html).not.toMatch(/<title>/i)
})
```

`web/src/index.css`：

```css
@import 'tailwindcss';

@theme {
  --font-sans: 'Noto Sans TC', system-ui, sans-serif;
  --font-serif: 'Noto Serif TC', 'Songti TC', serif;
  --color-brand: #2230b8;
  --color-brand-dark: #10185c;
  --color-brand-soft: #e9ecfa;
  --color-paper: #f7f8fc;
  --color-ink: #1c2233;
  --color-olive: #6e8b3d;
}

body {
  @apply bg-paper font-sans text-ink antialiased;
}

:focus-visible {
  outline: 2px solid var(--color-brand);
  outline-offset: 2px;
}

@media (prefers-reduced-motion: no-preference) {
  html {
    scroll-behavior: smooth;
  }
}
```

`web/src/admin/AdminApp.tsx`：

```tsx
export default function AdminApp() {
  return <p className="p-8">後台建置中</p>
}
```

`web/src/App.tsx`：

```tsx
import { lazy, Suspense } from 'react'
import { Route, Routes } from 'react-router'

const AdminApp = lazy(() => import('@/admin/AdminApp'))

export default function App() {
  return (
    <Routes>
      <Route
        path="/admin/*"
        element={
          <Suspense fallback={null}>
            <AdminApp />
          </Suspense>
        }
      />
      <Route path="*" element={<p className="p-8">前台建置中</p>} />
    </Routes>
  )
}
```

`web/src/main.tsx`：

```tsx
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import App from './App'
import './index.css'

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 60_000, retry: 1 } },
})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>,
)
```

Run: `pnpm --dir web test && pnpm --dir web build && grep -c 'sni:head' web/dist/index.html`
Expected: 測試 PASS；建置成功；`grep` 輸出 `1`（Vite 保留註解）。若 Vite 移除了註解，在 `vite.config.ts` 加入 `transformIndexHtml` plugin 於輸出後補回，並讓冒煙測試改檢查 `web/dist/index.html`。

- [ ] **Step 5: 更新 Makefile**

把 `Makefile` 整個改為（在後端版本上加入前端 target）：

```make
SQLC_VERSION ?= v1.30.0

.PHONY: up down gen test test-short lint build web dev sample-images db-reset e2e

up:
	docker compose up -d

down:
	docker compose down

gen:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION) generate

test:
	go test -race ./...
	pnpm --dir web test

test-short:
	go test -short ./...
	pnpm --dir web test

lint:
	@test -z "$$(gofmt -l backend)" || (gofmt -l backend; exit 1)
	go vet ./...
	pnpm --dir web typecheck

# 建置前端並複製到 Go embed 目錄
web:
	pnpm --dir web install --frozen-lockfile
	pnpm --dir web build
	find backend/internal/site/webdist/dist -mindepth 1 ! -name .gitkeep -exec rm -rf {} +
	cp -R web/dist/. backend/internal/site/webdist/dist/

build: web
	CGO_ENABLED=0 go build -trimpath -o bin/sniweb ./backend/cmd/sniweb

sample-images:
	mkdir -p tmp/picture
	cp legacy/img/bg1.jpg tmp/picture/sample-1.jpg
	cp legacy/img/bg3.jpg tmp/picture/sample-2.jpg

dev: sample-images
	@trap 'kill 0' INT TERM EXIT; \
	(set -a; . ./.env; set +a; go run ./backend/cmd/sniweb serve) & \
	pnpm --dir web dev & \
	wait

# 把本機資料庫重設為 seed 狀態（E2E 前使用）
db-reset:
	docker compose exec -T mysql sh -c '\
	  mysql -uroot -psniweb -e "DROP DATABASE IF EXISTS sniweb; CREATE DATABASE sniweb CHARACTER SET utf8mb4" && \
	  mysql -uroot -psniweb sniweb < /docker-entrypoint-initdb.d/01-schema.sql && \
	  mysql -uroot -psniweb sniweb < /docker-entrypoint-initdb.d/02-seed.sql'

e2e: web sample-images db-reset
	pnpm --dir web exec playwright install --with-deps chromium
	pnpm --dir web e2e
```

Run: `make web && ls backend/internal/site/webdist/dist && go build ./...`
Expected: 列出 `.gitkeep`、`index.html`、`assets`、`favicon.png` 等；Go 建置成功。

- [ ] **Step 6: Dockerfile**

`Dockerfile`：

```dockerfile
FROM node:24-alpine AS web
WORKDIR /src/web
RUN corepack enable
COPY web/package.json web/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY web/ ./
RUN pnpm build

FROM golang:1.25-alpine AS server
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY backend/ ./backend/
COPY --from=web /src/web/dist/ ./backend/internal/site/webdist/dist/
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/sniweb ./backend/cmd/sniweb

# 使用 root 映像，才能寫入 Zeabur 掛載的 /data Volume
FROM gcr.io/distroless/static-debian12
COPY --from=server /out/sniweb /sniweb
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/sniweb"]
CMD ["serve"]
```

Run:

```bash
docker build -t sniweb:dev .
docker run --rm sniweb:dev; echo "exit=$?"
```

Expected: 建置成功；容器預設執行 `serve`，因缺少環境變數印出「缺少環境變數：DATABASE_URL, PUBLIC_BASE_URL」且 `exit=1`。

- [ ] **Step 7: Commit**

```bash
git add web Makefile Dockerfile
git commit -m "[chore] 建立 web 前端專案、建置串接與 Dockerfile"
```

---

### Task 2: 共用 API client、公開設定與型別

**Files:**
- Create: `web/src/shared/api.ts`、`web/src/shared/config.ts`、`web/src/shared/types.ts`
- Test: `web/src/shared/api.test.ts`、`web/src/shared/config.test.ts`

**Interfaces:**
- Produces:
  ```ts
  // api.ts
  export class ApiError extends Error { status: number; code: string; requestId?: string }
  export function api<T>(path: string, opts?: { method?: string; body?: unknown; signal?: AbortSignal }): Promise<T>
  //   path 不含 /api/v1 前綴；body 為 FormData 時原樣送出；204 回傳 undefined
  //   網路錯誤 → ApiError(status 0, code 'network')

  // config.ts
  export type PublicConfig = { gaMeasurementId: string; counterScriptUrl: string }
  export function readPublicConfig(doc?: Document): PublicConfig
  export const publicConfig: PublicConfig

  // types.ts
  export type Settings = { web_title: string; web_sub_title: string; facebook_url: string }
  export type PageRef = { id: number; name: string }
  export type Group = { id: number; name: string; pages: PageRef[] }
  export type Page = { id: number; group_id: number; name: string; html: string }
  export type Carousel = { id: number; image: string; url: string }
  export type Marquee = { id: number; text: string; color: string }
  export type SiteData = Settings & { menu: Group[] }
  export type PageData = { page: Page | null; carousels: Carousel[]; marquees: Marquee[] }
  export type ImageItem = { name: string; url: string; size: number; mod_time: string }
  export type ImageUsages = { pages: PageRef[]; carousels: Carousel[] }
  export type User = { id: number; account: string; name: string }
  ```

- [ ] **Step 1: 寫失敗的測試**

`web/src/shared/api.test.ts`：

```ts
import { afterEach, describe, expect, it, vi } from 'vitest'
import { api, ApiError } from './api'

function mockFetch(impl: (...args: Parameters<typeof fetch>) => Promise<Response>) {
  const fn = vi.fn(impl)
  vi.stubGlobal('fetch', fn)
  return fn
}

const json = (status: number, body: unknown) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } })

afterEach(() => vi.unstubAllGlobals())

describe('api', () => {
  it('送出 JSON 並解析回應', async () => {
    const fetchMock = mockFetch(async () => json(200, { ok: true }))
    await expect(api('/admin/groups', { method: 'POST', body: { name: 'x' } })).resolves.toEqual({ ok: true })
    const [url, init] = fetchMock.mock.calls[0]
    expect(url).toBe('/api/v1/admin/groups')
    expect(init?.method).toBe('POST')
    expect(init?.body).toBe('{"name":"x"}')
    expect(init?.credentials).toBe('same-origin')
    expect(new Headers(init?.headers).get('Content-Type')).toBe('application/json')
  })

  it('FormData 原樣送出且不設定 Content-Type', async () => {
    const fetchMock = mockFetch(async () => json(201, {}))
    const form = new FormData()
    await api('/admin/images', { method: 'POST', body: form })
    const init = fetchMock.mock.calls[0][1]
    expect(init?.body).toBe(form)
    expect(new Headers(init?.headers).has('Content-Type')).toBe(false)
  })

  it('204 回傳 undefined', async () => {
    mockFetch(async () => new Response(null, { status: 204 }))
    await expect(api('/admin/pages/1', { method: 'DELETE' })).resolves.toBeUndefined()
  })

  it('錯誤回應轉成 ApiError', async () => {
    mockFetch(async () => json(409, { error: { code: 'group_not_empty', message: '此頁籤仍有頁面' } }))
    const err = await api('/x').catch((e) => e)
    expect(err).toBeInstanceOf(ApiError)
    expect(err).toMatchObject({ status: 409, code: 'group_not_empty', message: '此頁籤仍有頁面' })
  })

  it('500 訊息附上錯誤代碼', async () => {
    mockFetch(async () => json(500, { error: { code: 'internal', message: '系統發生錯誤', request_id: 'abc123' } }))
    const err = await api('/x').catch((e) => e)
    expect(err.message).toBe('系統發生錯誤（代碼 abc123）')
    expect(err.requestId).toBe('abc123')
  })

  it('非 JSON 錯誤回應也有可讀訊息', async () => {
    mockFetch(async () => new Response('bad gateway', { status: 502 }))
    const err = await api('/x').catch((e) => e)
    expect(err).toMatchObject({ status: 502, code: 'unknown', message: '伺服器錯誤（502）' })
  })

  it('網路錯誤', async () => {
    mockFetch(async () => {
      throw new TypeError('Failed to fetch')
    })
    const err = await api('/x').catch((e) => e)
    expect(err).toMatchObject({ status: 0, code: 'network' })
  })

  it('取消請求時保留 AbortError', async () => {
    mockFetch(async () => {
      throw new DOMException('aborted', 'AbortError')
    })
    const err = await api('/x').catch((e) => e)
    expect(err.name).toBe('AbortError')
  })
})
```

`web/src/shared/config.test.ts`：

```ts
import { describe, expect, it } from 'vitest'
import { readPublicConfig } from './config'

function doc(html: string) {
  return new DOMParser().parseFromString(`<html><head>${html}</head></html>`, 'text/html')
}

describe('readPublicConfig', () => {
  it('讀取注入的設定', () => {
    const d = doc('<script id="sni-config" type="application/json">{"gaMeasurementId":"G-1","counterScriptUrl":"https://c/x"}</script>')
    expect(readPublicConfig(d)).toEqual({ gaMeasurementId: 'G-1', counterScriptUrl: 'https://c/x' })
  })

  it('沒有注入（vite dev）或格式錯誤時回傳空設定', () => {
    const empty = { gaMeasurementId: '', counterScriptUrl: '' }
    expect(readPublicConfig(doc('<!--sni:head-->'))).toEqual(empty)
    expect(readPublicConfig(doc('<script id="sni-config" type="application/json">{oops</script>'))).toEqual(empty)
  })
})
```

Run: `pnpm --dir web test`
Expected: FAIL（找不到模組）

- [ ] **Step 2: 實作**

`web/src/shared/api.ts`：

```ts
export class ApiError extends Error {
  status: number
  code: string
  requestId?: string

  constructor(status: number, code: string, message: string, requestId?: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.requestId = requestId
  }
}

type Options = { method?: string; body?: unknown; signal?: AbortSignal }

export async function api<T>(path: string, { method = 'GET', body, signal }: Options = {}): Promise<T> {
  const isForm = body instanceof FormData
  let res: Response
  try {
    res = await fetch(`/api/v1${path}`, {
      method,
      signal,
      credentials: 'same-origin',
      headers: body === undefined || isForm ? undefined : { 'Content-Type': 'application/json' },
      body: body === undefined ? undefined : isForm ? body : JSON.stringify(body),
    })
  } catch (e) {
    if (e instanceof DOMException && e.name === 'AbortError') throw e
    throw new ApiError(0, 'network', '無法連線到伺服器，請檢查網路後再試')
  }
  if (res.status === 204) return undefined as T

  const data = await res.json().catch(() => null)
  if (!res.ok) {
    const err = data?.error
    let message: string = err?.message ?? `伺服器錯誤（${res.status}）`
    if (err?.request_id) message += `（代碼 ${err.request_id}）`
    throw new ApiError(res.status, err?.code ?? 'unknown', message, err?.request_id)
  }
  return data as T
}
```

`web/src/shared/config.ts`：

```ts
export type PublicConfig = { gaMeasurementId: string; counterScriptUrl: string }

const empty: PublicConfig = { gaMeasurementId: '', counterScriptUrl: '' }

// 設定由 Go 服務注入 index.html；vite dev 時不存在
export function readPublicConfig(doc: Document = document): PublicConfig {
  const text = doc.getElementById('sni-config')?.textContent
  if (!text) return empty
  try {
    return { ...empty, ...JSON.parse(text) }
  } catch {
    return empty
  }
}

export const publicConfig = readPublicConfig()
```

`web/src/shared/types.ts`：

```ts
export type Settings = { web_title: string; web_sub_title: string; facebook_url: string }
export type PageRef = { id: number; name: string }
export type Group = { id: number; name: string; pages: PageRef[] }
export type Page = { id: number; group_id: number; name: string; html: string }
export type Carousel = { id: number; image: string; url: string }
export type Marquee = { id: number; text: string; color: string }
export type SiteData = Settings & { menu: Group[] }
export type PageData = { page: Page | null; carousels: Carousel[]; marquees: Marquee[] }
export type ImageItem = { name: string; url: string; size: number; mod_time: string }
export type ImageUsages = { pages: PageRef[]; carousels: Carousel[] }
export type User = { id: number; account: string; name: string }
```

- [ ] **Step 3: 確認測試通過**

Run: `pnpm --dir web test && pnpm --dir web typecheck`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add web/src/shared
git commit -m "[feat] 前端 API client、公開設定與型別"
```

---

### Task 3: 舊網址與內文連結判斷、表格包裝

**Files:**
- Create: `web/src/site/legacyLinks.ts`、`web/src/site/wrapTables.ts`
- Modify: `web/src/main.tsx`（render 前轉址）
- Test: `web/src/site/legacyLinks.test.ts`、`web/src/site/wrapTables.test.ts`

**Interfaces:**
- Produces:
  ```ts
  export function legacyHashTarget(hash: string): string | null     // '#/5' → '/page/5'；'#/' → '/'
  export function internalPath(href: string, current: { href: string; hostname: string }): string | null
  export function wrapTables(root: HTMLElement): void                // 每個 table 包進 div.table-scroll（重複呼叫不會重複包）
  ```

- [ ] **Step 1: 寫失敗的測試**

`web/src/site/legacyLinks.test.ts`：

```ts
import { describe, expect, it } from 'vitest'
import { internalPath, legacyHashTarget } from './legacyLinks'

describe('legacyHashTarget', () => {
  it.each([
    ['#/5', '/page/5'],
    ['#/75/', '/page/75'],
    ['#/', '/'],
    ['#', null],
    ['', null],
    ['#/abc', null],
    ['#top', null],
  ])('%s → %s', (hash, want) => {
    expect(legacyHashTarget(hash)).toBe(want)
  })
})

describe('internalPath', () => {
  const current = { href: 'https://www.seicho-no-ie.org.tw/page/3', hostname: 'www.seicho-no-ie.org.tw' }

  it.each([
    ['舊站絕對網址（http）', 'http://www.seicho-no-ie.org.tw/#/5', '/page/5'],
    ['舊站首頁', 'http://www.seicho-no-ie.org.tw/#/', '/'],
    ['相對舊 hash', '/#/20', '/page/20'],
    ['新網址', '/page/7', '/page/7'],
    ['新網址帶查詢字串', 'https://www.seicho-no-ie.org.tw/page/7?x=1', '/page/7'],
    ['站內首頁', '/', '/'],
    ['同頁錨點', '#section', null],
    ['其他站內路徑', '/php/picture/a.jpg', null],
    ['後台', '/admin/groups', null],
    ['外部網站', 'https://www.facebook.com/seichonoie.tw', null],
    ['其他子網域', 'http://seicho-no-ie.org.tw/#/5', null],
    ['mailto', 'mailto:a@b.c', null],
    ['javascript', 'javascript:alert(1)', null],
    ['空字串', '', null],
  ])('%s', (_, href, want) => {
    expect(internalPath(href, current)).toBe(want)
  })
})
```

> 「空字串」會解析成目前頁面 `/page/3`，不帶 hash；但空 href 不是導覽意圖，實作需直接回傳 `null`。

`web/src/site/wrapTables.test.ts`：

```ts
import { expect, it } from 'vitest'
import { wrapTables } from './wrapTables'

it('每個表格包進可橫向捲動的容器，重複呼叫不重複包', () => {
  const root = document.createElement('div')
  root.innerHTML = '<p>x</p><table id="a"><tr><td><table id="b"><tr><td>1</td></tr></table></td></tr></table>'
  wrapTables(root)
  wrapTables(root)
  for (const id of ['a', 'b']) {
    const table = root.querySelector(`#${id}`)!
    expect(table.parentElement?.className).toBe('table-scroll')
    expect(table.parentElement?.parentElement?.className).not.toBe('table-scroll')
  }
  expect(root.querySelectorAll('.table-scroll')).toHaveLength(2)
})
```

Run: `pnpm --dir web test`
Expected: FAIL（找不到模組）

- [ ] **Step 2: 實作**

`web/src/site/legacyLinks.ts`：

```ts
// 舊站使用 hash 路由：/#/{頁面 id}
export function legacyHashTarget(hash: string): string | null {
  if (hash === '#/') return '/'
  const m = /^#\/(\d+)\/?$/.exec(hash)
  return m ? `/page/${m[1]}` : null
}

// internalPath 判斷內文連結是否該用前端路由切換；是的話回傳目標路徑。
// 只比對主機名稱，讓舊內容中的 http:// 絕對網址在 https 站上也能使用。
export function internalPath(href: string, current: { href: string; hostname: string }): string | null {
  if (!href) return null
  let url: URL
  try {
    url = new URL(href, current.href)
  } catch {
    return null
  }
  if (!/^https?:$/.test(url.protocol) || url.hostname !== current.hostname) return null
  if (url.hash) {
    return url.pathname === '/' ? legacyHashTarget(url.hash) : null
  }
  if (url.pathname === '/') return '/'
  return /^\/page\/\d+$/.test(url.pathname) ? url.pathname : null
}
```

`web/src/site/wrapTables.ts`：

```ts
export function wrapTables(root: HTMLElement) {
  for (const table of root.querySelectorAll('table')) {
    if (table.parentElement?.classList.contains('table-scroll')) continue
    const wrapper = document.createElement('div')
    wrapper.className = 'table-scroll'
    table.replaceWith(wrapper)
    wrapper.append(table)
  }
}
```

`web/src/main.tsx`：在 `const queryClient` 之前加入（並 import `legacyHashTarget`）：

```tsx
import { legacyHashTarget } from '@/site/legacyLinks'

// 舊網址 /#/5 → /page/5；在 Router 啟動前改寫，避免多一次導覽
const legacyTarget = legacyHashTarget(window.location.hash)
if (legacyTarget) window.history.replaceState(null, '', legacyTarget)
```

- [ ] **Step 3: 確認測試通過**

Run: `pnpm --dir web test && pnpm --dir web typecheck`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add web/src
git commit -m "[feat] 舊 hash 網址轉址、內文連結判斷與表格包裝"
```

---

### Task 4: 前台版面（Header、Footer、路由、GA4）

**Files:**
- Create: `web/src/site/queries.ts`、`analytics.ts`、`useDocumentTitle.ts`、`SiteApp.tsx`、`Header.tsx`、`Footer.tsx`、`NotFound.tsx`、`PageView.tsx`（本 task 先放簡化版，Task 5 補完）
- Modify: `web/src/App.tsx`
- Test: `web/src/site/analytics.test.ts`、`Header.test.tsx`、`Footer.test.tsx`

**Interfaces:**
- Consumes: `api`、`ApiError`、`publicConfig`、型別（Task 2）
- Produces:
  ```ts
  // queries.ts
  export function useSite(): UseQueryResult<SiteData>
  export function useHome(): UseQueryResult<PageData>
  export function usePage(id: string): UseQueryResult<PageData>
  // analytics.ts
  export function initAnalytics(id: string): void
  export function trackPageView(path: string): void
  // useDocumentTitle.ts
  export function useDocumentTitle(title: string | undefined): void
  // 元件
  export default function Header(props: { title: string; subtitle: string; menu: Group[] })
  export default function Footer(props: { title: string; facebookUrl: string; counterScriptUrl: string })
  export default function PageView(props: { query: UseQueryResult<PageData>; greeting?: string; groupName?: string })
  //   greeting：首頁問候語（明體大字）；groupName：內頁標題上方的所在頁籤。兩者都沒有時不顯示頁面標題（首頁）
  ```

- [ ] **Step 1: 寫失敗的測試**

`web/src/site/analytics.test.ts`：

```ts
import { beforeEach, expect, it, vi } from 'vitest'

beforeEach(() => {
  vi.resetModules()
  document.head.innerHTML = ''
  delete window.dataLayer
  delete window.gtag
})

it('沒有 ID 時不載入', async () => {
  const { initAnalytics, trackPageView } = await import('./analytics')
  initAnalytics('')
  trackPageView('/')
  expect(document.querySelector('script')).toBeNull()
  expect(window.dataLayer).toBeUndefined()
})

it('載入 gtag 一次並送出 page_view', async () => {
  const { initAnalytics, trackPageView } = await import('./analytics')
  initAnalytics('G-TEST')
  initAnalytics('G-TEST')
  const scripts = document.querySelectorAll('script')
  expect(scripts).toHaveLength(1)
  expect(scripts[0].src).toBe('https://www.googletagmanager.com/gtag/js?id=G-TEST')

  trackPageView('/page/3')
  const events = window.dataLayer!.map((args) => Array.from(args as ArrayLike<unknown>))
  expect(events).toContainEqual(['config', 'G-TEST', { send_page_view: false }])
  expect(events).toContainEqual(['event', 'page_view', expect.objectContaining({ page_path: '/page/3' })])
})
```

`web/src/site/Header.test.tsx`：

```tsx
import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes, useLocation } from 'react-router'
import { expect, it } from 'vitest'
import Header from './Header'
import type { Group } from '@/shared/types'

const menu: Group[] = [
  { id: 1, name: '首頁', pages: [{ id: 1, name: '最新消息' }] },
  { id: 2, name: '主要行事活動', pages: [{ id: 3, name: '練成會' }, { id: 4, name: '誌友會' }] },
]

function Where() {
  return <p data-testid="where">{useLocation().pathname}</p>
}

function setup(path = '/') {
  const user = userEvent.setup()
  render(
    <MemoryRouter initialEntries={[path]}>
      <Header title="生長之家" subtitle="一句簡單的感謝" menu={menu} />
      <Routes>
        <Route path="*" element={<Where />} />
      </Routes>
    </MemoryRouter>,
  )
  return user
}

it('滑過群組顯示頁面，點擊後導覽並關閉', async () => {
  const user = setup()
  const nav = screen.getByRole('navigation', { name: '主選單' })
  const trigger = within(nav).getByRole('button', { name: '主要行事活動' })
  expect(trigger).toHaveAttribute('aria-expanded', 'false')

  await user.hover(trigger)
  expect(trigger).toHaveAttribute('aria-expanded', 'true')
  await user.click(within(nav).getByRole('link', { name: '誌友會' }))
  expect(screen.getByTestId('where')).toHaveTextContent('/page/4')
  expect(within(nav).queryByRole('link', { name: '誌友會' })).toBeNull()
})

it('標示目前頁面所在的頁籤', () => {
  setup('/page/4')
  const nav = screen.getByRole('navigation', { name: '主選單' })
  expect(within(nav).getByRole('button', { name: '主要行事活動' })).toHaveAttribute('data-active', 'true')
  expect(within(nav).getByRole('button', { name: '首頁' })).not.toHaveAttribute('data-active')
  expect(screen.getByText('一句簡單的感謝')).toBeInTheDocument()
})

it('鍵盤：Enter 開啟、方向鍵進入選單、Escape 關閉並回到按鈕', async () => {
  const user = setup()
  const nav = screen.getByRole('navigation', { name: '主選單' })
  const trigger = within(nav).getByRole('button', { name: '主要行事活動' })
  trigger.focus()
  await user.keyboard('{ArrowDown}')
  await expect.poll(() => document.activeElement?.textContent).toBe('練成會')
  await user.keyboard('{Escape}')
  expect(trigger).toHaveAttribute('aria-expanded', 'false')
  expect(trigger).toHaveFocus()

  await user.keyboard('{Enter}')
  expect(trigger).toHaveAttribute('aria-expanded', 'true')
})

it('手機抽屜：手風琴展開群組並導覽', async () => {
  const user = setup()
  await user.click(screen.getByRole('button', { name: '開啟選單' }))
  const drawer = screen.getByRole('dialog', { name: '選單' })
  expect(within(drawer).queryByRole('link', { name: '練成會' })).toBeNull()
  await user.click(within(drawer).getByRole('button', { name: '主要行事活動' }))
  await user.click(within(drawer).getByRole('link', { name: '練成會' }))
  expect(screen.getByTestId('where')).toHaveTextContent('/page/3')
  expect(screen.queryByRole('dialog')).toBeNull()
})
```

`web/src/site/Footer.test.tsx`：

```tsx
import { render, screen } from '@testing-library/react'
import { expect, it } from 'vitest'
import Footer from './Footer'

it('有設定時顯示 Facebook 連結與計數器', () => {
  render(<Footer title="生長之家" facebookUrl="https://www.facebook.com/x" counterScriptUrl={'https://c.example/s.php?a=1&b="x"'} />)
  const fb = screen.getByRole('link', { name: /Facebook/ })
  expect(fb).toHaveAttribute('href', 'https://www.facebook.com/x')
  expect(fb).toHaveAttribute('target', '_blank')
  expect(fb).toHaveAttribute('rel', 'noopener noreferrer')
  const frame = screen.getByTitle('瀏覽次數')
  expect(frame.getAttribute('srcdoc')).toContain('src="https://c.example/s.php?a=1&amp;b=&quot;x&quot;"')
  expect(frame).toHaveAttribute('sandbox', 'allow-scripts')
})

it('沒有設定時不顯示', () => {
  render(<Footer title="生長之家" facebookUrl="" counterScriptUrl="" />)
  expect(screen.queryByRole('link')).toBeNull()
  expect(screen.queryByTitle('瀏覽次數')).toBeNull()
  expect(screen.getByText(/Seicho-No-Ie/)).toBeInTheDocument()
})
```

Run: `pnpm --dir web test`
Expected: FAIL（找不到模組）

- [ ] **Step 2: 實作資料與工具**

`web/src/site/queries.ts`：

```ts
import { useQuery } from '@tanstack/react-query'
import { api, ApiError } from '@/shared/api'
import type { PageData, SiteData } from '@/shared/types'

export const useSite = () =>
  useQuery({ queryKey: ['site'], queryFn: ({ signal }) => api<SiteData>('/site', { signal }) })

export const useHome = () =>
  useQuery({ queryKey: ['home'], queryFn: ({ signal }) => api<PageData>('/home', { signal }) })

export const usePage = (id: string) =>
  useQuery({
    queryKey: ['page', id],
    queryFn: ({ signal }) => api<PageData>(`/pages/${id}`, { signal }),
    retry: (count, err) => !(err instanceof ApiError && err.status < 500) && count < 1,
  })
```

`web/src/site/analytics.ts`：

```ts
declare global {
  interface Window {
    dataLayer?: unknown[]
    gtag?: (...args: unknown[]) => void
  }
}

let loaded = false

export function initAnalytics(id: string) {
  if (!id || loaded) return
  loaded = true
  window.dataLayer = window.dataLayer ?? []
  // gtag.js 需要 push arguments 物件，不能改成陣列
  window.gtag = function gtag() {
    window.dataLayer!.push(arguments)
  }
  window.gtag('js', new Date())
  window.gtag('config', id, { send_page_view: false })
  const script = document.createElement('script')
  script.async = true
  script.src = `https://www.googletagmanager.com/gtag/js?id=${encodeURIComponent(id)}`
  document.head.append(script)
}

export function trackPageView(path: string) {
  window.gtag?.('event', 'page_view', {
    page_path: path,
    page_location: window.location.origin + path,
    page_title: document.title,
  })
}
```

`web/src/site/useDocumentTitle.ts`：

```ts
import { useEffect } from 'react'

export function useDocumentTitle(title: string | undefined) {
  useEffect(() => {
    if (title) document.title = title
  }, [title])
}
```

- [ ] **Step 3: 實作 Header**

`web/src/site/Header.tsx`：

```tsx
import { useEffect, useRef, useState } from 'react'
import { Link, useLocation } from 'react-router'
import type { Group } from '@/shared/types'

type Props = { title: string; subtitle: string; menu: Group[] }

export default function Header({ title, subtitle, menu }: Props) {
  const [open, setOpen] = useState<number | null>(null)
  const [drawer, setDrawer] = useState(false)
  const [compact, setCompact] = useState(false)
  const navRef = useRef<HTMLElement>(null)
  const { pathname } = useLocation()
  const activeGroupId = menu.find((g) => g.pages.some((p) => pathname === `/page/${p.id}`))?.id

  // 換頁時關閉所有選單
  const [lastPath, setLastPath] = useState(pathname)
  if (lastPath !== pathname) {
    setLastPath(pathname)
    setOpen(null)
    setDrawer(false)
  }

  useEffect(() => {
    const onScroll = () => setCompact(window.scrollY > 40)
    onScroll()
    window.addEventListener('scroll', onScroll, { passive: true })
    return () => window.removeEventListener('scroll', onScroll)
  }, [])

  // 觸控裝置沒有 mouseleave，點選單外面時關閉
  useEffect(() => {
    if (open === null) return
    const onPointer = (e: PointerEvent) => {
      if (!navRef.current?.contains(e.target as Node)) setOpen(null)
    }
    document.addEventListener('pointerdown', onPointer)
    return () => document.removeEventListener('pointerdown', onPointer)
  }, [open])

  const focusFirst = (id: number) =>
    requestAnimationFrame(() => document.querySelector<HTMLElement>(`#menu-${id} a`)?.focus())

  return (
    <header className={`sticky top-0 z-40 border-b bg-white transition-[border-color] ${compact ? 'border-slate-200' : 'border-transparent'}`}>
      <div className={`mx-auto flex max-w-6xl items-center justify-between gap-4 px-4 transition-[height] ${compact ? 'h-14' : 'h-20'}`}>
        <Link to="/" className="flex min-w-0 items-center gap-3">
          <img src="/logo.png" alt="" className={`w-auto transition-[height] ${compact ? 'h-9' : 'h-12'}`} />
          <span className="truncate font-serif text-2xl font-bold tracking-wider text-brand-dark">{title}</span>
          {subtitle && <span className="hidden font-serif text-sm text-slate-500 xl:inline">{subtitle}</span>}
        </Link>

        <nav ref={navRef} aria-label="主選單" className="hidden lg:block">
          <ul className="flex gap-1">
            {menu.map((g) => (
              <li key={g.id} className="relative" onMouseEnter={() => setOpen(g.id)} onMouseLeave={() => setOpen(null)}>
                <button
                  type="button"
                  data-active={g.id === activeGroupId || undefined}
                  aria-expanded={open === g.id}
                  aria-controls={`menu-${g.id}`}
                  onClick={() => setOpen(g.id)}
                  onKeyDown={(e) => {
                    if (e.key === 'ArrowDown') {
                      e.preventDefault()
                      setOpen(g.id)
                      focusFirst(g.id)
                    }
                  }}
                  className="border-b-2 border-transparent px-3 py-2 font-medium text-ink hover:text-brand aria-expanded:text-brand data-active:border-olive"
                >
                  {g.name}
                </button>
                {open === g.id && g.pages.length > 0 && (
                  <ul
                    id={`menu-${g.id}`}
                    className="absolute right-0 top-full min-w-52 rounded-md border border-slate-200 bg-white py-2 shadow-lg"
                    onKeyDown={(e) => {
                      const links = [...e.currentTarget.querySelectorAll<HTMLElement>('a')]
                      const i = links.indexOf(document.activeElement as HTMLElement)
                      if (e.key === 'Escape') {
                        setOpen(null)
                        ;(e.currentTarget.previousElementSibling as HTMLElement | null)?.focus()
                      } else if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
                        e.preventDefault()
                        const next = e.key === 'ArrowDown' ? i + 1 : i - 1
                        links[(next + links.length) % links.length]?.focus()
                      }
                    }}
                  >
                    {g.pages.map((p) => (
                      <li key={p.id}>
                        <Link
                          to={`/page/${p.id}`}
                          className="block whitespace-nowrap px-4 py-2 text-ink hover:bg-brand-soft hover:text-brand focus-visible:bg-brand-soft focus-visible:outline-none"
                        >
                          {p.name}
                        </Link>
                      </li>
                    ))}
                  </ul>
                )}
              </li>
            ))}
          </ul>
        </nav>

        <button
          type="button"
          aria-label="開啟選單"
          aria-expanded={drawer}
          onClick={() => setDrawer(true)}
          className="rounded-md p-2 text-brand-dark lg:hidden"
        >
          <svg viewBox="0 0 24 24" className="h-7 w-7" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
            <path d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
      </div>
      {drawer && <MobileDrawer menu={menu} onClose={() => setDrawer(false)} />}
    </header>
  )
}

function MobileDrawer({ menu, onClose }: { menu: Group[]; onClose: () => void }) {
  const [expanded, setExpanded] = useState<number | null>(null)
  const closeRef = useRef<HTMLButtonElement>(null)

  useEffect(() => {
    closeRef.current?.focus()
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && onClose()
    window.addEventListener('keydown', onKey)
    document.body.style.overflow = 'hidden'
    return () => {
      window.removeEventListener('keydown', onKey)
      document.body.style.overflow = ''
    }
  }, [onClose])

  return (
    <div className="fixed inset-0 z-50 lg:hidden" role="dialog" aria-modal="true" aria-label="選單">
      <div className="absolute inset-0 bg-black/40" onClick={onClose} />
      <nav className="absolute right-0 top-0 flex h-full w-72 max-w-[85vw] flex-col bg-white shadow-xl">
        <div className="flex justify-end p-2">
          <button ref={closeRef} type="button" aria-label="關閉選單" onClick={onClose} className="rounded-md p-2 text-2xl leading-none">
            ×
          </button>
        </div>
        <ul className="overflow-y-auto px-4 pb-8">
          {menu.map((g) => (
            <li key={g.id} className="border-b border-slate-100">
              <button
                type="button"
                aria-expanded={expanded === g.id}
                onClick={() => setExpanded(expanded === g.id ? null : g.id)}
                className="flex w-full items-center justify-between py-3 text-left font-medium text-slate-800"
              >
                {g.name}
                <span aria-hidden="true" className="text-brand">{expanded === g.id ? '−' : '+'}</span>
              </button>
              {expanded === g.id && (
                <ul className="pb-2">
                  {g.pages.map((p) => (
                    <li key={p.id}>
                      <Link to={`/page/${p.id}`} onClick={onClose} className="block py-2 pl-4 text-slate-600">
                        {p.name}
                      </Link>
                    </li>
                  ))}
                </ul>
              )}
            </li>
          ))}
        </ul>
      </nav>
    </div>
  )
}
```

> 「換頁時關閉」使用 render 期間比較前一個 pathname 的寫法（React 官方建議，取代在 effect 中 setState）。

- [ ] **Step 4: 實作 Footer、NotFound、PageView（簡化版）、SiteApp**

`web/src/site/Footer.tsx`：

```tsx
type Props = { title: string; facebookUrl: string; counterScriptUrl: string }

const escapeAttr = (s: string) =>
  s.replaceAll('&', '&amp;').replaceAll('"', '&quot;').replaceAll('<', '&lt;').replaceAll('>', '&gt;')

export default function Footer({ title, facebookUrl, counterScriptUrl }: Props) {
  return (
    <footer className="mt-20 bg-brand-dark text-slate-200">
      <div className="mx-auto flex max-w-6xl flex-col gap-6 px-4 py-12 text-sm sm:flex-row sm:items-end sm:justify-between">
        <div className="space-y-2">
          {title && <p className="font-serif text-xl font-bold tracking-wider text-white">{title}</p>}
          <p>© {new Date().getFullYear()} Seicho-No-Ie R.O.C Missionary Headquarters All Rights Reserved</p>
        </div>
        <div className="flex flex-col items-start gap-3 sm:items-end">
          {facebookUrl && (
            <a
              href={facebookUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 underline-offset-4 hover:underline"
            >
              <svg viewBox="0 0 24 24" className="h-4 w-4" fill="currentColor" aria-hidden="true">
                <path d="M14 8h3V4h-3c-2.8 0-4 1.7-4 4.3V10H7v4h3v8h4v-8h3l1-4h-4V8.5c0-.3.2-.5.5-.5Z" />
              </svg>
              Facebook
            </a>
          )}
          {counterScriptUrl && (
            // 舊式計數器以 document.write 輸出，放在沙箱 iframe 內執行
            <iframe
              title="瀏覽次數"
              sandbox="allow-scripts"
              className="h-6 w-56 border-0"
              srcDoc={`<!doctype html><meta charset="utf-8"><body style="margin:0;font:13px sans-serif;color:#e2e8f0;white-space:nowrap">本站瀏覽次數：<script src="${escapeAttr(counterScriptUrl)}"></script></body>`}
            />
          )}
        </div>
      </div>
    </footer>
  )
}
```

`web/src/site/NotFound.tsx`：

```tsx
import { Link } from 'react-router'
import { useDocumentTitle } from './useDocumentTitle'

export default function NotFound() {
  useDocumentTitle('找不到頁面')
  return (
    <div className="mx-auto max-w-[42rem] px-4 py-24">
      <h1 className="font-serif text-3xl font-bold text-brand-dark">找不到這個頁面</h1>
      <p className="mt-4 text-lg text-slate-600">這個頁面可能已經移除，或網址有誤。可以從上方選單找到其他頁面。</p>
      <Link to="/" className="mt-8 inline-block rounded-md bg-brand px-5 py-2.5 text-white hover:bg-brand-dark">
        回到首頁
      </Link>
    </div>
  )
}
```

`web/src/site/PageView.tsx`（簡化版，Task 5 取代）：

```tsx
import type { UseQueryResult } from '@tanstack/react-query'
import type { PageData } from '@/shared/types'

type Props = { query: UseQueryResult<PageData>; greeting?: string; groupName?: string }

export default function PageView({ query, groupName }: Props) {
  if (!query.data) return null
  const { page } = query.data
  return <article className="mx-auto max-w-[42rem] px-4 py-8">{groupName && page && <h1>{page.name}</h1>}</article>
}
```

`web/src/site/SiteApp.tsx`：

```tsx
import { useEffect } from 'react'
import { Route, Routes, useLocation, useParams } from 'react-router'
import { ApiError } from '@/shared/api'
import { publicConfig } from '@/shared/config'
import type { SiteData } from '@/shared/types'
import { initAnalytics, trackPageView } from './analytics'
import Footer from './Footer'
import Header from './Header'
import NotFound from './NotFound'
import PageView from './PageView'
import { useHome, usePage, useSite } from './queries'
import { useDocumentTitle } from './useDocumentTitle'

export default function SiteApp() {
  const { data: site } = useSite()
  const { pathname } = useLocation()

  useEffect(() => initAnalytics(publicConfig.gaMeasurementId), [])
  useEffect(() => {
    window.scrollTo(0, 0)
    trackPageView(pathname)
  }, [pathname])

  return (
    <div className="flex min-h-screen flex-col">
      <Header title={site?.web_title ?? ''} subtitle={site?.web_sub_title ?? ''} menu={site?.menu ?? []} />
      <main className="flex-1">
        <Routes>
          <Route index element={<HomeRoute site={site} />} />
          <Route path="page/:id" element={<PageRoute site={site} />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </main>
      <Footer title={site?.web_title ?? ''} facebookUrl={site?.facebook_url ?? ''} counterScriptUrl={publicConfig.counterScriptUrl} />
    </div>
  )
}

function HomeRoute({ site }: { site?: SiteData }) {
  const query = useHome()
  useDocumentTitle(site && (site.web_sub_title ? `${site.web_title}｜${site.web_sub_title}` : site.web_title))
  return <PageView query={query} greeting={site?.web_sub_title} />
}

function PageRoute({ site }: { site?: SiteData }) {
  const { id = '' } = useParams()
  const query = usePage(id)
  const page = query.data?.page
  const groupName = site?.menu.find((g) => g.id === page?.group_id)?.name
  useDocumentTitle(page && site ? `${page.name}｜${site.web_title}` : undefined)
  if (query.error instanceof ApiError && (query.error.status === 404 || query.error.status === 400)) {
    return <NotFound />
  }
  return <PageView query={query} groupName={groupName ?? ''} />
}
```

`web/src/App.tsx` 中 `<Route path="*" element={<p className="p-8">前台建置中</p>} />` 換成 `<Route path="*" element={<SiteApp />} />`，並 `import SiteApp from '@/site/SiteApp'`。

- [ ] **Step 5: 確認測試通過**

Run: `pnpm --dir web test && pnpm --dir web typecheck`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add web/src
git commit -m "[feat] 前台版面：選單、手機抽屜、footer、路由與 GA4"
```

---

### Task 5: 輪播圖、跑馬燈與舊內容呈現

**Files:**
- Create: `web/src/site/Banner.tsx`、`Marquee.tsx`、`LegacyContent.tsx`
- Modify: `web/src/site/PageView.tsx`（整個取代）、`web/src/index.css`
- Test: `web/src/site/Marquee.test.tsx`、`LegacyContent.test.tsx`、`PageView.test.tsx`

**Interfaces:**
- Consumes: `internalPath`、`wrapTables`（Task 3）
- Produces:
  ```ts
  export default function Banner(props: { items: Carousel[] })
  export default function Marquee(props: { items: Marquee[] })
  export default function LegacyContent(props: { html: string })
  export function usePrefersReducedMotion(): boolean      // 由 Marquee.tsx 匯出，Banner 共用
  ```

- [ ] **Step 1: 寫失敗的測試**

`web/src/site/Marquee.test.tsx`：

```tsx
import { act, render, screen } from '@testing-library/react'
import { afterEach, expect, it, vi } from 'vitest'
import Marquee from './Marquee'

const items = [
  { id: 1, text: '合掌感謝！', color: '#1EFF00' },
  { id: 2, text: '改版上線', color: '#FFFFFF' },
]

function mockReducedMotion(matches: boolean) {
  vi.stubGlobal('matchMedia', (query: string) => ({
    matches,
    media: query,
    addEventListener: () => {},
    removeEventListener: () => {},
  }))
}

afterEach(() => {
  vi.unstubAllGlobals()
  vi.useRealTimers()
})

it('捲動模式：每則使用自己的顏色', () => {
  mockReducedMotion(false)
  render(<Marquee items={items} />)
  const region = screen.getByRole('region', { name: '最新公告' })
  const first = region.querySelector('span')!
  expect(first).toHaveTextContent('合掌感謝！')
  expect(first).toHaveStyle({ color: '#1EFF00' })
  expect(region.querySelector('.marquee-track')).not.toBeNull()
})

it('減少動態時改為靜態輪替', () => {
  vi.useFakeTimers()
  mockReducedMotion(true)
  render(<Marquee items={items} />)
  expect(screen.getByText('合掌感謝！')).toBeInTheDocument()
  expect(document.querySelector('.marquee-track')).toBeNull()
  act(() => vi.advanceTimersByTime(5000))
  expect(screen.getByText('改版上線')).toBeInTheDocument()
  expect(screen.queryByText('合掌感謝！')).toBeNull()
})

it('沒有資料時不顯示', () => {
  mockReducedMotion(false)
  const { container } = render(<Marquee items={[]} />)
  expect(container).toBeEmptyDOMElement()
})
```

`web/src/site/LegacyContent.test.tsx`：

```tsx
import { fireEvent, render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes, useLocation } from 'react-router'
import { expect, it } from 'vitest'
import LegacyContent from './LegacyContent'

function Where() {
  return <p data-testid="where">{useLocation().pathname}</p>
}

function setup(html: string) {
  render(
    <MemoryRouter initialEntries={['/page/1']}>
      <LegacyContent html={html} />
      <Routes>
        <Route path="*" element={<Where />} />
      </Routes>
    </MemoryRouter>,
  )
}

it('包裝表格', () => {
  setup('<table><tr><td>1</td></tr></table>')
  expect(document.querySelector('.legacy-content .table-scroll > table')).not.toBeNull()
})

it('站內舊連結改用前端路由', () => {
  setup(`<p><a href="http://${location.hostname}/#/3"><strong>練成會</strong></a></p>`)
  const notPrevented = fireEvent.click(screen.getByText('練成會'))
  expect(notPrevented).toBe(false)
  expect(screen.getByTestId('where')).toHaveTextContent('/page/3')
})

it.each([
  ['外部連結', '<a href="https://www.facebook.com/x">連結</a>', {}],
  ['開新分頁', '<a href="/#/3" target="_blank">連結</a>', {}],
  ['按住 Ctrl', '<a href="/#/3">連結</a>', { ctrlKey: true }],
  ['同頁錨點', '<a href="#top">連結</a>', {}],
])('%s維持瀏覽器預設行為', (_, html, init) => {
  setup(html)
  const notPrevented = fireEvent.click(screen.getByText('連結'), init)
  expect(notPrevented).toBe(true)
  expect(screen.getByTestId('where')).toHaveTextContent('/page/1')
})
```

`web/src/site/PageView.test.tsx`：

```tsx
import type { UseQueryResult } from '@tanstack/react-query'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router'
import { expect, it, vi } from 'vitest'
import { ApiError } from '@/shared/api'
import type { PageData } from '@/shared/types'
import PageView from './PageView'

vi.mock('./Banner', () => ({ default: () => <div data-testid="banner" /> }))

const q = (over: Partial<UseQueryResult<PageData>>) =>
  ({ isPending: false, isError: false, data: undefined, error: null, refetch: vi.fn(), ...over }) as UseQueryResult<PageData>

const renderView = (query: UseQueryResult<PageData>, props: { greeting?: string; groupName?: string } = { groupName: '主要行事活動' }) =>
  render(
    <MemoryRouter>
      <PageView query={query} {...props} />
    </MemoryRouter>,
  )

it('載入中顯示骨架', () => {
  renderView(q({ isPending: true }))
  expect(screen.getByLabelText('載入中')).toBeInTheDocument()
})

it('錯誤時顯示訊息並可重試', async () => {
  const query = q({ isError: true, error: new ApiError(0, 'network', '無法連線到伺服器') })
  renderView(query)
  expect(screen.getByRole('alert')).toHaveTextContent('無法連線到伺服器')
  await userEvent.click(screen.getByRole('button', { name: '重新載入' }))
  expect(query.refetch).toHaveBeenCalled()
})

it('內頁顯示所在頁籤、標題與內容', () => {
  renderView(q({ data: { page: { id: 3, group_id: 2, name: '練成會', html: '<p>日程</p>' }, carousels: [], marquees: [] } }))
  expect(screen.getByRole('heading', { name: '練成會' })).toBeInTheDocument()
  expect(screen.getByText('主要行事活動')).toBeInTheDocument()
  expect(screen.getByText('日程')).toBeInTheDocument()
})

it('首頁顯示問候語而不顯示頁面標題；沒有內容時顯示提示', () => {
  renderView(q({ data: { page: null, carousels: [], marquees: [] } }), { greeting: '一句簡單的感謝' })
  expect(screen.getByText('一句簡單的感謝')).toBeInTheDocument()
  expect(screen.queryByRole('heading')).toBeNull()
  expect(screen.getByText('這裡還沒有內容。')).toBeInTheDocument()
})
```

Run: `pnpm --dir web test`
Expected: FAIL

- [ ] **Step 2: 實作 Marquee**

`web/src/site/Marquee.tsx`：

```tsx
import { useEffect, useState, useSyncExternalStore } from 'react'
import type { Marquee as MarqueeItem } from '@/shared/types'

const reducedQuery = '(prefers-reduced-motion: reduce)'

export function usePrefersReducedMotion() {
  return useSyncExternalStore(
    (onChange) => {
      const mq = window.matchMedia?.(reducedQuery)
      mq?.addEventListener('change', onChange)
      return () => mq?.removeEventListener('change', onChange)
    },
    () => window.matchMedia?.(reducedQuery).matches ?? false,
    () => false,
  )
}

export default function Marquee({ items }: { items: MarqueeItem[] }) {
  const reduced = usePrefersReducedMotion()
  if (items.length === 0) return null
  return (
    <div role="region" aria-label="最新公告" className="bg-brand-dark text-sm font-medium sm:text-base">
      {reduced ? <Rotator items={items} /> : <Scroller items={items} />}
    </div>
  )
}

function Scroller({ items }: { items: MarqueeItem[] }) {
  const row = (copy: string) =>
    items.map((m) => (
      <span key={`${copy}-${m.id}`} style={{ color: m.color }} className="px-10">
        {m.text}
      </span>
    ))
  // 內容重複兩份、動畫位移 -50%，形成無縫循環
  return (
    <div className="group overflow-hidden py-2">
      <div
        className="marquee-track flex w-max whitespace-nowrap group-hover:[animation-play-state:paused]"
        style={{ animationDuration: `${Math.max(20, items.length * 12)}s` }}
      >
        {row('a')}
        <span aria-hidden="true" className="flex">
          {row('b')}
        </span>
      </div>
    </div>
  )
}

function Rotator({ items }: { items: MarqueeItem[] }) {
  const [index, setIndex] = useState(0)
  useEffect(() => {
    if (items.length < 2) return
    const timer = setInterval(() => setIndex((i) => (i + 1) % items.length), 5000)
    return () => clearInterval(timer)
  }, [items.length])
  const m = items[index % items.length]
  return (
    <p aria-live="polite" className="px-4 py-2 text-center" style={{ color: m.color }}>
      {m.text}
    </p>
  )
}
```

- [ ] **Step 3: 實作 Banner 與 LegacyContent**

`web/src/site/Banner.tsx`：

```tsx
import { useEffect, useState, type ReactNode } from 'react'
import { Link } from 'react-router'
import useEmblaCarousel from 'embla-carousel-react'
import Autoplay from 'embla-carousel-autoplay'
import type { Carousel } from '@/shared/types'
import { internalPath } from './legacyLinks'
import { usePrefersReducedMotion } from './Marquee'

export default function Banner({ items }: { items: Carousel[] }) {
  const reduced = usePrefersReducedMotion()
  const [viewportRef, embla] = useEmblaCarousel({ loop: true }, [
    Autoplay({ delay: 5000, active: !reduced, stopOnInteraction: false, stopOnMouseEnter: true }),
  ])
  const [selected, setSelected] = useState(0)

  useEffect(() => {
    if (!embla) return
    const onSelect = () => setSelected(embla.selectedScrollSnap())
    onSelect()
    embla.on('select', onSelect).on('reInit', onSelect)
    return () => {
      embla.off('select', onSelect).off('reInit', onSelect)
    }
  }, [embla])

  if (items.length === 0) return null
  return (
    <section aria-roledescription="carousel" aria-label="輪播圖" className="relative mx-auto max-w-6xl sm:px-4 sm:pt-6">
      <div ref={viewportRef} className="overflow-hidden sm:rounded-md">
        <div className="flex">
          {items.map((c, i) => (
            <div key={c.id} aria-roledescription="slide" aria-label={`第 ${i + 1} 張，共 ${items.length} 張`} className="min-w-0 flex-[0_0_100%]">
              <SlideLink url={c.url}>
                <img
                  src={c.image}
                  alt=""
                  loading={i === 0 ? 'eager' : 'lazy'}
                  className="aspect-[16/7] w-full bg-brand-soft object-cover"
                />
              </SlideLink>
            </div>
          ))}
        </div>
      </div>
      {items.length > 1 && (
        <>
          <ArrowButton label="上一張" side="left" onClick={() => embla?.scrollPrev()} />
          <ArrowButton label="下一張" side="right" onClick={() => embla?.scrollNext()} />
          <div className="absolute inset-x-0 bottom-3 flex justify-center gap-2">
            {items.map((c, i) => (
              <button
                key={c.id}
                type="button"
                aria-label={`第 ${i + 1} 張`}
                aria-current={i === selected}
                onClick={() => embla?.scrollTo(i)}
                className={`h-2.5 rounded-full transition-all ${i === selected ? 'w-6 bg-white' : 'w-2.5 bg-white/60'}`}
              />
            ))}
          </div>
        </>
      )}
    </section>
  )
}

function ArrowButton({ label, side, onClick }: { label: string; side: 'left' | 'right'; onClick: () => void }) {
  return (
    <button
      type="button"
      aria-label={label}
      onClick={onClick}
      className={`absolute top-1/2 hidden h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-black/30 text-2xl text-white hover:bg-black/50 sm:flex ${side === 'left' ? 'left-6' : 'right-6'}`}
    >
      {side === 'left' ? '‹' : '›'}
    </button>
  )
}

function SlideLink({ url, children }: { url: string; children: ReactNode }) {
  if (!url) return children
  const to = internalPath(url, window.location)
  if (to) return <Link to={to}>{children}</Link>
  return (
    <a href={url} target="_blank" rel="noopener noreferrer">
      {children}
    </a>
  )
}
```

`web/src/site/LegacyContent.tsx`：

```tsx
import { useLayoutEffect, useRef, type MouseEvent } from 'react'
import { useNavigate } from 'react-router'
import { internalPath } from './legacyLinks'
import { wrapTables } from './wrapTables'

// html 已由後端過濾，這裡原樣輸出
export default function LegacyContent({ html }: { html: string }) {
  const ref = useRef<HTMLDivElement>(null)
  const navigate = useNavigate()

  useLayoutEffect(() => {
    if (ref.current) wrapTables(ref.current)
  }, [html])

  const onClick = (e: MouseEvent<HTMLDivElement>) => {
    if (e.defaultPrevented || e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) return
    const a = (e.target as Element).closest('a')
    if (!a || (a.target && a.target !== '_self') || a.hasAttribute('download')) return
    const to = internalPath(a.getAttribute('href') ?? '', window.location)
    if (!to) return
    e.preventDefault()
    navigate(to)
  }

  return <div ref={ref} className="legacy-content" onClick={onClick} dangerouslySetInnerHTML={{ __html: html }} />
}
```

- [ ] **Step 4: 完整 PageView**

`web/src/site/PageView.tsx`（整個取代）：

```tsx
import type { UseQueryResult } from '@tanstack/react-query'
import type { PageData } from '@/shared/types'
import Banner from './Banner'
import LegacyContent from './LegacyContent'
import Marquee from './Marquee'

type Props = { query: UseQueryResult<PageData>; greeting?: string; groupName?: string }

export default function PageView({ query, greeting, groupName }: Props) {
  if (query.isPending) return <Skeleton />
  if (query.isError) {
    return (
      <div role="alert" className="mx-auto max-w-[42rem] px-4 py-24">
        <p className="text-lg text-ink">{query.error.message}</p>
        <button
          type="button"
          onClick={() => query.refetch()}
          className="mt-6 rounded-md bg-brand px-5 py-2.5 text-white hover:bg-brand-dark"
        >
          重新載入
        </button>
      </div>
    )
  }

  const { page, carousels, marquees } = query.data
  const isInner = groupName !== undefined
  return (
    <>
      <Banner items={carousels} />
      <Marquee items={marquees} />
      <article className="mx-auto max-w-[42rem] px-4 pb-8 pt-10 sm:pt-16">
        {greeting && (
          // 首頁唯一的強調：以明體大字排出網站副標題
          <p className="mb-12 font-serif text-4xl font-bold leading-snug tracking-wide text-brand-dark sm:text-6xl">
            {greeting}
          </p>
        )}
        {isInner && page && (
          <header className="mb-8">
            {groupName && <p className="text-sm text-slate-500">{groupName}</p>}
            <h1 className="mt-1 font-serif text-3xl font-bold leading-snug text-brand-dark sm:text-4xl">{page.name}</h1>
          </header>
        )}
        {page ? <LegacyContent html={page.html} /> : <p className="text-slate-500">這裡還沒有內容。</p>}
      </article>
    </>
  )
}

function Skeleton() {
  return (
    <div aria-label="載入中" aria-busy="true" className="animate-pulse">
      <div className="mx-auto aspect-[16/7] max-w-6xl bg-brand-soft sm:mt-6 sm:rounded-md" />
      <div className="mx-auto max-w-[42rem] space-y-4 px-4 py-12">
        <div className="h-8 w-1/3 rounded bg-slate-200" />
        <div className="h-4 rounded bg-slate-200" />
        <div className="h-4 w-5/6 rounded bg-slate-200" />
        <div className="h-4 w-2/3 rounded bg-slate-200" />
      </div>
    </div>
  )
}
```

- [ ] **Step 5: 內容樣式**

在 `web/src/index.css` 最後加入：

```css
@keyframes marquee {
  from { transform: translateX(0); }
  to { transform: translateX(-50%); }
}

.marquee-track {
  animation: marquee linear infinite;
}

/* 舊站內容：Tailwind preflight 會清掉預設樣式，這裡補回常用排版 */
.legacy-content {
  font-size: 1.0625rem;
  line-height: 1.85;
  overflow-wrap: anywhere;
}
@media (min-width: 640px) {
  .legacy-content { font-size: 1.125rem; }
}
.legacy-content p { margin: 0 0 1em; }
.legacy-content :is(h1, h2, h3, h4) { font-family: var(--font-serif); font-weight: 700; line-height: 1.4; margin: 1.6em 0 0.6em; }
.legacy-content h1 { font-size: 1.75em; }
.legacy-content h2 { font-size: 1.5em; }
.legacy-content h3 { font-size: 1.25em; }
.legacy-content h4 { font-size: 1.1em; }
.legacy-content a { color: var(--color-brand); text-decoration: underline; text-underline-offset: 0.2em; }
.legacy-content ul { list-style: disc; padding-left: 1.5em; margin: 0 0 1em; }
.legacy-content ol { list-style: decimal; padding-left: 1.5em; margin: 0 0 1em; }
.legacy-content blockquote { border-left: 4px solid var(--color-brand-soft); padding-left: 1em; color: #475569; }
.legacy-content hr { margin: 2em 0; border-color: #e2e8f0; }
.legacy-content img { display: inline-block; max-width: 100%; height: auto; }
/* 舊內容多寫死 width/height，以 !important 蓋過行內樣式，改為 16:9 自適應 */
.legacy-content iframe {
  display: block;
  width: 100% !important;
  max-width: 100%;
  height: auto !important;
  aspect-ratio: 16 / 9;
  border: 0;
}
.legacy-content .table-scroll { overflow-x: auto; margin: 1em 0; }
.legacy-content table { border-collapse: collapse; }
.legacy-content td,
.legacy-content th { padding: 0.4em 0.6em; vertical-align: top; }
```

- [ ] **Step 6: 確認測試通過並目視檢查**

Run: `pnpm --dir web test && pnpm --dir web typecheck`
Expected: PASS

目視檢查：

```bash
make up && make dev
```

開啟 `http://localhost:5173/`：輪播自動播放、跑馬燈顏色正確且 hover 暫停、首頁表格與 YouTube 影片寬度正常；把視窗縮到 375px 寬沒有橫向捲軸、漢堡選單可用；點內文「練成會」不會整頁重新載入。

- [ ] **Step 7: Commit**

```bash
git add web/src
git commit -m "[feat] 前台輪播圖、跑馬燈與舊內容呈現"
```

---

### Task 6: 前台 E2E（Playwright）

**Files:**
- Create: `web/playwright.config.ts`、`web/e2e/site.spec.ts`

**Interfaces:**
- Consumes: `make e2e`（Task 1）、`seed.sql` 資料（後端 Task 7）
- Produces: `web/playwright.config.ts` 的 `webServer` 設定；後台 plan 在 `web/e2e/` 加入 `admin.spec.ts` 共用此設定。

- [ ] **Step 1: Playwright 設定**

`web/playwright.config.ts`：

```ts
import { defineConfig, devices } from '@playwright/test'

// 需先 make up；make e2e 會建置前端、重設資料庫並放入範例圖片
export default defineConfig({
  testDir: './e2e',
  workers: 1,
  use: { baseURL: 'http://localhost:8080', trace: 'retain-on-failure' },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
  webServer: {
    command: 'go run ./backend/cmd/sniweb serve',
    cwd: '..',
    url: 'http://localhost:8080/readyz',
    reuseExistingServer: false,
    timeout: 120_000,
    env: {
      PORT: '8080',
      DATABASE_URL: 'root:sniweb@tcp(127.0.0.1:3306)/sniweb',
      PUBLIC_BASE_URL: 'http://localhost:8080',
      COOKIE_SECURE: 'false',
      STORAGE_DRIVER: 'disk',
      STORAGE_DISK_DIR: './tmp/picture',
    },
  },
})
```

- [ ] **Step 2: 寫 E2E 測試**

`web/e2e/site.spec.ts`：

```ts
import { expect, test } from '@playwright/test'

test('首頁顯示輪播、跑馬燈與內容', async ({ page }) => {
  await page.goto('/')
  await expect(page).toHaveTitle('生長之家｜一句簡單的感謝')
  await expect(page.getByRole('region', { name: '輪播圖' }).locator('img').first()).toBeVisible()
  await expect(page.getByRole('region', { name: '最新公告' }).getByText('合掌感謝！').first()).toBeVisible()
  await expect(page.getByText('●官方網站重新改版，歡迎瀏覽！')).toBeVisible()
  await expect(page.locator('.legacy-content .table-scroll table')).toHaveCount(1)
  await expect(page.locator('.legacy-content iframe')).toHaveAttribute('src', /youtube\.com\/embed/)
})

test('伺服器注入 SEO meta', async ({ request }) => {
  const html = await (await request.get('/page/3')).text()
  expect(html).toContain('<title>練成會｜生長之家</title>')
  expect(html).toContain('property="og:image" content="http://localhost:8080/php/picture/sample-1.jpg"')
})

test('舊 hash 網址轉到新網址', async ({ page }) => {
  await page.goto('/#/3')
  await expect(page).toHaveURL('/page/3')
  await expect(page.getByRole('heading', { name: '練成會' })).toBeVisible()
  await page.goto('/#/')
  await expect(page).toHaveURL('/')
})

test('內文舊連結以前端路由切換，不重新載入', async ({ page }) => {
  await page.goto('/')
  await page.evaluate(() => Object.assign(window, { __marker: 1 }))
  await page.locator('.legacy-content').getByRole('link', { name: '練成會' }).click()
  await expect(page).toHaveURL('/page/3')
  expect(await page.evaluate(() => (window as unknown as { __marker?: number }).__marker)).toBe(1)
})

test('桌機下拉選單', async ({ page }) => {
  await page.goto('/')
  const nav = page.getByRole('navigation', { name: '主選單' })
  await nav.getByRole('button', { name: '主要行事活動' }).hover()
  await nav.getByRole('link', { name: '誌友會' }).click()
  await expect(page).toHaveURL('/page/4')
  await expect(page).toHaveTitle('誌友會｜生長之家')
})

test('不存在的頁面回 404', async ({ page }) => {
  const res = await page.goto('/page/9999')
  expect(res?.status()).toBe(404)
  await expect(page.getByRole('heading', { name: '找不到這個頁面' })).toBeVisible()
})

test.describe('手機', () => {
  test.use({ viewport: { width: 375, height: 800 } })

  test('沒有橫向捲動，抽屜選單可用', async ({ page }) => {
    for (const path of ['/', '/page/3']) {
      await page.goto(path)
      await expect(page.locator('.legacy-content')).toBeVisible()
      const overflow = await page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth)
      expect(overflow, path).toBeLessThanOrEqual(0)
    }

    await page.getByRole('button', { name: '開啟選單' }).click()
    const drawer = page.getByRole('dialog', { name: '選單' })
    await drawer.getByRole('button', { name: '關於我們' }).click()
    await drawer.getByRole('link', { name: '生長之家簡介' }).click()
    await expect(page).toHaveURL('/page/2')
    await expect(drawer).toBeHidden()
  })
})
```

- [ ] **Step 3: 執行**

Run: `make up && make e2e`
Expected: 全部 PASS。失敗時用 `pnpm --dir web exec playwright show-trace web/test-results/<測試>/trace.zip` 查看。

- [ ] **Step 4: Commit**

```bash
git add web/playwright.config.ts web/e2e
git commit -m "[test] 前台 E2E：首頁、舊網址、選單、404 與手機版"
```

---

### Task 7: README 前端與部署說明

**Files:**
- Modify: `README.md`

- [ ] **Step 1: 在「本機開發」段落後加入**

````markdown
### 前後端一起開發

```bash
make up
cp .env.example .env
make dev          # Go :8080 + Vite :5173（/api、/php/picture 代理到 Go）
```

開啟 http://localhost:5173 。前端指令：

| 指令 | 說明 |
|---|---|
| `pnpm --dir web test` | Vitest |
| `pnpm --dir web typecheck` | TypeScript 型別檢查 |
| `make web` | 建置前端並放入 Go embed 目錄 |
| `make e2e` | 重設本機資料庫後執行 Playwright |

## 部署（Zeabur）

1. 建立 MySQL 服務，匯入舊資料並執行上方「資料庫」段落的 SQL。
2. 以 repo 根目錄的 `Dockerfile` 建立 app 服務，設定環境變數：
   `DATABASE_URL`、`PUBLIC_BASE_URL`（例如 `https://www.seicho-no-ie.org.tw`），
   以及 `STORAGE_DRIVER`（`disk` 或 `s3`）與對應變數、`GA_MEASUREMENT_ID`、`COUNTER_SCRIPT_URL`。
3. 使用 `disk` 時掛載 Volume 到 `/data`：首次掛載會清空目錄，掛載後再把舊的 `php/picture/*` 放進 `/data/picture/`；重新部署時會短暫停機。
4. 在 app 服務的終端機執行 `/sniweb user create --account <帳號> --name <暱稱>` 建立管理者。
5. 服務不會自動 migrate，schema 變更需手動執行。

`COUNTER_SCRIPT_URL` 若只有 `http://` 版本，在 https 網站上會被瀏覽器封鎖（混合內容），需改用供應商的 https 網址。
````

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "[docs] README：前端開發與 Zeabur 部署"
```

---

## 自我檢查紀錄

- spec 第 7 節（前台）每一項都有對應：路由與舊 hash（Task 3）、內文連結攔截（Task 3、5）、Header 下拉／抽屜／sticky（Task 4）、Banner（Task 5）、跑馬燈與 reduced motion（Task 5）、`.legacy-content` 規則（Task 5）、Footer 與計數器（Task 4）、GA4（Task 4）、主色／字型／骨架／錯誤提示（Task 1、5）、公開設定（Task 2）。
- spec 第 11 節 Dockerfile 在 Task 1；第 12 節 `make dev` 在 Task 1；第 13 節前台 E2E 在 Task 6（後台 E2E 在第三份 plan）。
- Banner 沒有單元測試（Embla 在 jsdom 需大量 mock），由 E2E 覆蓋。
- 設計方向檢查：已移除漸層底、卡片陰影、膠囊按鈕與「大數字 404」；強調只放在首頁問候語。輪播依 spec 使用 `object-cover`，若輪播圖多為含文字的活動海報，實作後需目視確認是否裁切到文字（必要時改為 `object-contain` 配夜藍底，需先與需求方確認）。
