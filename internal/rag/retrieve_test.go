package rag

import (
	"context"
	"testing"
)

func TestLocalRetrieverRanksTranscriptChunksFirst(t *testing.T) {
	retriever := NewLocalRetriever([]Chunk{
		chunk("lms-1", "oc-online-resources", "Online Resources", "Find online course access and myOkanagan resources.", RiskMedium),
		chunk("transcript-1", "oc-transcript-request-2005-onwards", "Transcript Request Guidance", "Order an official transcript request through the public Registrar guidance.", RiskHigh),
	})

	results, err := retriever.Retrieve(context.Background(), "How do I order my official transcript?", 3)
	if err != nil {
		t.Fatalf("Retrieve returned error: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("Retrieve returned no results")
	}
	if results[0].Chunk.SourceID != "oc-transcript-request-2005-onwards" {
		t.Fatalf("top source = %q, want transcript source; results=%+v", results[0].Chunk.SourceID, results)
	}
	if results[0].Confidence < DefaultMinimumConfidence {
		t.Fatalf("confidence = %.2f, want >= %.2f", results[0].Confidence, DefaultMinimumConfidence)
	}
}

func TestLocalRetrieverUnrelatedQueryReturnsLowConfidence(t *testing.T) {
	retriever := NewLocalRetriever([]Chunk{
		chunk("transcript-1", "oc-transcript-request-2005-onwards", "Transcript Request Guidance", "Order an official transcript request through the public Registrar guidance.", RiskHigh),
	})

	results, err := retriever.Retrieve(context.Background(), "Where can I park a bicycle downtown?", 3)
	if err != nil {
		t.Fatalf("Retrieve returned error: %v", err)
	}

	if len(results) > 0 && results[0].Confidence >= DefaultMinimumConfidence {
		t.Fatalf("top confidence = %.2f, want low confidence for unrelated query", results[0].Confidence)
	}
	sources, err := retriever.RetrieveSources(context.Background(), "Where can I park a bicycle downtown?")
	if err != nil {
		t.Fatalf("RetrieveSources returned error: %v", err)
	}
	if len(sources) != 0 {
		t.Fatalf("sources = %+v, want none for unrelated query", sources)
	}
}

func TestLocalRetrieverRespectsLimitAndReturnsSourceMetadata(t *testing.T) {
	retriever := NewLocalRetriever([]Chunk{
		chunk("transcript-1", "oc-transcript-request-2005-onwards", "Transcript Request Guidance", "Order an official transcript request.", RiskHigh),
		chunk("transcript-2", "oc-registrar-office", "Office of the Registrar", "Registrar transcript services and public learner service context.", RiskMedium),
		chunk("transcript-3", "oc-online-resources", "Online Resources", "Transcript-related online resource navigation.", RiskMedium),
	})

	results, err := retriever.Retrieve(context.Background(), "transcript registrar request", 2)
	if err != nil {
		t.Fatalf("Retrieve returned error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("results = %d, want limit of 2", len(results))
	}
	for _, result := range results {
		if result.Chunk.SourceTitle == "" || result.Chunk.SourceURL == "" || result.Chunk.RiskLevel == "" {
			t.Fatalf("result missing source metadata: %+v", result)
		}
	}
}

func chunk(id, sourceID, title, text string, risk RiskLevel) Chunk {
	return Chunk{
		ID:              id,
		SourceID:        sourceID,
		SourceTitle:     title,
		SourceURL:       "https://www.okanagancollege.ca/" + sourceID,
		RiskLevel:       risk,
		FreshnessStatus: FreshnessFresh,
		Text:            text,
	}
}
