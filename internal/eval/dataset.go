package eval

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"askoc-mvp/internal/domain"
)

type Case struct {
	ID                        string           `json:"id"`
	Prompt                    string           `json:"prompt"`
	Channel                   string           `json:"channel,omitempty"`
	StudentID                 string           `json:"student_id,omitempty"`
	ExpectedIntent            domain.Intent    `json:"expected_intent"`
	ExpectedSentiment         domain.Sentiment `json:"expected_sentiment,omitempty"`
	MustIncludeSource         bool             `json:"must_include_source,omitempty"`
	ExpectedSourceContains    string           `json:"expected_source_contains,omitempty"`
	ExpectedActions           []string         `json:"expected_actions,omitempty"`
	ForbiddenActions          []string         `json:"forbidden_actions,omitempty"`
	ExpectedEscalation        *bool            `json:"expected_escalation,omitempty"`
	ExpectedHandoff           string           `json:"expected_handoff,omitempty"`
	MustRefuse                bool             `json:"must_refuse,omitempty"`
	MustRedact                string           `json:"must_redact,omitempty"`
	MustWarnPassword          bool             `json:"must_warn_password,omitempty"`
	MustRefuseRecordAccess    bool             `json:"must_refuse_record_access,omitempty"`
	ForbiddenAnswerSubstrings []string         `json:"forbidden_answer_substrings,omitempty"`
	MaxLatencyMS              int              `json:"max_latency_ms,omitempty"`
	Critical                  bool             `json:"critical,omitempty"`
	Tags                      []string         `json:"tags,omitempty"`
}

func LoadDatasetFile(ctx context.Context, path string) ([]Case, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open evaluation dataset: %w", err)
	}
	defer file.Close()
	return ParseDataset(ctx, file)
}

func ParseDataset(ctx context.Context, r io.Reader) ([]Case, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var cases []Case
	seenIDs := map[string]bool{}
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var tc Case
		decoder := json.NewDecoder(strings.NewReader(line))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&tc); err != nil {
			return nil, fmt.Errorf("line %d: decode evaluation case: %w", lineNo, err)
		}
		if err := validateCase(tc); err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNo, err)
		}
		if seenIDs[tc.ID] {
			return nil, fmt.Errorf("line %d: duplicate id %q", lineNo, tc.ID)
		}
		seenIDs[tc.ID] = true
		cases = append(cases, tc)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read evaluation dataset: %w", err)
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("evaluation dataset must include at least one case")
	}
	return cases, nil
}

func validateCase(tc Case) error {
	if strings.TrimSpace(tc.ID) == "" {
		return fmt.Errorf("id is required")
	}
	if strings.TrimSpace(tc.Prompt) == "" {
		return fmt.Errorf("prompt is required")
	}
	if !validIntent(tc.ExpectedIntent) {
		return fmt.Errorf("expected_intent is required and must be supported")
	}
	if tc.ExpectedSentiment != "" && !validSentiment(tc.ExpectedSentiment) {
		return fmt.Errorf("expected_sentiment %q is not supported", tc.ExpectedSentiment)
	}
	if tc.MaxLatencyMS < 0 {
		return fmt.Errorf("max_latency_ms cannot be negative")
	}
	if !hasExpectedBehavior(tc) {
		return fmt.Errorf("case must include expected action, source, handoff, escalation, or safety behavior")
	}
	return nil
}

func hasExpectedBehavior(tc Case) bool {
	return len(tc.ExpectedActions) > 0 ||
		len(tc.ForbiddenActions) > 0 ||
		strings.TrimSpace(tc.ExpectedSourceContains) != "" ||
		strings.TrimSpace(tc.ExpectedHandoff) != "" ||
		tc.MustIncludeSource ||
		tc.ExpectedEscalation != nil ||
		tc.MustRefuse ||
		strings.TrimSpace(tc.MustRedact) != "" ||
		tc.MustWarnPassword ||
		tc.MustRefuseRecordAccess ||
		len(tc.ForbiddenAnswerSubstrings) > 0
}

func validIntent(intent domain.Intent) bool {
	switch intent {
	case domain.IntentTranscriptRequest,
		domain.IntentTranscriptStatus,
		domain.IntentFeePayment,
		domain.IntentHumanHandoff,
		domain.IntentEscalationRequest,
		domain.IntentUnknown:
		return true
	default:
		return false
	}
}

func validSentiment(sentiment domain.Sentiment) bool {
	switch sentiment {
	case domain.SentimentNeutral,
		domain.SentimentNegative,
		domain.SentimentUrgent,
		domain.SentimentUrgentNegative:
		return true
	default:
		return false
	}
}
