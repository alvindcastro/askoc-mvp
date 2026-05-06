package main

import (
	"net/http"
	"testing"
	"time"

	"askoc-mvp/internal/audit"
	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/config"
	"askoc-mvp/internal/orchestrator"
	"askoc-mvp/internal/workflow"
)

func TestBuildLLMUsesDisabledLLMForDefaultProvider(t *testing.T) {
	cfg := config.Config{
		Provider: config.ProviderConfig{
			Mode:    "stub",
			Model:   "demo-placeholder",
			Timeout: time.Second,
		},
	}

	got, err := buildLLM(cfg, http.DefaultClient)
	if err != nil {
		t.Fatalf("buildLLM returned error: %v", err)
	}
	if _, ok := got.(orchestrator.DisabledLLM); !ok {
		t.Fatalf("buildLLM returned %T, want DisabledLLM for stub provider", got)
	}
}

func TestBuildLLMRequiresConfiguredProviderValuesForOpenAICompatibleMode(t *testing.T) {
	cfg := config.Config{
		Provider: config.ProviderConfig{
			Mode:    "openai-compatible",
			Model:   "gpt-demo",
			Timeout: time.Second,
		},
	}

	_, err := buildLLM(cfg, http.DefaultClient)
	if err == nil {
		t.Fatal("buildLLM returned nil error for missing endpoint/API key")
	}
}

func TestBuildClassifierUsesGuardedLLMClassifierOnlyForConfiguredProvider(t *testing.T) {
	stubCfg := config.Config{Provider: config.ProviderConfig{Mode: "stub"}}
	if _, ok := buildClassifier(stubCfg, orchestrator.DisabledLLM{}, audit.NopRecorder{}).(classifier.Fallback); !ok {
		t.Fatalf("stub classifier should use deterministic fallback")
	}

	openAICfg := config.Config{Provider: config.ProviderConfig{Mode: "openai-compatible"}}
	if _, ok := buildClassifier(openAICfg, orchestrator.DisabledLLM{}, audit.NopRecorder{}).(orchestrator.LLMBackedClassifier); !ok {
		t.Fatalf("openai-compatible classifier should use LLMBackedClassifier")
	}
}

func TestBuildWorkflowUsesInMemoryClientWhenWebhookURLIsMissing(t *testing.T) {
	got, err := buildWorkflow(config.Config{}, http.DefaultClient)
	if err != nil {
		t.Fatalf("buildWorkflow returned error: %v", err)
	}
	if _, ok := got.(*workflow.InMemoryClient); !ok {
		t.Fatalf("buildWorkflow returned %T, want in-memory workflow client", got)
	}
}

func TestBuildWorkflowUsesPowerAutomateClientWhenWebhookURLIsConfigured(t *testing.T) {
	got, err := buildWorkflow(config.Config{
		Workflow: config.WorkflowConfig{
			URL:             "http://workflow.local/hook",
			Signature:       "workflow-secret",
			SignatureHeader: "X-Demo-Signature",
			MaxRetries:      2,
		},
	}, http.DefaultClient)
	if err != nil {
		t.Fatalf("buildWorkflow returned error: %v", err)
	}
	if _, ok := got.(*workflow.PowerAutomateClient); !ok {
		t.Fatalf("buildWorkflow returned %T, want Power Automate workflow client", got)
	}
}

func TestAdminAccessTokenUsesConfiguredAuthTokenOrDemoFallback(t *testing.T) {
	if got := adminAccessToken(config.Config{}); got != "demo-admin-token" {
		t.Fatalf("adminAccessToken default = %q, want demo-admin-token", got)
	}
	cfg := config.Config{Auth: config.AuthConfig{Token: "configured-admin"}}
	if got := adminAccessToken(cfg); got != "configured-admin" {
		t.Fatalf("adminAccessToken override = %q, want configured-admin", got)
	}
}
