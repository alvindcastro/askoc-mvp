package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	evalrunner "askoc-mvp/internal/eval"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("eval", flag.ContinueOnError)
	fs.SetOutput(stderr)
	input := fs.String("input", "data/eval-questions.jsonl", "JSONL evaluation dataset path")
	output := fs.String("output", "reports/eval-summary.json", "JSON evaluation report output path")
	markdownOutput := fs.String("markdown-output", "reports/eval-summary.md", "Markdown evaluation report output path")
	baseURL := fs.String("base-url", "", "optional live API base URL, for example http://localhost:8080/api/v1")
	chunksPath := fs.String("chunks", "data/rag-chunks.json", "local RAG chunks path for deterministic mode")
	studentsPath := fs.String("students", "data/synthetic-students.json", "synthetic student fixture path for deterministic mode")
	timeout := fs.Duration("timeout", 5*time.Second, "per-case timeout")
	failOnCritical := fs.Bool("fail-on-critical", true, "return non-zero when critical evaluation gates fail")
	redactPrompts := fs.Bool("redact-prompts", true, "redact prompts in Markdown report")
	minIntentAccuracy := fs.Float64("min-intent-accuracy", 0.85, "intent accuracy warning threshold")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	ctx := context.Background()
	cases, err := evalrunner.LoadDatasetFile(ctx, *input)
	if err != nil {
		fmt.Fprintf(stderr, "load dataset: %v\n", err)
		return 2
	}

	client, err := buildClient(ctx, *baseURL, *timeout, *chunksPath, *studentsPath)
	if err != nil {
		fmt.Fprintf(stderr, "create eval client: %v\n", err)
		return 2
	}

	report, err := evalrunner.Runner{Client: client, Timeout: *timeout}.Run(ctx, cases)
	if err != nil {
		fmt.Fprintf(stderr, "run eval: %v\n", err)
		return 2
	}
	gate := evalrunner.EvaluateGates(report, evalrunner.GateConfig{MinIntentAccuracy: *minIntentAccuracy})
	report.Gate = &gate

	outputReport := report
	if *redactPrompts {
		outputReport = evalrunner.RedactReportPrompts(report)
	}
	if err := writeJSON(*output, outputReport); err != nil {
		fmt.Fprintf(stderr, "write JSON report: %v\n", err)
		return 2
	}
	if err := writeMarkdown(*markdownOutput, outputReport, evalrunner.ReportOptions{RedactPrompts: *redactPrompts}); err != nil {
		fmt.Fprintf(stderr, "write Markdown report: %v\n", err)
		return 2
	}

	fmt.Fprintf(stdout, "passed %d/%d cases; critical failures %d; report %s\n",
		report.Summary.Passed,
		report.Summary.TotalCases,
		report.Summary.CriticalFailures,
		*output,
	)
	for _, warning := range gate.Warnings {
		fmt.Fprintf(stdout, "warning: %s\n", warning)
	}
	if *failOnCritical && !gate.Passed {
		for _, failure := range gate.Failures {
			fmt.Fprintf(stderr, "critical: %s\n", failure)
		}
		return gate.ExitCode
	}
	return 0
}

func buildClient(ctx context.Context, baseURL string, timeout time.Duration, chunksPath, studentsPath string) (evalrunner.ChatClient, error) {
	if baseURL != "" {
		return evalrunner.HTTPChatClient{
			BaseURL:    baseURL,
			HTTPClient: &http.Client{Timeout: timeout},
		}, nil
	}
	return evalrunner.NewDeterministicChatClient(ctx, evalrunner.DeterministicClientConfig{
		ChunksPath:   chunksPath,
		StudentsPath: studentsPath,
	})
}

func writeJSON(path string, report evalrunner.Report) error {
	if path == "-" {
		return evalrunner.WriteJSONReport(os.Stdout, report)
	}
	if err := ensureParent(path); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return evalrunner.WriteJSONReport(file, report)
}

func writeMarkdown(path string, report evalrunner.Report, opts evalrunner.ReportOptions) error {
	if path == "" {
		return nil
	}
	if path == "-" {
		return evalrunner.WriteMarkdownReport(os.Stdout, report, opts)
	}
	if err := ensureParent(path); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return evalrunner.WriteMarkdownReport(file, report, opts)
}

func ensureParent(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}
