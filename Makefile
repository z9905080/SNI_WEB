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

.PHONY: install up down gen test test-short lint fmt build

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

test-short:
	$(MISE) go test -short ./...

lint:
	$(MISE) golangci-lint run ./...

fmt:
	$(MISE) golangci-lint fmt ./...

build:
	$(MISE) env CGO_ENABLED=0 go build -trimpath -o bin/sniweb ./backend/cmd/sniweb
