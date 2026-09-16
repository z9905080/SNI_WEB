SQLC_VERSION ?= v1.30.0

.PHONY: up down gen test test-short lint build

up:
	docker compose up -d

down:
	docker compose down

gen:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION) generate

test:
	go test -race ./...

test-short:
	go test -short ./...

lint:
	@test -z "$$(gofmt -l backend)" || (gofmt -l backend; exit 1)
	go vet ./...

build:
	CGO_ENABLED=0 go build -trimpath -o bin/sniweb ./backend/cmd/sniweb
