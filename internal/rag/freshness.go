package rag

import "time"

const HighRiskConfidenceThreshold = 0.85

type SourceDecision struct {
	Caution         bool
	Stale           bool
	HighRisk        bool
	RequiresHandoff bool
	Reason          string
}

func AssessSource(source Source, confidence float64, now time.Time) SourceDecision {
	decision := SourceDecision{
		HighRisk: source.RiskLevel == RiskHigh,
	}
	if source.FreshnessStatus == FreshnessStale {
		decision.Stale = true
	}
	if source.RequiresFreshnessCheck && source.StaleAfterDays > 0 && !source.RetrievedAt.IsZero() {
		staleAfter := source.RetrievedAt.AddDate(0, 0, source.StaleAfterDays)
		if now.After(staleAfter) {
			decision.Stale = true
		}
	}
	if decision.Stale {
		decision.Caution = true
		decision.RequiresHandoff = true
		decision.Reason = "source is stale"
		return decision
	}
	if decision.HighRisk && confidence < HighRiskConfidenceThreshold {
		decision.Caution = true
		decision.RequiresHandoff = true
		decision.Reason = "high-risk source confidence is below threshold"
	}
	return decision
}
