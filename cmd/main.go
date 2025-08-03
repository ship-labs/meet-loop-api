package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ship-labs/meet-loop-api/config"
	"github.com/ship-labs/meet-loop-api/database"
	"github.com/ship-labs/meet-loop-api/middleware"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger = logger.With("app", "MeetLoop")
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		exit(err, "config.LoadConfig()")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	conn, err := database.Dial(timeoutCtx, cfg.DBURL)
	if err != nil {
		exit(err, "database.Dial")
	}
	defer conn.Close()
	slog.InfoContext(ctx, "main", "message", "Connected to database successfully")

	port := cmp.Or(cfg.Port, config.DefaultPort)
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", middleware.Auth(func(w http.ResponseWriter, r *http.Request) middleware.Handler {
		return middleware.JSON(middleware.Response{Message: http.StatusText(http.StatusOK)})
	}))

	handler := middleware.CorsMiddleware(mux)
	handler = middleware.LoggingMiddleware(handler)
	server := http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%d", port),
	}

	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			exit(err, "server.ListenAndServe")
		}
	}()

	slog.Info("main", "message", "Server started successfully", "port", server.Addr, "numCPUS", runtime.NumCPU())

	<-ctx.Done()
	timeoutCtx, cancel = context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		exit(err, "server.Shutdown")
	}

	slog.Info("main", "message", "Server shutdown successfully")
}

func exit(err error, origin string) {
	slog.Error(origin, "error", err)
	os.Exit(1)
}
