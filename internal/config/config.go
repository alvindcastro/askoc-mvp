package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr     string
	LogLevel     string
	Auth         AuthConfig
	Workflow     WorkflowConfig
	Provider     ProviderConfig
	Integrations IntegrationConfig
	RAG          RAGConfig
}

type AuthConfig struct {
	Enabled bool
	Token   string
}

type WorkflowConfig struct {
	URL             string
	Timeout         time.Duration
	Signature       string
	SignatureHeader string
	MaxRetries      int
}

type ProviderConfig struct {
	Mode     string
	Model    string
	Endpoint string
	APIKey   string
	Timeout  time.Duration
}

type IntegrationConfig struct {
	BannerURL  string
	PaymentURL string
	CRMURL     string
}

type RAGConfig struct {
	ChunksPath string
}

func Load() (Config, error) {
	return LoadFromEnv(snapshotEnv())
}

func LoadFromEnv(env map[string]string) (Config, error) {
	cfg := Config{
		HTTPAddr: value(env, "ASKOC_HTTP_ADDR", ":8080"),
		LogLevel: value(env, "ASKOC_LOG_LEVEL", "info"),
		Auth: AuthConfig{
			Enabled: false,
			Token:   value(env, "ASKOC_AUTH_TOKEN", ""),
		},
		Workflow: WorkflowConfig{
			URL:             value(env, "ASKOC_WORKFLOW_URL", ""),
			Timeout:         5 * time.Second,
			Signature:       value(env, "ASKOC_WORKFLOW_SIGNATURE", ""),
			SignatureHeader: value(env, "ASKOC_WORKFLOW_SIGNATURE_HEADER", "X-AskOC-Workflow-Signature"),
			MaxRetries:      1,
		},
		Provider: ProviderConfig{
			Mode:     value(env, "ASKOC_PROVIDER", "stub"),
			Model:    value(env, "ASKOC_PROVIDER_MODEL", "demo-placeholder"),
			Endpoint: value(env, "ASKOC_PROVIDER_ENDPOINT", ""),
			APIKey:   value(env, "ASKOC_PROVIDER_API_KEY", ""),
			Timeout:  5 * time.Second,
		},
		Integrations: IntegrationConfig{
			BannerURL:  value(env, "ASKOC_BANNER_URL", "http://localhost:8081"),
			PaymentURL: value(env, "ASKOC_PAYMENT_URL", "http://localhost:8082"),
			CRMURL:     value(env, "ASKOC_CRM_URL", "http://localhost:8083"),
		},
		RAG: RAGConfig{
			ChunksPath: value(env, "ASKOC_RAG_CHUNKS_PATH", "data/rag-chunks.json"),
		},
	}

	if strings.TrimSpace(cfg.HTTPAddr) == "" {
		return Config{}, fmt.Errorf("ASKOC_HTTP_ADDR must not be empty")
	}

	authEnabled, err := parseBool(env, "ASKOC_AUTH_ENABLED", false)
	if err != nil {
		return Config{}, err
	}
	cfg.Auth.Enabled = authEnabled

	timeout, err := parsePositiveSeconds(env, "ASKOC_WORKFLOW_TIMEOUT_SECONDS", 5)
	if err != nil {
		return Config{}, err
	}
	cfg.Workflow.Timeout = timeout
	maxRetries, err := parseNonNegativeInt(env, "ASKOC_WORKFLOW_MAX_RETRIES", 1)
	if err != nil {
		return Config{}, err
	}
	cfg.Workflow.MaxRetries = maxRetries
	if strings.TrimSpace(cfg.Workflow.SignatureHeader) == "" {
		return Config{}, fmt.Errorf("ASKOC_WORKFLOW_SIGNATURE_HEADER must not be empty")
	}

	if !validLogLevel(cfg.LogLevel) {
		return Config{}, fmt.Errorf("ASKOC_LOG_LEVEL must be one of debug, info, warn, error")
	}
	if !validProviderMode(cfg.Provider.Mode) {
		return Config{}, fmt.Errorf("ASKOC_PROVIDER must be one of stub, openai-compatible")
	}
	providerTimeout, err := parsePositiveSeconds(env, "ASKOC_PROVIDER_TIMEOUT_SECONDS", 5)
	if err != nil {
		return Config{}, err
	}
	cfg.Provider.Timeout = providerTimeout
	if cfg.Provider.Mode == "openai-compatible" {
		if strings.TrimSpace(cfg.Provider.Endpoint) == "" {
			return Config{}, fmt.Errorf("ASKOC_PROVIDER_ENDPOINT must not be empty when ASKOC_PROVIDER=openai-compatible")
		}
		if strings.TrimSpace(cfg.Provider.APIKey) == "" {
			return Config{}, fmt.Errorf("ASKOC_PROVIDER_API_KEY must not be empty when ASKOC_PROVIDER=openai-compatible")
		}
	}
	if strings.TrimSpace(cfg.RAG.ChunksPath) == "" {
		return Config{}, fmt.Errorf("ASKOC_RAG_CHUNKS_PATH must not be empty")
	}

	return cfg, nil
}

func (c Config) String() string {
	authToken := ""
	if c.Auth.Token != "" {
		authToken = "REDACTED"
	}
	providerKey := ""
	if c.Provider.APIKey != "" {
		providerKey = "REDACTED"
	}
	return fmt.Sprintf(
		"http_addr:%s log_level:%s auth_enabled:%t auth_token:%s workflow_url:%s workflow_timeout:%s workflow_signature:%s workflow_signature_header:%s workflow_max_retries:%d banner_url:%s payment_url:%s crm_url:%s rag_chunks_path:%s provider:%s provider_model:%s provider_endpoint:%s provider_timeout:%s provider_api_key:%s",
		c.HTTPAddr,
		c.LogLevel,
		c.Auth.Enabled,
		authToken,
		redactIfPresent(c.Workflow.URL),
		c.Workflow.Timeout,
		redactIfPresent(c.Workflow.Signature),
		c.Workflow.SignatureHeader,
		c.Workflow.MaxRetries,
		c.Integrations.BannerURL,
		c.Integrations.PaymentURL,
		c.Integrations.CRMURL,
		c.RAG.ChunksPath,
		c.Provider.Mode,
		c.Provider.Model,
		c.Provider.Endpoint,
		c.Provider.Timeout,
		providerKey,
	)
}

func snapshotEnv() map[string]string {
	env := make(map[string]string)
	for _, pair := range os.Environ() {
		key, val, ok := strings.Cut(pair, "=")
		if ok {
			env[key] = val
		}
	}
	return env
}

func value(env map[string]string, key, fallback string) string {
	if got, ok := env[key]; ok {
		return got
	}
	return fallback
}

func parseBool(env map[string]string, key string, fallback bool) (bool, error) {
	raw, ok := env[key]
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback, nil
	}
	got, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean: %w", key, err)
	}
	return got, nil
}

func parsePositiveSeconds(env map[string]string, key string, fallback int) (time.Duration, error) {
	raw, ok := env[key]
	if !ok || strings.TrimSpace(raw) == "" {
		return time.Duration(fallback) * time.Second, nil
	}
	seconds, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be a whole number of seconds: %w", key, err)
	}
	if seconds <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", key)
	}
	return time.Duration(seconds) * time.Second, nil
}

func parseNonNegativeInt(env map[string]string, key string, fallback int) (int, error) {
	raw, ok := env[key]
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback, nil
	}
	got, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be a whole number: %w", key, err)
	}
	if got < 0 {
		return 0, fmt.Errorf("%s must not be negative", key)
	}
	return got, nil
}

func redactIfPresent(value string) string {
	if value == "" {
		return ""
	}
	return "REDACTED"
}

func validLogLevel(level string) bool {
	switch level {
	case "debug", "info", "warn", "error":
		return true
	default:
		return false
	}
}

func validProviderMode(mode string) bool {
	switch mode {
	case "stub", "openai-compatible":
		return true
	default:
		return false
	}
}
