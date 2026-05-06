package rag

import (
	"testing"
	"time"
)

func TestAssessSourceFlagsStaleSource(t *testing.T) {
	source := Source{
		ID:                     "oc-transcript-request-2005-onwards",
		RiskLevel:              RiskMedium,
		RequiresFreshnessCheck: true,
		StaleAfterDays:         14,
		RetrievedAt:            time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		FreshnessStatus:        FreshnessFresh,
	}

	decision := AssessSource(source, 0.90, time.Date(2026, 5, 6, 0, 0, 0, 0, time.UTC))

	if !decision.Caution || !decision.Stale {
		t.Fatalf("decision = %+v, want stale caution", decision)
	}
	if !decision.RequiresHandoff {
		t.Fatalf("decision = %+v, want handoff for stale policy source", decision)
	}
}

func TestAssessSourceRequiresHigherConfidenceForHighRiskSource(t *testing.T) {
	source := Source{
		ID:                     "oc-transcript-request-2005-onwards",
		RiskLevel:              RiskHigh,
		RequiresFreshnessCheck: true,
		StaleAfterDays:         14,
		RetrievedAt:            time.Date(2026, 5, 6, 0, 0, 0, 0, time.UTC),
		FreshnessStatus:        FreshnessFresh,
	}

	lowConfidence := AssessSource(source, 0.70, time.Date(2026, 5, 6, 0, 0, 0, 0, time.UTC))
	if !lowConfidence.Caution || !lowConfidence.HighRisk || !lowConfidence.RequiresHandoff {
		t.Fatalf("low-confidence decision = %+v, want high-risk handoff", lowConfidence)
	}

	sufficientConfidence := AssessSource(source, HighRiskConfidenceThreshold, time.Date(2026, 5, 6, 0, 0, 0, 0, time.UTC))
	if sufficientConfidence.RequiresHandoff {
		t.Fatalf("sufficient-confidence decision = %+v, want no handoff", sufficientConfidence)
	}
}

func TestAssessSourceAllowsFreshMediumRiskSource(t *testing.T) {
	source := Source{
		ID:                     "oc-online-resources",
		RiskLevel:              RiskMedium,
		RequiresFreshnessCheck: true,
		StaleAfterDays:         30,
		RetrievedAt:            time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		FreshnessStatus:        FreshnessFresh,
	}

	decision := AssessSource(source, 0.70, time.Date(2026, 5, 6, 0, 0, 0, 0, time.UTC))

	if decision.Caution || decision.RequiresHandoff {
		t.Fatalf("decision = %+v, want fresh medium-risk source to pass", decision)
	}
}
