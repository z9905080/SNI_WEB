# SNI_WEB

生長之家網站（前台 + 後台管理）。Go 後端 + React 前端，單一服務部署。

## 本機開發

需要 Docker 與 [mise](https://mise.jdx.dev/)。工具鏈版本固定在 `mise.toml`（Go 1.27、golangci-lint v2、Node.js 24、pnpm 10）：

```bash
brew install mise   # 其他安裝方式見 mise 官網
mise install
```

```bash
make up                      # MySQL（含 schema 與範例資料）+ MinIO
cp .env.example .env
set -a; . ./.env; set +a
go run ./backend/cmd/sniweb serve
```

範例資料的管理者帳號：`admin` / `admin1234`（只供本機使用）。

在 macOS 上若使用 OrbStack，需要 Docker 的測試（`make test`、`go test ./backend/internal/store/`、`go test ./backend/cmd/sniweb/` 等）要先設定：

```bash
export DOCKER_HOST=unix://$HOME/.orbstack/run/docker.sock
```

常用指令：

| 指令 | 說明 |
|---|---|
| `make gen` | 由 `backend/internal/db/queries/*.sql` 產生 sqlc 程式碼 |
| `make test` | 全部測試（需要 Docker） |
| `make test-short` | 略過需要 Docker 的測試 |
| `make lint` | golangci-lint（含 gofumpt 排版檢查，設定在 `.golangci.yml`） |
| `golangci-lint fmt ./...` | 以 gofumpt 自動排版 |
| `make build` | 建置 `bin/sniweb` |

## 資料庫

沿用舊站的 7 張表，應用程式不會自動修改 schema。部署前需手動執行：

```sql
ALTER TABLE `user` MODIFY `pwd` varchar(255) NOT NULL COMMENT '密碼（bcrypt）';
INSERT INTO `web_config` (`data_key`, `data_value`) VALUES ('facebook_url', '');
```

舊密碼（MD5）與舊 session 無法沿用，請用下列指令重設。

## 帳號管理

```bash
sniweb user create --account <帳號> --name <暱稱>      # 互動輸入密碼
echo -n '<密碼>' | sniweb user passwd --account <帳號> --password-stdin
```

## 環境變數

見 `.env.example` 與 `docs/superpowers/specs/2026-09-16-sni-web-rewrite-design.md` 第 10 節。

## 安全提醒

舊的 `legacy/php/main.php` 內含資料庫密碼且已存在於 git 歷史，請更換該密碼。
