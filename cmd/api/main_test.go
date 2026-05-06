package main

import (
	"net/http"
	"testing"
	"time"

	"askoc-mvp/internal/audit"
	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/config"
	"askoc-mvp/internal/orchestrator"
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
