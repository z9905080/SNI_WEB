# 生長之家網站重寫設計（React + Go）

- 日期：2026-09-16
- 分支：`rewrite/react-go`
- 參考：現行線上網站 https://www.seicho-no-ie.org.tw/#/

## 1. 目標與範圍

以 React（前端）與 Go（後端）重寫現有「生長之家前後台管理系統」。**功能與舊站對等，介面與架構全部重新設計**；前台保留「下拉選單 → banner 輪播 → 跑馬燈 → 內容 → footer」的結構，外觀現代化並修正手機版。

### 保留的功能

前台：
- 首頁：輪播圖、跑馬燈、首頁內容（`page_group_id = 1` 的第一個頁面）
- 內頁：依選單群組顯示頁面內容（內頁同樣顯示輪播圖與跑馬燈）
- 網站標題／副標題

後台：
- 登入／登出
- 選單群組（頁籤）CRUD 與拖曳排序
- 群組內頁面 CRUD、排序與 HTML 內容編輯
- 圖庫：上傳、列表、刪除
- 輪播圖：CRUD（從圖庫選圖 + 連結網址）
- 跑馬燈：CRUD（文字 + 顏色）

### 新增（小幅）

- 後台「網站設定」頁：編輯 `web_title`、`web_sub_title`、`facebook_url`（舊版只能改資料庫）
- Footer 固定顯示 Facebook 連結
- 帳號管理 CLI：`sniweb user create`、`sniweb user passwd`

### 不做

- 舊資料匯入（由專案擁有者另行撰寫 SQL migrate）
- 多語系、前台會員、頁面草稿／版本紀錄、輪播圖排序
- 舊後台中未接後端的模板頁面：`Register`、`FundList`、`InfoShow`、`ArticleList`

## 2. 資料庫

MySQL，沿用現有 7 張表，**結構不變**，唯一例外：

```sql
ALTER TABLE `user` MODIFY `pwd` varchar(255) NOT NULL COMMENT '密碼（bcrypt）';
```

| 表 | 引擎 | 用途 |
|---|---|---|
| `carousel` (id, image, url) | MyISAM | 輪播圖 |
| `marquee` (id, text, color) | MyISAM | 跑馬燈 |
| `page_group` (id, group_name, page_sort) | InnoDB | 選單群組；`page_sort` 為頁面 id 的 JSON 陣列 |
| `page_content` (id, page_group_id, page_name, html_context) | InnoDB | 頁面 |
| `user` (id, account, pwd, identity, user_name) | MyISAM | 管理者 |
| `user_token` (user_id PK, token, expire_time) | MyISAM | 登入 session |
| `web_config` (id, data_key, data_value) | InnoDB | 設定：`web_title`、`web_sub_title`、`page_group_sort`（群組 id JSON 陣列）、`facebook_url`（新增資料列） |

- 應用程式不修改 schema；`backend/internal/db/schema.sql`（舊 dump + `pwd varchar(255)`）僅供 sqlc、本機開發與測試使用。
- 連線使用 `utf8mb4`。建議（非必要）在 migrate 時把各表轉為 `utf8mb4` 與 InnoDB；目前 `utf8` 無法儲存 emoji。
- 涉及排序的異動只碰 InnoDB 表（`page_content`、`page_group`、`web_config`），可放在同一 transaction。

## 3. 架構與技術選型

單一 Go 服務以 `embed` 提供 React 編譯產物與 API；Zeabur 上部署一個 app 服務 + 一個 MySQL。前後端同網域，不需 CORS。

### 後端（Go 1.25+）

- 路由：標準庫 `net/http`（method + path pattern）
- 資料存取：`sqlc` + `go-sql-driver/mysql`
- 設定：環境變數；log：`slog`
- 密碼：bcrypt
- HTML 過濾：`bluemonday`
- 物件儲存：AWS SDK v2（S3 相容）

### 前端（`web/`）

- Vite + React 19 + TypeScript，React Router 7
- TanStack Query
- Tailwind CSS 4；後台使用 shadcn/ui
- 編輯器：Tiptap 3 + 擴充 + HTML 原始碼模式（CodeMirror）
- 拖曳排序：dnd-kit；輪播：Embla Carousel
- `/admin/*` 以 lazy load 載入，前台 bundle 不含後台與 Tiptap

### 目錄結構

```
backend/
  cmd/sniweb/          main：serve / user 子指令
  internal/config/     環境變數載入與驗證
  internal/db/         schema.sql、queries/*.sql、sqlc 產出
  internal/auth/       密碼、token、登入限流、middleware
  internal/storage/    Storage 介面、disk 與 s3 實作
  internal/content/    排序合併、HTML 過濾、上傳檔名與類型判斷
  internal/api/        HTTP handlers（public / admin）
  internal/site/       SPA 靜態檔、meta 注入、/php/picture 相容路由
web/
  src/site/            前台
  src/admin/           後台
  src/shared/          API client、型別
legacy/                舊程式碼（front/、client/、php/、node-vue-ele-app/、舊 build 產物）
Dockerfile
docker-compose.yml     本機 MySQL + MinIO
Makefile
.env.example
```

## 4. API

前綴 `/api/v1`，JSON。錯誤格式：

```json
{ "error": { "code": "group_not_empty", "message": "此頁籤仍有頁面，請先移除或刪除頁面" } }
```

`message` 為可直接顯示的中文。

### 公開

| Method | 路徑 | 回應 |
|---|---|---|
| GET | `/site` | 標題、副標題、facebook_url、排序後的選單（群組 → 頁面 id/名稱） |
| GET | `/home` | 首頁內容、輪播圖、跑馬燈 |
| GET | `/pages/{id}` | 頁面（id、group_id、name、html）、輪播圖、跑馬燈；不存在回 404 |

### 後台（`/api/v1/admin`，需登入）

| Method | 路徑 | 說明 |
|---|---|---|
| POST | `/auth/login` | 帳密登入，設定 cookie |
| POST | `/auth/logout` | 刪除 token |
| GET | `/auth/me` | 目前使用者 |
| GET/POST | `/groups` | 列表（已排序，含頁面）／新增 |
| PATCH/DELETE | `/groups/{id}` | 改名／刪除 |
| PUT | `/groups/order` | 寫入 `web_config.page_group_sort` |
| PUT | `/groups/{id}/pages/order` | 寫入 `page_group.page_sort` |
| POST | `/pages` | 新增（name、group_id、html） |
| GET/PATCH/DELETE | `/pages/{id}` | 取得／更新（可換群組）／刪除 |
| GET/POST | `/carousels` | 列表／新增 |
| PATCH/DELETE | `/carousels/{id}` | 更新／刪除 |
| GET/POST | `/marquees` | 列表／新增 |
| PATCH/DELETE | `/marquees/{id}` | 更新／刪除 |
| GET/POST | `/images` | 列表／上傳（multipart，可多檔） |
| GET | `/images/{name}/usages` | 引用此圖的頁面與輪播圖 |
| DELETE | `/images/{name}` | 刪除 |
| GET/PUT | `/settings` | `web_title`、`web_sub_title`、`facebook_url` |

## 5. 業務規則

### 排序

- 依 JSON 陣列順序排列；不在陣列中的項目依 id 遞增接在最後；陣列中已不存在的 id 忽略（與舊 `GetNav.php` 一致）。
- 空字串或無效 JSON 視為空陣列。
- 新增群組／頁面：id 附加到對應陣列尾端。刪除：從陣列移除。頁面換群組：從舊群組陣列移除、附加到新群組陣列。以上皆在同一 transaction 中完成。
- 排序 API 只接受整數陣列，且必須與實際項目集合相同（不多不少），否則回 400。

### 群組刪除

- 群組內仍有頁面時回 409 `group_not_empty`（舊版會留下孤兒頁面）。
- `id = 1`（首頁群組）不可刪除，回 409 `group_protected`。

### 首頁內容

- 取 `page_group_id = 1` 依排序後的第一個頁面；若無則回空內容。

### 認證

- 登入：以 `account` 查詢 `user`，bcrypt 比對。失敗訊息不區分帳號不存在或密碼錯誤。
- token：`crypto/rand` 32 bytes，base64url 編碼放在 cookie；資料庫 `user_token.token` 存 SHA-256 hex。
- cookie：`sni_session`，`HttpOnly`、`Secure`（本機開發可關閉）、`SameSite=Lax`、`Path=/`。
- 效期 1 小時；每次通過驗證且剩餘時間少於 30 分鐘時延長為 1 小時（sliding）。
- `user_token` 主鍵為 `user_id`，同一帳號同時只有一個 session；新登入會取代舊 session（受 schema 限制，刻意保留）。
- 過期 token 一律拒絕（修正舊版因未定義變數而永不過期的 bug）。
- 登入限流：同一 IP 每 15 分鐘 20 次、同一帳號每 15 分鐘 10 次失敗，超過回 429（記憶體實作，單一實例即可）。
- CSRF：後台所有非 GET 請求檢查 `Origin`（或 `Referer`）必須等於 `PUBLIC_BASE_URL` 或本機開發來源。

### 圖片上傳

- 允許 jpg、png、gif、webp，以 `http.DetectContentType` 依內容判斷，不信任客戶端 MIME。
- 單檔上限 5MB。
- 檔名：`YYYY-MM-DD_HH-mm-ss.<ext>`（Asia/Taipei）；重複時為 `YYYY-MM-DD_HH-mm-ss-<6碼亂數>.<ext>`（修正舊版把亂數接在副檔名後的 bug）。
- 刪除前可查詢使用情況：在 `page_content.html_context` 與 `carousel.image` 中搜尋檔名。前端在有引用時要求再次確認；API 本身不阻擋。
- 檔名只允許 `[A-Za-z0-9._-]`，拒絕路徑穿越。

### 圖片網址與儲存

- 對外網址一律為 `/php/picture/{檔名}`；上傳 API 回傳此相對路徑，內容與 `carousel.image` 也存這個值。
- `Storage` 介面：`Put`、`Delete`、`List`、`Open`／`URL`。
  - `disk`：存於 `STORAGE_DISK_DIR`，Go 直接回傳檔案（`Cache-Control: public, max-age=31536000, immutable`）。
  - `s3`：物件 key 為 `picture/{檔名}`；`/php/picture/{檔名}` 回 302 導向 `S3_PUBLIC_BASE_URL/picture/{檔名}`。
- 以 `STORAGE_DRIVER` 切換。舊內容中的 `http://www.seicho-no-ie.org.tw/php/picture/...` 絕對網址在同網域下可直接使用。

### HTML 內容

- 儲存時以 `bluemonday` 過濾：移除 `<script>`、`on*` 屬性、`javascript:` 網址；保留 `style`、`class`、表格、圖片、常見排版標籤。
- iframe 僅允許 `src` 網域白名單：`www.youtube.com`、`youtube.com`、`www.youtube-nocookie.com`、`players.brightcove.net`、`drive.google.com`、`www.facebook.com`；可由程式常數調整。
- 前台原樣輸出已過濾的 HTML。

## 6. SEO 與 meta 注入

Go 在回應 SPA 的 `index.html` 時，依路徑替換 `<head>` 中的佔位內容：

| 路徑 | title | description | og:image |
|---|---|---|---|
| `/` | `{web_title}｜{web_sub_title}` | 首頁內容純文字前 120 字 | 第一張輪播圖 |
| `/page/{id}` | `{page_name}｜{web_title}` | 頁面純文字前 120 字 | 內容第一張圖，否則第一張輪播圖 |
| 其他 | `{web_title}` | `web_sub_title` | 第一張輪播圖 |

- 同時輸出 `og:title`、`og:description`、`og:url`、`og:type`、`twitter:card`；圖片轉為以 `PUBLIC_BASE_URL` 開頭的絕對網址。
- 所有值經 HTML escape。
- `/page/{id}` 不存在時仍回 SPA，但 HTTP status 為 404；未知路徑同樣回 404。
- `/admin/*` 回 SPA 並加 `noindex`。
- 靜態資源（`/assets/*`）加長效快取；`index.html` 不快取。

## 7. 前台

### 路由

- `/` 首頁、`/page/:id` 內頁、其他顯示 404。
- 舊 hash 網址：載入時若 `location.hash` 符合 `#/(\d+)` 則 `replace` 到 `/page/$1`；`#/` 則到 `/`。
- 內文連結攔截：點擊內容區中的 `<a>` 時，若解析後為同網域且為舊 hash 路由或 `/page/:id`，改用前端路由切換；外部連結維持原行為。

### 版面

- Header：logo + 網站名稱；桌機為下拉選單（hover 與點擊皆可、支援鍵盤）；手機為右側抽屜 + 群組手風琴；捲動時 sticky 並縮小。
- Banner：Embla 輪播，自動播放、箭頭與圓點；固定寬高比、`object-cover`；有 `url` 者可點擊（外部連結開新分頁）。
- 跑馬燈：CSS 動畫，每則使用自己的 `color`，hover 暫停；`prefers-reduced-motion` 時改為靜態輪替。
- 內容區（`.legacy-content`）：`img { max-width: 100%; height: auto }`；iframe 以 16:9 自適應；渲染後將每個 `<table>` 包在 `overflow-x: auto` 容器內。
- Footer：版權文字、Facebook 連結、瀏覽計數器（以 `iframe srcdoc` 載入 `COUNTER_SCRIPT_URL`）。
- GA4：`GA_MEASUREMENT_ID` 有值時載入，路由變更時送出 `page_view`。
- 視覺：主色取自 logo 深藍紫，輔以淺藍漸層；Noto Sans TC；內文 17–18px、行高 1.8；載入時顯示骨架畫面，錯誤時顯示友善提示。
- 前端所需的公開設定（GA ID、計數器網址）由 Go 注入 `index.html` 的 `window.__SNI_CONFIG__`。

## 8. 後台

- `/admin/login`；其餘頁面需要登入。收到 401 時導向登入頁並帶上 `redirect` 參數。
- 左側選單：頁籤與頁面、輪播圖、跑馬燈、圖庫、網站設定、登出。

### 頁籤與頁面

- 左欄是群組列表，右欄是選中群組的頁面；兩欄都可用 dnd-kit 拖曳排序，放開即儲存（樂觀更新，失敗時回復並提示）。
- 群組新增、改名、刪除使用 dialog；刪除失敗時顯示 API 訊息。

### 頁面編輯（`/admin/pages/:id`、`/admin/pages/new?group=:gid`）

- 欄位：頁面名稱、所屬群組、內容。
- Tiptap 擴充：StarterKit、Underline、TextStyle + Color、Highlight、TextAlign、Link、Image、Table 系列、iframe（自訂 node，套用與後端相同的網域白名單）、全域 `style` / `class` 屬性保留（套用於 paragraph、heading、table 系列、image、span）。
- 工具列：標題、粗體、斜體、底線、刪除線、文字顏色、螢光筆、對齊、清單、連結、插入圖片（從圖庫選或上傳）、插入影片（YouTube 或 iframe 網址）、表格操作、復原與重做、切換 HTML 原始碼。
- HTML 原始碼模式：使用 CodeMirror（HTML 語法）。在原始碼模式下存檔會直接送出原始碼，不經過 Tiptap。
- 遺失偵測：載入內容時，將原始 HTML 與「經 Tiptap 解析再輸出」的 HTML 都正規化（DOM 解析、移除空白差異、屬性排序）後比較；若有標籤、屬性或文字遺失，預設進入原始碼模式並顯示警告，說明切換到視覺模式可能遺失格式。
- 預覽（在新分頁開啟 `/page/:id`）；有未存檔的變更時，離開頁面前提示。

### 輪播圖、跑馬燈、圖庫、網站設定

- 輪播圖：縮圖列表（依 id 排序）；新增或編輯的 dialog 可從圖庫選圖或上傳，並填入連結網址。
- 跑馬燈：列表；文字 + 顏色選擇器，即時預覽。
- 圖庫：格狀縮圖，可拖曳多檔上傳（顯示進度），可複製網址；刪除前呼叫 usages API，有引用時列出並要求再次確認。
- 網站設定：表單編輯三個欄位。
- 所有操作結果以 toast 顯示。

## 9. 錯誤處理與安全

- HTTP status：400 驗證錯誤、401 未登入、404 找不到、409 衝突、413 檔案太大、429 過多請求、500 伺服器錯誤。
- 500 時 log 記錄錯誤與 request id，回應只含一般訊息與 request id。
- Middleware：recover、request id、存取 log、body 大小限制（JSON 1MB、上傳每檔 5MB）、安全性 header（CSP 需允許 GA、計數器與 iframe 白名單網域；`X-Content-Type-Options: nosniff`；`Referrer-Policy: strict-origin-when-cross-origin`）。
- `/healthz`：程序存活；`/readyz`：ping 資料庫。
- 程式碼與 repo 中不存放任何密碼。舊 `php/main.php` 內的資料庫密碼已在 git 歷史中外洩，建議更換。

## 10. 設定

| 變數 | 說明 | 預設 |
|---|---|---|
| `PORT` | 監聽 port | `8080` |
| `DATABASE_URL` | MySQL DSN | 必填 |
| `PUBLIC_BASE_URL` | 對外網址 | 必填 |
| `COOKIE_SECURE` | cookie 是否加上 Secure | `true` |
| `DEV_ORIGINS` | 本機開發允許的 Origin（逗號分隔） | 空 |
| `STORAGE_DRIVER` | `disk` 或 `s3` | `disk` |
| `STORAGE_DISK_DIR` | 磁碟目錄 | `/data/picture` |
| `S3_ENDPOINT`、`S3_REGION`、`S3_BUCKET`、`S3_ACCESS_KEY`、`S3_SECRET_KEY`、`S3_PUBLIC_BASE_URL` | S3 相容儲存 | `s3` 時必填 |
| `GA_MEASUREMENT_ID` | GA4 | 空（不載入） |
| `COUNTER_SCRIPT_URL` | 瀏覽計數器 | 空（不顯示） |

啟動時驗證必填值，缺少時直接結束並列出缺少的變數。

## 11. 部署（Zeabur）

- 根目錄 `Dockerfile` 多階段建置：node（pnpm build）→ golang（embed、`CGO_ENABLED=0`）→ `gcr.io/distroless/static`。
- Zeabur：app 服務 + MySQL 服務。使用 `disk` 時掛載 Volume 到 `/data`（重新部署時會短暫停機；首次掛載會清空目錄，舊圖片需在掛載後再放入）；使用 `s3` 時不需要 Volume。
- 服務不會自動執行 migrate。
- 建立帳號：`sniweb user create --account <帳號> --name <暱稱>`（密碼互動輸入，或由 `--password-stdin` 讀取）；重設密碼：`sniweb user passwd --account <帳號>`。

## 12. 本機開發

- `docker compose up`：MySQL 8（自動匯入 `schema.sql` 與 `seed.sql`）+ MinIO。
- `make dev`：同時啟動 `go run ./backend/cmd/sniweb serve` 與 `pnpm --dir web dev`；Vite 將 `/api`、`/php/picture` 代理到 Go。
- `make gen`（sqlc）、`make test`、`make build`、`make lint`。
- `seed.sql`：一個管理者帳號、數個群組與頁面（含表格、iframe、行內樣式的舊式 HTML 範例）、輪播圖與跑馬燈。

## 13. 測試

- 開發一律採 TDD。
- Go 單元測試：排序合併、token／密碼、登入限流、上傳類型與檔名、檔名驗證、HTML 過濾（含 iframe 白名單）、meta 注入、disk storage。
- Go 整合測試：`testcontainers-go` 啟動 MySQL 並匯入 `schema.sql`，涵蓋所有 API 流程（登入、過期、CRUD、排序、刪除規則、usages）；s3 storage 對 MinIO 容器測試。
- 前端（Vitest + Testing Library）：hash 轉址、內文連結攔截、Tiptap 遺失偵測、表格包裝、API client 錯誤處理。
- E2E（Playwright）：前台首頁與內頁、舊 hash 網址轉址、手機寬度無橫向捲動、後台登入 → 編輯頁面 → 前台看到更新、拖曳排序。

## 14. 遷移步驟（repo）

1. 將舊程式碼與 build 產物（`front/`、`client/`、`php/`、`node-vue-ele-app/`、`css/`、`js/`、`fonts/`、`img/`、`index.html`、`favicon.*`、`20191225.sql`）移到 `legacy/`，logo 等需要的素材複製到 `web/public/`。
2. 建立 `backend/`、`web/` 與根目錄的工具檔案。
3. 更新 `README.md`：專案說明、開發、部署、帳號管理。
