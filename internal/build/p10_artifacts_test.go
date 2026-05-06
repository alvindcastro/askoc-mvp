package build

import (
	"os"
	"strings"
	"testing"
)

func TestP10DockerfilePackagesGoServicesSafely(t *testing.T) {
	text := readRepoFile(t, "Dockerfile")

	mustContainAll(t, text,
		"FROM golang:",
		"ARG APP=api",
		"CGO_ENABLED=0",
		"go build",
		"./cmd/${APP}",
		"USER ",
		"ENTRYPOINT",
	)
	if strings.Contains(text, "COPY .env") {
		t.Fatalf("Dockerfile must not copy local .env files into images")
	}

	ignore := readRepoFile(t, ".dockerignore")
	mustContainAll(t, ignore, ".env", ".git", "coverage")
}

func TestP10ComposeStackUsesSyntheticDeterministicDefaults(t *testing.T) {
	text := readRepoFile(t, "docker-compose.yml")

	mustContainAll(t, text,
		"api:",
		"mock-banner:",
		"mock-payment:",
		"mock-crm:",
		"mock-lms:",
		"workflow-sim:",
		"ASKOC_PROVIDER: stub",
		"ASKOC_BANNER_URL: http://mock-banner:8081",
		"ASKOC_PAYMENT_URL: http://mock-payment:8082",
		"ASKOC_CRM_URL: http://mock-crm:8083",
		"ASKOC_WORKFLOW_URL: http://workflow-sim:8084/api/v1/automation/payment-reminder",
		"${ASKOC_API_PORT:-8080}:8080",
		"${ASKOC_BANNER_PORT:-8081}:8081",
		"${ASKOC_PAYMENT_PORT:-8082}:8082",
		"${ASKOC_CRM_PORT:-8083}:8083",
		"${ASKOC_WORKFLOW_PORT:-8084}:8084",
		"${ASKOC_LMS_PORT:-8085}:8085",
		"/healthz",
	)
	if strings.Contains(text, "ASKOC_PROVIDER_API_KEY:") {
		t.Fatalf("docker-compose.yml must not configure live provider secrets by default")
	}
}

func TestP10CIWorkflowRunsOfflineQualityGates(t *testing.T) {
	text := readRepoFile(t, ".github/workflows/ci.yml")

	mustContainAll(t, text,
		"go test ./...",
		"go vet ./...",
		"make eval",
		"ASKOC_PROVIDER: stub",
	)
	for _, forbidden := range []string{"ASKOC_PROVIDER_API_KEY", "OPENAI_API_KEY", "AZURE_OPENAI_API_KEY"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("CI workflow must not require live provider secret %s", forbidden)
		}
	}
}

func TestP10EnvironmentSampleAndGitignoreAreSecretSafe(t *testing.T) {
	env := readRepoFile(t, ".env.example")
	mustContainAll(t, env,
		"ASKOC_HTTP_ADDR=:8080",
		"ASKOC_PROVIDER=stub",
		"ASKOC_PROVIDER_API_KEY=",
		"ASKOC_WORKFLOW_URL=",
		"ASKOC_BANNER_URL=http://localhost:8081",
		"ASKOC_PAYMENT_URL=http://localhost:8082",
		"ASKOC_CRM_URL=http://localhost:8083",
	)
	for _, forbidden := range []string{"sk-", "xoxb-", "ghp_", "AIza"} {
		if strings.Contains(env, forbidden) {
			t.Fatalf(".env.example contains secret-looking token marker %q", forbidden)
		}
	}

	gitignore := readRepoFile(t, ".gitignore")
	mustContainAll(t, gitignore, ".env", "!.env.example")

	makefile := readRepoFile(t, "Makefile")
	mustContainAll(t, makefile, "secret-check:", "scripts/check-secrets.sh")

	script := readRepoFile(t, "scripts/check-secrets.sh")
	mustContainAll(t, script, "ghp_", "xoxb-", "AIza", "sk-live-", ".env.example")
}

func TestP10MakeTargetsAndSmokeScriptDocumentRepeatableDemo(t *testing.T) {
	makefile := readRepoFile(t, "Makefile")
	mustContainAll(t, makefile,
		"docker-build:",
		"compose-up:",
		"compose-test:",
		"smoke:",
		"scripts/smoke.sh",
	)

	script := readRepoFile(t, "scripts/smoke.sh")
	mustContainAll(t, script,
		"/healthz",
		"/api/v1/chat",
		"S100002",
		"S100003",
		"payment_reminder_triggered",
		"crm_case_created",
	)
	if !strings.HasPrefix(script, "#!/usr/bin/env sh") {
		t.Fatalf("scripts/smoke.sh must be directly executable with /usr/bin/env sh")
	}
}

func readRepoFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile("../../" + path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func mustContainAll(t *testing.T, text string, needles ...string) {
	t.Helper()
	for _, needle := range needles {
		if !strings.Contains(text, needle) {
			t.Fatalf("content missing %q", needle)
		}
	}
}
