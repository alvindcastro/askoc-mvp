package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunExitsZeroForPassingDeterministicDataset(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "eval.jsonl")
	output := filepath.Join(dir, "summary.json")
	markdown := filepath.Join(dir, "summary.md")
	dataset := `{"id":"D01","prompt":"How do I order my official transcript?","expected_intent":"transcript_request","expected_source_contains":"oc-transcript-request-2005-onwards","expected_actions":["grounded_answer_returned"],"expected_handoff":"none","critical":true}` + "\n"
	if err := os.WriteFile(input, []byte(dataset), 0o600); err != nil {
		t.Fatalf("write dataset: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{
		"-input", input,
		"-output", output,
		"-markdown-output", markdown,
		"-chunks", "../../data/rag-chunks.json",
		"-students", "../../data/synthetic-students.json",
		"-fail-on-critical",
	}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "passed") {
		t.Fatalf("stdout = %q, want summary", stdout.String())
	}
	if _, err := os.Stat(output); err != nil {
		t.Fatalf("json output not written: %v", err)
	}
	if _, err := os.Stat(markdown); err != nil {
		t.Fatalf("markdown output not written: %v", err)
	}
}
