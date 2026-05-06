package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"askoc-mvp/internal/domain"
)

type LocalRetriever struct {
	chunks []Chunk
}

func NewLocalRetriever(chunks []Chunk) *LocalRetriever {
	copied := append([]Chunk(nil), chunks...)
	return &LocalRetriever{chunks: copied}
}

func LoadChunks(ctx context.Context, path string) ([]Chunk, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load RAG chunks: %w", err)
	}
	var chunks []Chunk
	if err := json.Unmarshal(body, &chunks); err != nil {
		return nil, fmt.Errorf("parse RAG chunks: %w", err)
	}
	return chunks, nil
}

func (r *LocalRetriever) Retrieve(ctx context.Context, query string, limit int) ([]RetrievalResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 3
	}
	queryTerms := tokenize(query)
	if len(queryTerms) == 0 {
		return nil, nil
	}
	results := make([]RetrievalResult, 0, len(r.chunks))
	for _, chunk := range r.chunks {
		score := scoreChunk(queryTerms, chunk)
		if score == 0 {
			continue
		}
		confidence := float64(score) / float64(len(queryTerms))
		if confidence > 1 {
			confidence = 1
		}
		results = append(results, RetrievalResult{
			Chunk:      chunk,
			Score:      score,
			Confidence: confidence,
		})
	}
	sort.SliceStable(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Chunk.ID < results[j].Chunk.ID
		}
		return results[i].Score > results[j].Score
	})
	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

func (r *LocalRetriever) RetrieveSources(ctx context.Context, query string) ([]domain.Source, error) {
	results, err := r.Retrieve(ctx, query, 3)
	if err != nil {
		return nil, err
	}
	sources := make([]domain.Source, 0, len(results))
	seen := map[string]bool{}
	for _, result := range results {
		if result.Confidence < DefaultMinimumConfidence {
			continue
		}
		key := normalizeURL(result.Chunk.SourceURL)
		if seen[key] {
			continue
		}
		seen[key] = true
		sources = append(sources, domain.Source{
			ID:              result.Chunk.SourceID,
			Title:           result.Chunk.SourceTitle,
			URL:             result.Chunk.SourceURL,
			ChunkID:         result.Chunk.ID,
			Confidence:      result.Confidence,
			RiskLevel:       string(result.Chunk.RiskLevel),
			FreshnessStatus: string(result.Chunk.FreshnessStatus),
		})
	}
	return sources, nil
}

var tokenPattern = regexp.MustCompile(`[a-z0-9]+`)

func tokenize(text string) []string {
	raw := tokenPattern.FindAllString(strings.ToLower(text), -1)
	terms := make([]string, 0, len(raw))
	seen := map[string]bool{}
	for _, term := range raw {
		if stopWords[term] || len(term) < 2 {
			continue
		}
		if seen[term] {
			continue
		}
		seen[term] = true
		terms = append(terms, term)
	}
	return terms
}

func scoreChunk(queryTerms []string, chunk Chunk) int {
	searchable := strings.ToLower(chunk.SourceTitle + " " + chunk.SourceID + " " + chunk.Text)
	score := 0
	for _, term := range queryTerms {
		if strings.Contains(searchable, term) {
			score++
		}
		if strings.Contains(strings.ToLower(chunk.SourceTitle), term) {
			score++
		}
	}
	return score
}

var stopWords = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "can": true, "do": true,
	"for": true, "how": true, "i": true, "is": true, "it": true, "my": true,
	"of": true, "on": true, "or": true, "the": true, "to": true, "where": true,
	"with": true, "you": true,
}
