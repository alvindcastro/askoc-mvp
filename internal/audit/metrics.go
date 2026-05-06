package audit

import (
	"sort"
	"strconv"
	"strings"
)

type DashboardMetrics struct {
	ByType          map[string]int
	ByAction        map[string]int
	ByStatus        map[string]int
	Escalations     int
	LowConfidence   int
	GuardrailEvents int
	WorkflowEvents  int
}

type Summary struct {
	TotalEvents        int               `json:"total_events"`
	TotalConversations int               `json:"total_conversations"`
	ContainmentRate    float64           `json:"containment_rate"`
	EscalationRate     float64           `json:"escalation_rate"`
	Escalations        int               `json:"escalations"`
	TopIntents         []IntentCount     `json:"top_intents"`
	Automation         AutomationMetrics `json:"automation"`
	ReviewQueue        ReviewQueue       `json:"review_queue"`
	ByType             map[string]int    `json:"by_type"`
	ByAction           map[string]int    `json:"by_action"`
	ByStatus           map[string]int    `json:"by_status"`
}

type IntentCount struct {
	Intent string `json:"intent"`
	Count  int    `json:"count"`
}

type AutomationMetrics struct {
	PaymentRemindersSent int `json:"payment_reminders_sent"`
	WorkflowFailures     int `json:"workflow_failures"`
	WorkflowEvents       int `json:"workflow_events"`
}

type ReviewQueue struct {
	LowConfidenceAnswers int          `json:"low_confidence_answers"`
	StaleSourceQuestions int          `json:"stale_source_questions"`
	Items                []ReviewItem `json:"items"`
}

type ReviewItem struct {
	TraceID        string `json:"trace_id,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	Reason         string `json:"reason"`
	Question       string `json:"question,omitempty"`
}

func DashboardMetricsFromEvents(events []Event) DashboardMetrics {
	metrics := DashboardMetrics{
		ByType:   map[string]int{},
		ByAction: map[string]int{},
		ByStatus: map[string]int{},
	}

	for _, event := range events {
		if event.Type != "" {
			metrics.ByType[event.Type]++
		}
		if event.Action != "" {
			metrics.ByAction[event.Action]++
		}
		if event.Status != "" {
			metrics.ByStatus[event.Status]++
		}
		switch event.Type {
		case EventTypeEscalation:
			metrics.Escalations++
		case EventTypeGuardrail:
			metrics.GuardrailEvents++
		case EventTypeWorkflow:
			metrics.WorkflowEvents++
		}
		if isLowConfidence(event) {
			metrics.LowConfidence++
		}
	}

	return metrics
}

func SummaryFromEvents(events []Event) Summary {
	base := DashboardMetricsFromEvents(events)
	summary := Summary{
		TotalEvents: len(events),
		Escalations: base.Escalations,
		Automation: AutomationMetrics{
			WorkflowEvents: base.WorkflowEvents,
		},
		ReviewQueue: ReviewQueue{},
		ByType:      base.ByType,
		ByAction:    base.ByAction,
		ByStatus:    base.ByStatus,
	}

	conversations := map[string]bool{}
	escalatedConversations := map[string]bool{}
	intents := map[string]int{}
	seenReviewItems := map[string]bool{}

	for _, event := range events {
		key := conversationKey(event)
		if key != "" {
			conversations[key] = true
		}
		if isEscalation(event) {
			if key != "" {
				escalatedConversations[key] = true
			}
		}
		if intent := strings.TrimSpace(event.Metadata["intent"]); intent != "" {
			intents[intent]++
		}
		if isPaymentReminderSent(event) {
			summary.Automation.PaymentRemindersSent++
		}
		if event.Type == EventTypeWorkflow && event.Status == StatusFailed {
			summary.Automation.WorkflowFailures++
		}
		if isLowConfidence(event) {
			summary.ReviewQueue.LowConfidenceAnswers++
			addReviewItem(&summary.ReviewQueue, seenReviewItems, event, "low_confidence")
		}
		if isStaleSource(event) {
			summary.ReviewQueue.StaleSourceQuestions++
			addReviewItem(&summary.ReviewQueue, seenReviewItems, event, "stale_source")
		}
	}

	summary.TotalConversations = len(conversations)
	if summary.TotalConversations > 0 {
		summary.EscalationRate = float64(len(escalatedConversations)) / float64(summary.TotalConversations)
		summary.ContainmentRate = float64(summary.TotalConversations-len(escalatedConversations)) / float64(summary.TotalConversations)
	}
	summary.TopIntents = sortedIntentCounts(intents)
	return summary
}

func isLowConfidence(event Event) bool {
	if event.Metadata == nil {
		return false
	}
	if event.Metadata["low_confidence"] == "true" {
		return true
	}
	confidence, ok := event.Metadata["confidence"]
	if !ok {
		return false
	}
	value, err := strconv.ParseFloat(confidence, 64)
	if err != nil {
		return false
	}
	return value < 0.5
}

func conversationKey(event Event) string {
	if got := strings.TrimSpace(event.ConversationID); got != "" {
		return got
	}
	return strings.TrimSpace(event.TraceID)
}

func isEscalation(event Event) bool {
	return event.Type == EventTypeEscalation || event.Action == ActionCreateCRMCase
}

func isPaymentReminderSent(event Event) bool {
	return event.Type == EventTypeWorkflow && event.Action == ActionPaymentReminder && event.Status == StatusCompleted
}

func isStaleSource(event Event) bool {
	if event.Metadata != nil && event.Metadata["stale_source"] == "true" {
		return true
	}
	return event.Type == EventTypeGuardrail && event.Action == ActionSourceCheck
}

func addReviewItem(queue *ReviewQueue, seen map[string]bool, event Event, reason string) {
	key := reason + ":" + conversationKey(event) + ":" + event.TraceID
	if seen[key] {
		return
	}
	seen[key] = true
	queue.Items = append(queue.Items, ReviewItem{
		TraceID:        event.TraceID,
		ConversationID: event.ConversationID,
		Reason:         reason,
		Question:       reviewQuestion(event),
	})
}

func reviewQuestion(event Event) string {
	if event.Metadata != nil {
		if question := strings.TrimSpace(event.Metadata["question"]); question != "" {
			return question
		}
	}
	return event.Message
}

func sortedIntentCounts(intents map[string]int) []IntentCount {
	counts := make([]IntentCount, 0, len(intents))
	for intent, count := range intents {
		counts = append(counts, IntentCount{Intent: intent, Count: count})
	}
	sort.Slice(counts, func(i, j int) bool {
		if counts[i].Count == counts[j].Count {
			return counts[i].Intent < counts[j].Intent
		}
		return counts[i].Count > counts[j].Count
	})
	return counts
}
