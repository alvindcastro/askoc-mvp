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

	"askoc-mvp/internal/audit"
	"askoc-mvp/internal/handlers"
	"askoc-mvp/internal/middleware"
	"askoc-mvp/internal/workflow"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	auditStore := audit.NewMemoryStore()

	mux := http.NewServeMux()
	mux.Handle(workflow.PaymentReminderPath, workflow.NewSimulatorHandler(workflow.NewInMemoryClient(), auditStore))
	mux.Handle("/api/v1/admin/metrics", handlers.AdminMetricsHandler(auditStore, "demo-admin-token"))
	mux.Handle("/api/v1/admin/audit/export", handlers.AdminAuditExportHandler(auditStore, "demo-admin-token"))
	mux.Handle("/healthz", handlers.HealthHandler())
	mux.Handle("/readyz", handlers.ReadyHandler())

	server := &http.Server{
		Addr:              ":9084",
		Handler:           middleware.Chain(mux, middleware.TraceID, middleware.Recover, middleware.RequestLogger(logger, middleware.BasicRedactor)),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("workflow simulator listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("workflow simulator stopped", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("workflow simulator shutdown failed", "error", err)
		os.Exit(1)
	}
}
