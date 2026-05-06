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
	if cfg.Provider.Mode != "stub" {
		t.Fatalf("Provider.Mode = %q, want stub", cfg.Provider.Mode)
	}
	if cfg.Provider.Model != "demo-placeholder" {
		t.Fatalf("Provider.Model = %q, want demo-placeholder", cfg.Provider.Model)
	}
}

func TestLoadFromEnvUsesOverrides(t *testing.T) {
	cfg, err := LoadFromEnv(map[string]string{
		"ASKOC_HTTP_ADDR":                "127.0.0.1:9090",
		"ASKOC_AUTH_ENABLED":             "true",
		"ASKOC_AUTH_TOKEN":               "demo-token",
		"ASKOC_LOG_LEVEL":                "debug",
		"ASKOC_WORKFLOW_URL":             "http://workflow.local/hook",
		"ASKOC_WORKFLOW_TIMEOUT_SECONDS": "12",
		"ASKOC_PROVIDER":                 "openai-compatible",
		"ASKOC_PROVIDER_MODEL":           "gpt-demo",
		"ASKOC_PROVIDER_API_KEY":         "sk-demo-secret",
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
	if cfg.Provider.Mode != "openai-compatible" || cfg.Provider.Model != "gpt-demo" {
		t.Fatalf("Provider = %+v", cfg.Provider)
	}
	if cfg.Provider.APIKey != "sk-demo-secret" {
		t.Fatalf("Provider.APIKey was not loaded")
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
			name:    "unsupported log level",
			env:     map[string]string{"ASKOC_LOG_LEVEL": "verbose"},
			wantErr: "ASKOC_LOG_LEVEL",
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
		"ASKOC_AUTH_ENABLED":     "true",
		"ASKOC_AUTH_TOKEN":       "demo-token",
		"ASKOC_PROVIDER_API_KEY": "sk-demo-secret",
	})
	if err != nil {
		t.Fatalf("LoadFromEnv returned error: %v", err)
	}

	got := cfg.String()
	if strings.Contains(got, "demo-token") || strings.Contains(got, "sk-demo-secret") {
		t.Fatalf("Config.String leaked secrets: %s", got)
	}
	if !strings.Contains(got, "auth_token:REDACTED") || !strings.Contains(got, "provider_api_key:REDACTED") {
		t.Fatalf("Config.String did not mark redacted secrets: %s", got)
	}
}
