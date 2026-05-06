package rag

import "testing"

func TestChunkDocumentShortContentProducesOneChunk(t *testing.T) {
	doc := testDocument("short transcript guidance")

	chunks, err := ChunkDocument(doc, ChunkOptions{MaxWords: 20})
	if err != nil {
		t.Fatalf("ChunkDocument returned error: %v", err)
	}

	if len(chunks) != 1 {
		t.Fatalf("chunks = %d, want 1", len(chunks))
	}
	if chunks[0].Text != "short transcript guidance" {
		t.Fatalf("chunk text = %q", chunks[0].Text)
	}
	if chunks[0].ID == "" {
		t.Fatal("chunk ID was empty")
	}
}

func TestChunkDocumentLongContentProducesBoundedStableChunks(t *testing.T) {
	doc := testDocument("one two three four five six seven eight nine ten eleven twelve thirteen fourteen fifteen")

	first, err := ChunkDocument(doc, ChunkOptions{MaxWords: 5, OverlapWords: 1})
	if err != nil {
		t.Fatalf("ChunkDocument returned error: %v", err)
	}
	second, err := ChunkDocument(doc, ChunkOptions{MaxWords: 5, OverlapWords: 1})
	if err != nil {
		t.Fatalf("second ChunkDocument returned error: %v", err)
	}

	if len(first) < 3 {
		t.Fatalf("chunks = %d, want multiple bounded chunks", len(first))
	}
	if len(first) != len(second) {
		t.Fatalf("stable run chunk count mismatch: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if wordCount(first[i].Text) > 5 {
			t.Fatalf("chunk %d has %d words, want <= 5: %q", i, wordCount(first[i].Text), first[i].Text)
		}
		if first[i].ID != second[i].ID {
			t.Fatalf("chunk %d ID = %q then %q, want stable IDs", i, first[i].ID, second[i].ID)
		}
		if first[i].Text == "" {
			t.Fatalf("chunk %d text was empty", i)
		}
	}
}

func TestChunkDocumentCopiesMetadataToEveryChunk(t *testing.T) {
	doc := testDocument("transcript payment registrar guidance needs chunk metadata copied")

	chunks, err := ChunkDocument(doc, ChunkOptions{MaxWords: 3})
	if err != nil {
		t.Fatalf("ChunkDocument returned error: %v", err)
	}

	for _, chunk := range chunks {
		if chunk.SourceID != doc.Source.ID || chunk.SourceURL != doc.Source.URL || chunk.SourceTitle != doc.Source.Title {
			t.Fatalf("chunk metadata = %+v, want document source metadata", chunk)
		}
		if chunk.RiskLevel != doc.Source.RiskLevel || chunk.FreshnessStatus != doc.Source.FreshnessStatus {
			t.Fatalf("chunk risk/freshness = %+v, want source metadata", chunk)
		}
	}
}

func testDocument(text string) Document {
	return Document{
		Source: Source{
			ID:              "oc-transcript-request-2005-onwards",
			Title:           "Transcript Request Guidance",
			URL:             "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
			Department:      "Registrar",
			RiskLevel:       RiskHigh,
			FreshnessStatus: FreshnessFresh,
		},
		Title:       "Transcript Request Guidance",
		URL:         "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
		Text:        text,
		ContentHash: "hash-test",
	}
}
