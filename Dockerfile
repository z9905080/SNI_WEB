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

# 使用 root 映像，才能寫入 Zeabur 掛載的 /data Volume
FROM gcr.io/distroless/static-debian12
COPY --from=server /out/sniweb /sniweb
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/sniweb"]
CMD ["serve"]
