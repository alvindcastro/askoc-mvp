package build

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGoTestAllCompilesFromRepoRoot(t *testing.T) {
	if os.Getenv("ASKOC_NESTED_GO_TEST") == "1" {
		t.Skip("nested go test invocation")
	}

	cmd := exec.Command("go", "test", "./...")
	cmd.Env = append(os.Environ(), "ASKOC_NESTED_GO_TEST=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go test ./... failed: %v\n%s", err, out)
	}
}

func TestMakefileDocumentsTestAndDevTargets(t *testing.T) {
	content, err := os.ReadFile("../../Makefile")
	if err != nil {
		t.Fatalf("read Makefile: %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "test:") || !strings.Contains(text, "go test ./...") {
		t.Fatalf("Makefile test target must execute go test ./...")
	}
	if !strings.Contains(text, "dev:") || !strings.Contains(text, "go run ./cmd/api") {
		t.Fatalf("Makefile dev target must run the API skeleton")
	}
}
