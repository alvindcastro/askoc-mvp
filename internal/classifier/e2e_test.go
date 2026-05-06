package classifier_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/domain"
)

const (
	classificationFixturePath = "../../data/classification-fixtures.jsonl"
	minimumFixturesPerIntent  = 5
	// P6-T06 target: every synthetic fixture must keep its expected intent.
	fixtureIntentAccuracyTarget      = 1.0
	requiredNonNeutralSentimentCount = 1
)

var fixtureCoveredIntents = []domain.Intent{
	domain.IntentTranscriptRequest,
	domain.IntentTranscriptStatus,
	domain.IntentFeePayment,
	domain.IntentHumanHandoff,
	domain.IntentEscalationRequest,
	domain.IntentUnknown,
}

func TestClassificationFixturesCoverSupportedIntentsAndSentiments(t *testing.T) {
	fixtures := loadClassificationFixtures(t)
	intentCounts := map[domain.Intent]int{}
	sentimentCounts := map[domain.Sentiment]int{}

	for _, fixture := range fixtures {
		intentCounts[fixture.ExpectedIntent]++
		sentimentCounts[fixture.ExpectedSentiment]++
	}

	for _, intent := range fixtureCoveredIntents {
		if intentCounts[intent] < minimumFixturesPerIntent {
			t.Fatalf("intent %q has %d fixtures, want at least %d", intent, intentCounts[intent], minimumFixturesPerIntent)
		}
	}

	for _, sentiment := range []domain.Sentiment{domain.SentimentNegative, domain.SentimentUrgent, domain.SentimentUrgentNegative} {
		if sentimentCounts[sentiment] < requiredNonNeutralSentimentCount {
			t.Fatalf("sentiment %q has %d fixtures, want at least %d", sentiment, sentimentCounts[sentiment], requiredNonNeutralSentimentCount)
		}
	}
}

func TestFallbackClassifierMatchesClassificationFixtures(t *testing.T) {
	fixtures := loadClassificationFixtures(t)
	fallback := classifier.Fallback{}
	intentMatches := 0

	for _, fixture := range fixtures {
		fixture := fixture
		t.Run(fmt.Sprintf("%s/%s", fixture.ExpectedIntent, fixture.ID), func(t *testing.T) {
			got, err := fallback.Classify(context.Background(), fixture.Message)
			if err != nil {
				t.Fatalf("Classify returned error: %v", err)
			}

			if got.Intent == fixture.ExpectedIntent {
				intentMatches++
			} else {
				t.Fatalf("fixture %s regressed intent %q: got %q for message %q", fixture.ID, fixture.ExpectedIntent, got.Intent, fixture.Message)
			}
			if got.Sentiment != fixture.ExpectedSentiment {
				t.Fatalf("fixture %s sentiment = %q, want %q", fixture.ID, got.Sentiment, fixture.ExpectedSentiment)
			}
			if got.CanTriggerSensitiveTools() != fixture.ExpectedCanTriggerTools {
				t.Fatalf("fixture %s can trigger tools = %t, want %t", fixture.ID, got.CanTriggerSensitiveTools(), fixture.ExpectedCanTriggerTools)
			}
			if fixture.ExpectedIntent == domain.IntentUnknown && got.CanTriggerSensitiveTools() {
				t.Fatalf("unknown/off-topic fixture %s should not trigger actions: %+v", fixture.ID, got)
			}
		})
	}

	accuracy := float64(intentMatches) / float64(len(fixtures))
	if accuracy < fixtureIntentAccuracyTarget {
		t.Fatalf("fixture intent accuracy = %.2f, want %.2f", accuracy, fixtureIntentAccuracyTarget)
	}
}

type classificationFixture struct {
	ID                      string           `json:"id"`
	Message                 string           `json:"message"`
	ExpectedIntent          domain.Intent    `json:"expected_intent"`
	ExpectedSentiment       domain.Sentiment `json:"expected_sentiment"`
	ExpectedCanTriggerTools bool             `json:"expected_can_trigger_tools"`
}

func loadClassificationFixtures(t *testing.T) []classificationFixture {
	t.Helper()

	path := filepath.Clean(classificationFixturePath)
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open classification fixtures %s: %v", path, err)
	}
	defer file.Close()

	var fixtures []classificationFixture
	scanner := bufio.NewScanner(file)
	for line := 1; scanner.Scan(); line++ {
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" {
			continue
		}

		var fixture classificationFixture
		if err := json.Unmarshal([]byte(raw), &fixture); err != nil {
			t.Fatalf("decode classification fixture %s line %d: %v", path, line, err)
		}
		if fixture.ID == "" || fixture.Message == "" || fixture.ExpectedIntent == "" || fixture.ExpectedSentiment == "" {
			t.Fatalf("classification fixture %s line %d has required empty fields: %+v", path, line, fixture)
		}
		fixtures = append(fixtures, fixture)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan classification fixtures %s: %v", path, err)
	}
	if len(fixtures) == 0 {
		t.Fatalf("classification fixtures %s were empty", path)
	}

	return fixtures
}
