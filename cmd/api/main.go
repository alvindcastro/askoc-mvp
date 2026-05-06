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
	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/config"
	"askoc-mvp/internal/handlers"
	"askoc-mvp/internal/llm"
	"askoc-mvp/internal/middleware"
	"askoc-mvp/internal/orchestrator"
	"askoc-mvp/internal/rag"
	"askoc-mvp/internal/tools"
	"askoc-mvp/internal/workflow"
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
	toolHTTPClient := &http.Client{Timeout: cfg.Workflow.Timeout}
	llmPort, err := buildLLM(cfg, &http.Client{Timeout: cfg.Provider.Timeout})
	if err != nil {
		logger.Error("create llm provider", "error", err)
		os.Exit(1)
	}
	auditStore := audit.NewMemoryStore()
	retriever := buildRetriever(context.Background(), cfg.RAG.ChunksPath, logger)
	workflowPort, err := buildWorkflow(cfg, toolHTTPClient)
	if err != nil {
		logger.Error("create workflow client", "error", err)
		os.Exit(1)
	}
	chatService, err := orchestrator.New(orchestrator.Dependencies{
		Classifier: buildClassifier(cfg, llmPort, auditStore),
		Retriever:  retriever,
		LLM:        llmPort,
		Banner:     tools.NewBannerClient(cfg.Integrations.BannerURL, toolHTTPClient),
		Payment:    tools.NewPaymentClient(cfg.Integrations.PaymentURL, toolHTTPClient),
		Workflow:   workflowPort,
		CRM:        tools.NewCRMClient(cfg.Integrations.CRMURL, toolHTTPClient),
		Audit:      auditStore,
	})
	if err != nil {
		logger.Error("create orchestrator", "error", err)
		os.Exit(1)
	}
	mux.Handle("/", handlers.ChatPageHandler("web/templates/chat.html"))
	mux.Handle("/chat", handlers.ChatPageHandler("web/templates/chat.html"))
	mux.Handle("/admin", handlers.AdminPageHandler("web/templates/admin.html"))
	mux.Handle("/static/", http.StripPrefix("/static/", handlers.StaticFileHandler("web/static")))
	mux.Handle("/api/v1/chat", handlers.ChatHandler(chatService))
	mux.Handle("/api/v1/admin/metrics", handlers.AdminMetricsHandler(auditStore, adminAccessToken(cfg)))
	mux.Handle("/api/v1/admin/audit/export", handlers.AdminAuditExportHandler(auditStore, adminAccessToken(cfg)))
	mux.Handle("/api/v1/admin/audit/reset", handlers.AdminAuditResetHandler(auditStore, adminAccessToken(cfg)))
	mux.Handle("/api/v1/admin/audit/purge", handlers.AdminAuditPurgeHandler(auditStore, adminAccessToken(cfg)))
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

func buildLLM(cfg config.Config, httpClient *http.Client) (orchestrator.LLM, error) {
	if cfg.Provider.Mode != "openai-compatible" {
		return orchestrator.DisabledLLM{}, nil
	}
	provider, err := llm.NewOpenAIClient(llm.OpenAIClientConfig{
		Endpoint:   cfg.Provider.Endpoint,
		APIKey:     cfg.Provider.APIKey,
		Model:      cfg.Provider.Model,
		HTTPClient: httpClient,
		Timeout:    cfg.Provider.Timeout,
	})
	if err != nil {
		return nil, err
	}
	return providerLLM{
		provider: provider,
		timeout:  cfg.Provider.Timeout,
	}, nil
}

func buildClassifier(cfg config.Config, llmPort orchestrator.LLM, recorder audit.Recorder) orchestrator.IntentClassifier {
	fallback := classifier.Fallback{}
	if cfg.Provider.Mode != "openai-compatible" {
		return fallback
	}
	return orchestrator.LLMBackedClassifier{
		LLM:      llmPort,
		Fallback: fallback,
		Parse:    classifier.ParseLLMClassificationOutput,
		Audit:    recorder,
	}
}

func buildWorkflow(cfg config.Config, httpClient *http.Client) (orchestrator.WorkflowTool, error) {
	if cfg.Workflow.URL == "" {
		return workflow.NewInMemoryClient(), nil
	}
	return workflow.NewPowerAutomateClient(workflow.PowerAutomateClientConfig{
		WebhookURL:      cfg.Workflow.URL,
		HTTPClient:      httpClient,
		Signature:       cfg.Workflow.Signature,
		SignatureHeader: cfg.Workflow.SignatureHeader,
		MaxRetries:      cfg.Workflow.MaxRetries,
	})
}

func adminAccessToken(cfg config.Config) string {
	if cfg.Auth.Token != "" {
		return cfg.Auth.Token
	}
	return "demo-admin-token"
}

type providerLLM struct {
	provider llm.Provider
	timeout  time.Duration
}

func (p providerLLM) GenerateAnswer(ctx context.Context, prompt string) (string, error) {
	resp, err := p.provider.GenerateAnswer(ctx, llm.AnswerRequest{
		Messages: []llm.Message{
			{Role: llm.RoleUser, Content: prompt},
		},
		MaxTokens:   700,
		Temperature: 0.2,
		Timeout:     p.timeout,
	})
	if err != nil {
		return "", err
	}
	return resp.Answer, nil
}

func buildRetriever(ctx context.Context, chunksPath string, logger *slog.Logger) orchestrator.Retriever {
	chunks, err := rag.LoadChunks(ctx, chunksPath)
	if err != nil {
		logger.Warn("local RAG chunks unavailable; source-grounded answers will use safe fallback", "path", chunksPath, "error", err)
		return orchestrator.DisabledRetriever{}
	}
	if len(chunks) == 0 {
		logger.Warn("local RAG chunks file is empty; source-grounded answers will use safe fallback", "path", chunksPath)
		return orchestrator.DisabledRetriever{}
	}
	return rag.NewLocalRetriever(chunks)
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
