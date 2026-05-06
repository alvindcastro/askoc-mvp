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

	"askoc-mvp/internal/fixtures"
	"askoc-mvp/internal/handlers"
	"askoc-mvp/internal/middleware"
	paymentmock "askoc-mvp/internal/mock/payment"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	records, err := fixtures.Load(context.Background(), "data/synthetic-students.json")
	if err != nil {
		logger.Error("load synthetic fixtures", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/students/", paymentmock.NewHandler(records))
	mux.Handle("/healthz", handlers.HealthHandler())
	mux.Handle("/readyz", handlers.ReadyHandler())

	server := &http.Server{
		Addr:              ":8082",
		Handler:           middleware.Chain(mux, middleware.TraceID, middleware.Recover, middleware.RequestLogger(logger, middleware.BasicRedactor)),
		ReadHeaderTimeout: 5 * time.Second,
	}
	run(server, logger, "mock payment")
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
