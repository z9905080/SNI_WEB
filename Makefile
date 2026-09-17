SQLC_VERSION ?= v1.30.0
MISE := mise exec --

# OrbStack 的 Docker socket 不會被 testcontainers-go 自動偵測到；
# 若尚未設定 DOCKER_HOST 且該 socket 存在，自動補上（其他 Docker 環境不受影響）
ifndef DOCKER_HOST
ORBSTACK_SOCK := $(HOME)/.orbstack/run/docker.sock
ifneq ($(wildcard $(ORBSTACK_SOCK)),)
export DOCKER_HOST := unix://$(ORBSTACK_SOCK)
endif
endif

.PHONY: install up down gen test test-short lint fmt build web dev sample-images db-reset e2e

# 安裝 mise.toml 固定的工具鏈版本（Go、golangci-lint、Node.js、pnpm）
install:
	mise install

up:
	docker compose up -d

down:
	docker compose down

gen:
	$(MISE) go run github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION) generate

test:
	$(MISE) go test -race ./...
	$(MISE) pnpm --dir web test

test-short:
	$(MISE) go test -short ./...
	$(MISE) pnpm --dir web test

lint:
	$(MISE) golangci-lint run ./...
	$(MISE) pnpm --dir web typecheck

fmt:
	$(MISE) golangci-lint fmt ./...

# 建置前端並複製到 Go embed 目錄
web:
	$(MISE) pnpm --dir web install --frozen-lockfile
	$(MISE) pnpm --dir web build
	find backend/internal/site/webdist/dist -mindepth 1 ! -name .gitkeep -exec rm -rf {} +
	cp -R web/dist/. backend/internal/site/webdist/dist/

build: web
	$(MISE) env CGO_ENABLED=0 go build -trimpath -o bin/sniweb ./backend/cmd/sniweb

sample-images:
	mkdir -p tmp/picture
	cp legacy/img/bg1.jpg tmp/picture/sample-1.jpg
	cp legacy/img/bg3.jpg tmp/picture/sample-2.jpg

dev: sample-images
	@trap 'kill 0' INT TERM EXIT; \
	(set -a; . ./.env; set +a; $(MISE) go run ./backend/cmd/sniweb serve) & \
	$(MISE) pnpm --dir web dev & \
	wait

# 把本機資料庫重設為 seed 狀態（E2E 前使用）
db-reset:
	docker compose exec -T mysql sh -c '\
	  mysql -uroot -psniweb -e "DROP DATABASE IF EXISTS sniweb; CREATE DATABASE sniweb CHARACTER SET utf8mb4" && \
	  mysql -uroot -psniweb sniweb < /docker-entrypoint-initdb.d/01-schema.sql && \
	  mysql -uroot -psniweb sniweb < /docker-entrypoint-initdb.d/02-seed.sql'

e2e: web sample-images db-reset
	$(MISE) pnpm --dir web exec playwright install --with-deps chromium
	$(MISE) pnpm --dir web e2e
