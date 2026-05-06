package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadFromEnvUsesSafeDefaults(t *testing.T) {
	cfg, err := LoadFromEnv(map[string]string{})
	if err != nil {
		t.Fatalf("LoadFromEnv returned error: %v", err)
	}

	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("HTTPAddr = %q, want :8080", cfg.HTTPAddr)
	}
	if cfg.Auth.Enabled {
		t.Fatalf("Auth.Enabled = true, want false for local demo default")
	}
	if cfg.Auth.Token != "" {
		t.Fatalf("Auth.Token = %q, want empty default", cfg.Auth.Token)
	}
	if cfg.LogLevel != "info" {
		t.Fatalf("LogLevel = %q, want info", cfg.LogLevel)
	}
	if cfg.Workflow.URL != "" {
		t.Fatalf("Workflow.URL = %q, want empty default", cfg.Workflow.URL)
	}
	if cfg.Workflow.Timeout != 5*time.Second {
		t.Fatalf("Workflow.Timeout = %s, want 5s", cfg.Workflow.Timeout)
	}
	if cfg.Workflow.Signature != "" {
		t.Fatalf("Workflow.Signature = %q, want empty default", cfg.Workflow.Signature)
	}
	if cfg.Workflow.SignatureHeader != "X-AskOC-Workflow-Signature" {
		t.Fatalf("Workflow.SignatureHeader = %q, want default signature header", cfg.Workflow.SignatureHeader)
	}
	if cfg.Workflow.MaxRetries != 1 {
		t.Fatalf("Workflow.MaxRetries = %d, want one retry by default", cfg.Workflow.MaxRetries)
	}
	if cfg.Integrations.BannerURL != "http://localhost:8081" {
		t.Fatalf("Integrations.BannerURL = %q, want mock Banner default", cfg.Integrations.BannerURL)
	}
	if cfg.Integrations.PaymentURL != "http://localhost:8082" {
		t.Fatalf("Integrations.PaymentURL = %q, want mock Payment default", cfg.Integrations.PaymentURL)
	}
	if cfg.Integrations.CRMURL != "http://localhost:8083" {
		t.Fatalf("Integrations.CRMURL = %q, want mock CRM default", cfg.Integrations.CRMURL)
	}
	if cfg.Provider.Mode != "stub" {
		t.Fatalf("Provider.Mode = %q, want stub", cfg.Provider.Mode)
	}
	if cfg.Provider.Model != "demo-placeholder" {
		t.Fatalf("Provider.Model = %q, want demo-placeholder", cfg.Provider.Model)
	}
	if cfg.Provider.Endpoint != "" {
		t.Fatalf("Provider.Endpoint = %q, want empty default", cfg.Provider.Endpoint)
	}
	if cfg.Provider.Timeout != 5*time.Second {
		t.Fatalf("Provider.Timeout = %s, want 5s", cfg.Provider.Timeout)
	}
	if cfg.RAG.ChunksPath != "data/rag-chunks.json" {
		t.Fatalf("RAG.ChunksPath = %q, want data/rag-chunks.json", cfg.RAG.ChunksPath)
	}
}

func TestLoadFromEnvUsesOverrides(t *testing.T) {
	cfg, err := LoadFromEnv(map[string]string{
		"ASKOC_HTTP_ADDR":                 "127.0.0.1:9090",
		"ASKOC_AUTH_ENABLED":              "true",
		"ASKOC_AUTH_TOKEN":                "demo-token",
		"ASKOC_LOG_LEVEL":                 "debug",
		"ASKOC_WORKFLOW_URL":              "http://workflow.local/hook",
		"ASKOC_WORKFLOW_TIMEOUT_SECONDS":  "12",
		"ASKOC_WORKFLOW_SIGNATURE":        "workflow-secret",
		"ASKOC_WORKFLOW_SIGNATURE_HEADER": "X-Demo-Signature",
		"ASKOC_WORKFLOW_MAX_RETRIES":      "2",
		"ASKOC_BANNER_URL":                "http://banner.local",
		"ASKOC_PAYMENT_URL":               "http://payment.local",
		"ASKOC_CRM_URL":                   "http://crm.local",
		"ASKOC_RAG_CHUNKS_PATH":           "tmp/test-rag-chunks.json",
		"ASKOC_PROVIDER":                  "openai-compatible",
		"ASKOC_PROVIDER_MODEL":            "gpt-demo",
		"ASKOC_PROVIDER_ENDPOINT":         "http://llm.local/v1/chat/completions",
		"ASKOC_PROVIDER_API_KEY":          "sk-demo-secret",
		"ASKOC_PROVIDER_TIMEOUT_SECONDS":  "9",
	})
	if err != nil {
		t.Fatalf("LoadFromEnv returned error: %v", err)
	}

	if cfg.HTTPAddr != "127.0.0.1:9090" {
		t.Fatalf("HTTPAddr = %q", cfg.HTTPAddr)
	}
	if !cfg.Auth.Enabled || cfg.Auth.Token != "demo-token" {
		t.Fatalf("Auth = %+v, want enabled with token", cfg.Auth)
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("LogLevel = %q", cfg.LogLevel)
	}
	if cfg.Workflow.URL != "http://workflow.local/hook" {
		t.Fatalf("Workflow.URL = %q", cfg.Workflow.URL)
	}
	if cfg.Workflow.Timeout != 12*time.Second {
		t.Fatalf("Workflow.Timeout = %s", cfg.Workflow.Timeout)
	}
	if cfg.Workflow.Signature != "workflow-secret" {
		t.Fatalf("Workflow.Signature was not loaded")
	}
	if cfg.Workflow.SignatureHeader != "X-Demo-Signature" {
		t.Fatalf("Workflow.SignatureHeader = %q", cfg.Workflow.SignatureHeader)
	}
	if cfg.Workflow.MaxRetries != 2 {
		t.Fatalf("Workflow.MaxRetries = %d", cfg.Workflow.MaxRetries)
	}
	if cfg.Integrations.BannerURL != "http://banner.local" || cfg.Integrations.PaymentURL != "http://payment.local" || cfg.Integrations.CRMURL != "http://crm.local" {
		t.Fatalf("Integrations = %+v", cfg.Integrations)
	}
	if cfg.Provider.Mode != "openai-compatible" || cfg.Provider.Model != "gpt-demo" {
		t.Fatalf("Provider = %+v", cfg.Provider)
	}
	if cfg.Provider.Endpoint != "http://llm.local/v1/chat/completions" {
		t.Fatalf("Provider.Endpoint = %q", cfg.Provider.Endpoint)
	}
	if cfg.Provider.APIKey != "sk-demo-secret" {
		t.Fatalf("Provider.APIKey was not loaded")
	}
	if cfg.Provider.Timeout != 9*time.Second {
		t.Fatalf("Provider.Timeout = %s", cfg.Provider.Timeout)
	}
	if cfg.RAG.ChunksPath != "tmp/test-rag-chunks.json" {
		t.Fatalf("RAG.ChunksPath = %q", cfg.RAG.ChunksPath)
	}
}

func TestLoadFromEnvRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr string
	}{
		{
			name:    "invalid auth boolean",
			env:     map[string]string{"ASKOC_AUTH_ENABLED": "sometimes"},
			wantErr: "ASKOC_AUTH_ENABLED",
		},
		{
			name:    "invalid workflow timeout",
			env:     map[string]string{"ASKOC_WORKFLOW_TIMEOUT_SECONDS": "soon"},
			wantErr: "ASKOC_WORKFLOW_TIMEOUT_SECONDS",
		},
		{
			name:    "negative workflow timeout",
			env:     map[string]string{"ASKOC_WORKFLOW_TIMEOUT_SECONDS": "-1"},
			wantErr: "ASKOC_WORKFLOW_TIMEOUT_SECONDS",
		},
		{
			name:    "invalid workflow retries",
			env:     map[string]string{"ASKOC_WORKFLOW_MAX_RETRIES": "many"},
			wantErr: "ASKOC_WORKFLOW_MAX_RETRIES",
		},
		{
			name:    "negative workflow retries",
			env:     map[string]string{"ASKOC_WORKFLOW_MAX_RETRIES": "-1"},
			wantErr: "ASKOC_WORKFLOW_MAX_RETRIES",
		},
		{
			name:    "unsupported log level",
			env:     map[string]string{"ASKOC_LOG_LEVEL": "verbose"},
			wantErr: "ASKOC_LOG_LEVEL",
		},
		{
			name:    "unsupported provider mode",
			env:     map[string]string{"ASKOC_PROVIDER": "live"},
			wantErr: "ASKOC_PROVIDER",
		},
		{
			name:    "openai-compatible missing endpoint",
			env:     map[string]string{"ASKOC_PROVIDER": "openai-compatible", "ASKOC_PROVIDER_API_KEY": "sk-demo-secret"},
			wantErr: "ASKOC_PROVIDER_ENDPOINT",
		},
		{
			name:    "openai-compatible missing api key",
			env:     map[string]string{"ASKOC_PROVIDER": "openai-compatible", "ASKOC_PROVIDER_ENDPOINT": "http://llm.local/v1/chat/completions"},
			wantErr: "ASKOC_PROVIDER_API_KEY",
		},
		{
			name:    "invalid provider timeout",
			env:     map[string]string{"ASKOC_PROVIDER_TIMEOUT_SECONDS": "soon"},
			wantErr: "ASKOC_PROVIDER_TIMEOUT_SECONDS",
		},
		{
			name:    "empty RAG chunks path",
			env:     map[string]string{"ASKOC_RAG_CHUNKS_PATH": " "},
			wantErr: "ASKOC_RAG_CHUNKS_PATH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadFromEnv(tt.env)
			if err == nil {
				t.Fatal("LoadFromEnv returned nil error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %q, want it to mention %s", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestConfigStringRedactsSecrets(t *testing.T) {
	cfg, err := LoadFromEnv(map[string]string{
		"ASKOC_AUTH_ENABLED":       "true",
		"ASKOC_AUTH_TOKEN":         "demo-token",
		"ASKOC_PROVIDER_API_KEY":   "sk-demo-secret",
		"ASKOC_WORKFLOW_URL":       "http://workflow.local/hook?sig=url-secret",
		"ASKOC_WORKFLOW_SIGNATURE": "workflow-secret",
	})
	if err != nil {
		t.Fatalf("LoadFromEnv returned error: %v", err)
	}

	got := cfg.String()
	if strings.Contains(got, "demo-token") || strings.Contains(got, "sk-demo-secret") || strings.Contains(got, "workflow-secret") || strings.Contains(got, "url-secret") {
		t.Fatalf("Config.String leaked secrets: %s", got)
	}
	if !strings.Contains(got, "auth_token:REDACTED") || !strings.Contains(got, "provider_api_key:REDACTED") || !strings.Contains(got, "workflow_url:REDACTED") || !strings.Contains(got, "workflow_signature:REDACTED") {
		t.Fatalf("Config.String did not mark redacted secrets: %s", got)
	}
}
