FROM node:24-alpine AS web
WORKDIR /src/web
RUN corepack enable
COPY web/package.json web/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY web/ ./
RUN pnpm build

FROM golang:1.27-alpine AS server
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY backend/ ./backend/
COPY --from=web /src/web/dist/ ./backend/internal/site/webdist/dist/
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/sniweb ./backend/cmd/sniweb

# 用 alpine 而非 distroless：Zeabur 的 service exec 在 distroless 上會回 INTERNAL_ERROR，
# 而部署後套用資料庫遷移要靠它。以 root 執行才能寫入 Zeabur 掛載的 /data Volume。
# alpine 不像 distroless/static 內建 CA 憑證，STORAGE_DRIVER=s3 需要，必須自己裝。
FROM alpine:3.22
RUN apk add --no-cache ca-certificates
COPY --from=server /out/sniweb /sniweb
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/sniweb"]
CMD ["serve"]
