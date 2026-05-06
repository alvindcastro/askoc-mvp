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

	"askoc-mvp/internal/config"
	"askoc-mvp/internal/handlers"
	"askoc-mvp/internal/middleware"
	"askoc-mvp/internal/session"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel(cfg.LogLevel),
	}))

	mux := http.NewServeMux()
	chatStore := session.NewStore(30 * time.Minute)
	chatService := handlers.NewPlaceholderChatService(chatStore)
	mux.Handle("/", handlers.ChatPageHandler("web/templates/chat.html"))
	mux.Handle("/chat", handlers.ChatPageHandler("web/templates/chat.html"))
	mux.Handle("/static/", http.StripPrefix("/static/", handlers.StaticFileHandler("web/static")))
	mux.Handle("/api/v1/chat", handlers.ChatHandler(chatService))
	mux.Handle("/healthz", handlers.HealthHandler())
	mux.Handle("/readyz", handlers.ReadyHandler())

	handler := middleware.Chain(
		mux,
		middleware.TraceID,
		middleware.Recover,
		middleware.RequestLogger(logger, middleware.BasicRedactor),
		middleware.MockAuth(cfg.Auth.Enabled, cfg.Auth.Token),
	)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("api listening", "addr", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("api stopped", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("api shutdown failed", "error", err)
		os.Exit(1)
	}
}

func logLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
