package eval

import "fmt"

type GateConfig struct {
	MinIntentAccuracy float64
}

type GateResult struct {
	Passed   bool     `json:"passed"`
	ExitCode int      `json:"exit_code"`
	Failures []string `json:"failures,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

func EvaluateGates(report Report, cfg GateConfig) GateResult {
	gate := GateResult{Passed: true}
	for _, result := range report.Results {
		for _, failure := range result.Score.CriticalFailures {
			gate.Failures = append(gate.Failures, fmt.Sprintf("%s: %s", result.ID, failure))
		}
	}
	if report.Summary.CriticalHallucinations > 0 {
		gate.Failures = append(gate.Failures, fmt.Sprintf("critical_hallucinations: %d", report.Summary.CriticalHallucinations))
	}
	if cfg.MinIntentAccuracy > 0 && report.Summary.IntentAccuracy < cfg.MinIntentAccuracy {
		gate.Warnings = append(gate.Warnings, fmt.Sprintf("intent_accuracy_below_target: got %.2f want %.2f", report.Summary.IntentAccuracy, cfg.MinIntentAccuracy))
	}
	if len(gate.Failures) > 0 {
		gate.Passed = false
		gate.ExitCode = 1
	}
	return gate
}
