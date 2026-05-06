package eval

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"askoc-mvp/internal/domain"
	"askoc-mvp/internal/privacy"
)

type Report struct {
	Summary Summary      `json:"summary"`
	Results []CaseResult `json:"results"`
	Gate    *GateResult  `json:"gate,omitempty"`
}

type Summary struct {
	TotalCases               int     `json:"total_cases"`
	Passed                   int     `json:"passed"`
	Failed                   int     `json:"failed"`
	CriticalFailures         int     `json:"critical_failures"`
	MinorFailures            int     `json:"minor_failures"`
	Warnings                 int     `json:"warnings"`
	CriticalHallucinations   int     `json:"critical_hallucinations"`
	IntentAccuracy           float64 `json:"intent_accuracy"`
	SourceRecallAt3          float64 `json:"source_recall_at_3"`
	WorkflowDecisionAccuracy float64 `json:"workflow_decision_accuracy"`
	EscalationAccuracy       float64 `json:"escalation_accuracy"`
	SafetyPassRate           float64 `json:"safety_pass_rate"`
	AverageLatencyMS         int64   `json:"average_latency_ms"`
}

type CaseResult struct {
	ID        string              `json:"id"`
	Prompt    string              `json:"prompt,omitempty"`
	Critical  bool                `json:"critical"`
	Response  domain.ChatResponse `json:"response"`
	LatencyMS int64               `json:"latency_ms"`
	Error     string              `json:"error,omitempty"`
	Score     CaseScore           `json:"score"`
}

type ReportOptions struct {
	RedactPrompts bool
}

func BuildSummary(cases []Case, results []CaseResult) Summary {
	summary := Summary{TotalCases: len(results)}
	var (
		intentExpected, intentMatched         int
		sourceExpected, sourceMatched         int
		actionExpected, actionMatched         int
		escalationExpected, escalationMatched int
		safetyExpected, safetyMatched         int
		totalLatency                          int64
	)
	for i, result := range results {
		tc := Case{}
		if i < len(cases) {
			tc = cases[i]
		}
		if result.Score.Passed {
			summary.Passed++
		} else {
			summary.Failed++
		}
		summary.CriticalFailures += len(result.Score.CriticalFailures)
		summary.MinorFailures += len(result.Score.MinorFailures)
		summary.Warnings += len(result.Score.Warnings)
		if result.Score.CriticalHallucination {
			summary.CriticalHallucinations++
		}
		totalLatency += result.LatencyMS

		if tc.ExpectedIntent != "" {
			intentExpected++
			if result.Score.IntentMatched {
				intentMatched++
			}
		}
		if tc.MustIncludeSource || strings.TrimSpace(tc.ExpectedSourceContains) != "" {
			sourceExpected++
			if result.Score.SourceMatched {
				sourceMatched++
			}
		}
		if len(tc.ExpectedActions) > 0 || len(tc.ForbiddenActions) > 0 {
			actionExpected++
			if result.Score.ActionsMatched {
				actionMatched++
			}
		}
		if tc.ExpectedEscalation != nil || strings.TrimSpace(tc.ExpectedHandoff) != "" {
			escalationExpected++
			if result.Score.HandoffMatched {
				escalationMatched++
			}
		}
		if tc.MustRefuse || tc.MustRedact != "" || tc.MustWarnPassword || tc.MustRefuseRecordAccess || len(tc.ForbiddenAnswerSubstrings) > 0 {
			safetyExpected++
			if result.Score.SafetyPassed {
				safetyMatched++
			}
		}
	}
	summary.IntentAccuracy = ratio(intentMatched, intentExpected)
	summary.SourceRecallAt3 = ratio(sourceMatched, sourceExpected)
	summary.WorkflowDecisionAccuracy = ratio(actionMatched, actionExpected)
	summary.EscalationAccuracy = ratio(escalationMatched, escalationExpected)
	summary.SafetyPassRate = ratio(safetyMatched, safetyExpected)
	if len(results) > 0 {
		summary.AverageLatencyMS = totalLatency / int64(len(results))
	}
	return summary
}

func WriteJSONReport(w io.Writer, report Report) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func RedactReportPrompts(report Report) Report {
	report.Results = append([]CaseResult(nil), report.Results...)
	for i := range report.Results {
		report.Results[i].Prompt = privacy.Redact(report.Results[i].Prompt)
	}
	return report
}

func WriteMarkdownReport(w io.Writer, report Report, opts ReportOptions) error {
	_, err := io.WriteString(w, MarkdownReport(report, opts))
	return err
}

func MarkdownReport(report Report, opts ReportOptions) string {
	var b strings.Builder
	b.WriteString("# AskOC Evaluation Summary\n\n")
	b.WriteString("| Metric | Value |\n")
	b.WriteString("|---|---:|\n")
	b.WriteString(fmt.Sprintf("| Total cases | %d |\n", report.Summary.TotalCases))
	b.WriteString(fmt.Sprintf("| Passed | %d |\n", report.Summary.Passed))
	b.WriteString(fmt.Sprintf("| Failed | %d |\n", report.Summary.Failed))
	b.WriteString(fmt.Sprintf("| Intent accuracy | %.2f |\n", report.Summary.IntentAccuracy))
	b.WriteString(fmt.Sprintf("| Source recall@3 | %.2f |\n", report.Summary.SourceRecallAt3))
	b.WriteString(fmt.Sprintf("| Workflow decision accuracy | %.2f |\n", report.Summary.WorkflowDecisionAccuracy))
	b.WriteString(fmt.Sprintf("| Escalation accuracy | %.2f |\n", report.Summary.EscalationAccuracy))
	b.WriteString(fmt.Sprintf("| Safety pass rate | %.2f |\n", report.Summary.SafetyPassRate))
	b.WriteString(fmt.Sprintf("| Critical hallucinations | %d |\n", report.Summary.CriticalHallucinations))
	b.WriteString(fmt.Sprintf("| Average latency ms | %d |\n", report.Summary.AverageLatencyMS))

	if report.Gate != nil {
		b.WriteString("\n## Quality Gate\n\n")
		status := "pass"
		if !report.Gate.Passed {
			status = "fail"
		}
		b.WriteString(fmt.Sprintf("- Status: %s\n", status))
		b.WriteString(fmt.Sprintf("- Exit code: %d\n", report.Gate.ExitCode))
		for _, failure := range report.Gate.Failures {
			b.WriteString("- Failure: " + escapeMarkdownCell(failure) + "\n")
		}
		for _, warning := range report.Gate.Warnings {
			b.WriteString("- Warning: " + escapeMarkdownCell(warning) + "\n")
		}
	}

	b.WriteString("\n## Cases\n\n")
	b.WriteString("| ID | Result | Critical | Prompt | Failures |\n")
	b.WriteString("|---|---|---:|---|---|\n")
	for _, result := range report.Results {
		status := "pass"
		if !result.Score.Passed {
			status = "fail"
		}
		prompt := result.Prompt
		if opts.RedactPrompts {
			prompt = privacy.Redact(prompt)
		}
		failures := append([]string{}, result.Score.CriticalFailures...)
		failures = append(failures, result.Score.MinorFailures...)
		b.WriteString(fmt.Sprintf("| %s | %s | %t | %s | %s |\n",
			escapeMarkdownCell(result.ID),
			status,
			result.Critical,
			escapeMarkdownCell(prompt),
			escapeMarkdownCell(strings.Join(failures, "; ")),
		))
	}
	return b.String()
}

func ratio(numerator, denominator int) float64 {
	if denominator == 0 {
		return 1
	}
	return float64(numerator) / float64(denominator)
}

func escapeMarkdownCell(value string) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "|", `\|`)
	return value
}
