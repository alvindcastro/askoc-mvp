package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"askoc-mvp/internal/handlers"
	"askoc-mvp/internal/middleware"
	"askoc-mvp/internal/mock/crm"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	mux.Handle("/api/v1/crm/cases", crm.NewHandler())
	mux.Handle("/healthz", handlers.HealthHandler())
	mux.Handle("/readyz", handlers.ReadyHandler())

	server := &http.Server{
		Addr:              ":9083",
		Handler:           middleware.Chain(mux, middleware.TraceID, middleware.Recover, middleware.RequestLogger(logger, middleware.BasicRedactor)),
		ReadHeaderTimeout: 5 * time.Second,
	}
	run(server, logger, "mock crm")
}

func run(server *http.Server, logger *slog.Logger, name string) {
	go func() {
		logger.Info(name+" listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(name+" stopped", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(name+" shutdown failed", "error", err)
		os.Exit(1)
	}
}
