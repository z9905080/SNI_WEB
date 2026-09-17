// Command sniweb 是生長之家網站的伺服器與管理工具。
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/z9905080/SNI_WEB/backend/internal/api"
	"github.com/z9905080/SNI_WEB/backend/internal/auth"
	"github.com/z9905080/SNI_WEB/backend/internal/config"
	"github.com/z9905080/SNI_WEB/backend/internal/db"
	"github.com/z9905080/SNI_WEB/backend/internal/db/dbgen"
	"github.com/z9905080/SNI_WEB/backend/internal/site"
	"github.com/z9905080/SNI_WEB/backend/internal/site/webdist"
	"github.com/z9905080/SNI_WEB/backend/internal/storage"
	"github.com/z9905080/SNI_WEB/backend/internal/store"
)

const usage = `用法：
  sniweb serve
  sniweb user create --account <帳號> --name <暱稱> [--password-stdin]
  sniweb user passwd --account <帳號> [--password-stdin]
`

var errUsage = errors.New("usage")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := run(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr, os.Getenv)
	stop()
	os.Exit(code)
}

func run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer, getenv func(string) string) int {
	err := errUsage
	switch {
	case len(args) == 1 && args[0] == "serve":
		err = serve(ctx, getenv, stderr)
	case len(args) >= 2 && args[0] == "user":
		err = userCmd(ctx, args[1], args[2:], stdin, stdout, stderr, getenv)
	}
	switch {
	case errors.Is(err, errUsage):
		fmt.Fprint(stderr, usage)
		return 2
	case err != nil:
		fmt.Fprintln(stderr, "錯誤：", err)
		return 1
	}
	return 0
}

func serve(ctx context.Context, getenv func(string) string, logOut io.Writer) error {
	cfg, err := config.Load(getenv)
	if err != nil {
		return err
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(logOut, nil)))

	conn, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	var images storage.Storage
	if cfg.StorageDriver == "s3" {
		images = storage.NewS3(storage.S3Options(cfg.S3))
	} else {
		disk, err := storage.NewDisk(cfg.StorageDiskDir)
		if err != nil {
			return err
		}
		images = disk
	}

	content := store.New(conn)
	apiHandler := api.New(api.Deps{
		Store:          content,
		Sessions:       auth.NewSessions(dbgen.New(conn), time.Now),
		Storage:        images,
		CookieSecure:   cfg.CookieSecure,
		AllowedOrigins: cfg.AllowedOrigins(),
		Ping:           conn.PingContext,
	}).Routes()

	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler)
	mux.Handle("GET /healthz", apiHandler)
	mux.Handle("GET /readyz", apiHandler)
	mux.Handle("/", site.New(webdist.FS, content, images, cfg.PublicBaseURL, site.PublicConfig{
		GAMeasurementID:  cfg.GAMeasurementID,
		CounterScriptURL: cfg.CounterScriptURL,
	}))

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           api.Wrap(mux, api.BuildCSP(cfg.CounterScriptURL)),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       2 * time.Minute, // 多檔上傳
		WriteTimeout:      2 * time.Minute,
		IdleTimeout:       2 * time.Minute,
	}
	errCh := make(chan error, 1)
	go func() {
		slog.Info("listening", "addr", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
