package rag

import (
	"context"

	"askoc-mvp/internal/domain"
)

const DefaultMinimumConfidence = 0.35

type RetrievalResult struct {
	Chunk      Chunk
	Score      int
	Confidence float64
}

type Retriever interface {
	Retrieve(ctx context.Context, query string, limit int) ([]RetrievalResult, error)
	RetrieveSources(ctx context.Context, query string) ([]domain.Source, error)
}
