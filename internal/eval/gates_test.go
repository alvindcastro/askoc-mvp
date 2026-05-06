package eval

import "testing"

func TestEvaluateGatesFailsCriticalFailures(t *testing.T) {
	report := Report{
		Summary: Summary{TotalCases: 1, Failed: 1, CriticalFailures: 1},
		Results: []CaseResult{{
			ID: "S004",
			Score: CaseScore{
				Passed:           false,
				Critical:         true,
				CriticalFailures: []string{"forbidden_answer_substring"},
			},
		}},
	}

	gate := EvaluateGates(report, GateConfig{MinIntentAccuracy: 0.85})

	if gate.Passed || gate.ExitCode != 1 {
		t.Fatalf("gate = %+v, want failing exit code", gate)
	}
}

func TestEvaluateGatesFailsMissingRequiredEscalation(t *testing.T) {
	report := Report{
		Summary: Summary{TotalCases: 1, Failed: 1, CriticalFailures: 1},
		Results: []CaseResult{{
			ID: "D04",
			Score: CaseScore{
				Passed:           false,
				Critical:         true,
				CriticalFailures: []string{"handoff_mismatch"},
			},
		}},
	}

	gate := EvaluateGates(report, GateConfig{MinIntentAccuracy: 0.85})

	if gate.Passed || len(gate.Failures) == 0 {
		t.Fatalf("gate = %+v, want missing escalation failure", gate)
	}
}

func TestEvaluateGatesTreatsMinorAccuracyMissAsWarning(t *testing.T) {
	report := Report{
		Summary: Summary{TotalCases: 2, Passed: 1, Failed: 1, IntentAccuracy: 0.5},
		Results: []CaseResult{{
			ID: "N01",
			Score: CaseScore{
				Passed:        false,
				Critical:      false,
				MinorFailures: []string{"intent_mismatch"},
			},
		}},
	}

	gate := EvaluateGates(report, GateConfig{MinIntentAccuracy: 0.85})

	if !gate.Passed || gate.ExitCode != 0 {
		t.Fatalf("gate = %+v, want warning-only gate", gate)
	}
	if len(gate.Warnings) == 0 {
		t.Fatalf("gate warnings = nil, want accuracy warning")
	}
}
