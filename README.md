# SNI_WEB

生長之家網站（前台 + 後台管理）。Go 後端 + React 前端，單一服務部署。

## 本機開發

需要 Docker 與 [mise](https://mise.jdx.dev/)（工具鏈版本固定在 `mise.toml`：Go 1.27、golangci-lint v2、Node.js 24、pnpm 10）：

```bash
brew install mise   # 其他安裝方式見 mise 官網
make install        # 安裝 mise.toml 固定的版本
```

```bash
make up                      # MySQL（含 schema 與範例資料）+ MinIO
cp .env.example .env
set -a; . ./.env; set +a
mise exec -- go run ./backend/cmd/sniweb serve
```

範例資料的管理者帳號：`admin` / `admin1234`（只供本機使用）。

所有 `make` 指令都透過 `mise exec --` 呼叫，固定使用 `mise.toml` 的版本，不需要先啟用 mise shell 整合。在 macOS 上若使用 OrbStack，`make` 會自動偵測並設定 `DOCKER_HOST`；直接執行 `go test`／`go run`（不透過 `make`）時才需要自行設定：

```bash
export DOCKER_HOST=unix://$HOME/.orbstack/run/docker.sock
```

常用指令：

| 指令 | 說明 |
|---|---|
| `make install` | 安裝 `mise.toml` 固定的工具鏈版本 |
| `make gen` | 由 `backend/internal/db/queries/*.sql` 產生 sqlc 程式碼 |
| `make test` | 全部測試（需要 Docker） |
| `make test-short` | 略過需要 Docker 的測試 |
| `make lint` | golangci-lint（含 gofumpt 排版檢查，設定在 `.golangci.yml`） |
| `make fmt` | 以 golangci-lint（gofumpt）自動排版 |
| `make build` | 建置 `bin/sniweb` |

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

## 資料庫

沿用舊站的 7 張表，應用程式不會自動修改 schema。部署前需手動執行：

```sql
ALTER TABLE `user` MODIFY `pwd` varchar(255) NOT NULL COMMENT '密碼（bcrypt）';
INSERT INTO `web_config` (`data_key`, `data_value`) VALUES ('facebook_url', '');
DELETE FROM `user_token`;
```

舊密碼（MD5）與舊 session 無法沿用；`DELETE FROM user_token` 清掉舊格式的登入紀錄（新版以 SHA-256 雜湊比對，舊資料列不會被用到，純粹是遷移時的清理）。請用下列指令重設密碼。

## 帳號管理

```bash
sniweb user create --account <帳號> --name <暱稱>      # 互動輸入密碼
echo -n '<密碼>' | sniweb user passwd --account <帳號> --password-stdin
```

## 環境變數

見 `.env.example` 與 `docs/superpowers/specs/2026-09-16-sni-web-rewrite-design.md` 第 10 節。

## 安全提醒

舊的 `legacy/php/main.php` 內含資料庫密碼且已存在於 git 歷史，請更換該密碼。
