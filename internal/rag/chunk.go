package rag

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	DefaultMaxChunkWords     = 140
	DefaultChunkOverlapWords = 20
)

type ChunkOptions struct {
	MaxWords     int
	OverlapWords int
}

type Chunk struct {
	ID              string          `json:"id"`
	SourceID        string          `json:"source_id"`
	SourceTitle     string          `json:"source_title"`
	SourceURL       string          `json:"source_url"`
	Department      string          `json:"department,omitempty"`
	RiskLevel       RiskLevel       `json:"risk_level"`
	FreshnessStatus FreshnessStatus `json:"freshness_status"`
	ContentHash     string          `json:"content_hash,omitempty"`
	Index           int             `json:"index"`
	Text            string          `json:"text"`
}

func ChunkDocument(doc Document, opts ChunkOptions) ([]Chunk, error) {
	words := strings.Fields(normalizeText(doc.Text))
	if len(words) == 0 {
		return nil, nil
	}
	defaultedMaxWords := opts.MaxWords == 0
	if opts.MaxWords == 0 {
		opts.MaxWords = DefaultMaxChunkWords
	}
	if opts.OverlapWords == 0 && defaultedMaxWords {
		opts.OverlapWords = DefaultChunkOverlapWords
	}
	if opts.MaxWords <= 0 {
		return nil, errors.New("max words must be greater than zero")
	}
	if opts.OverlapWords < 0 {
		return nil, errors.New("overlap words must not be negative")
	}
	if opts.OverlapWords >= opts.MaxWords {
		return nil, errors.New("overlap words must be smaller than max words")
	}

	step := opts.MaxWords - opts.OverlapWords
	chunks := make([]Chunk, 0, (len(words)/step)+1)
	for start, index := 0, 0; start < len(words); start, index = start+step, index+1 {
		end := start + opts.MaxWords
		if end > len(words) {
			end = len(words)
		}
		text := strings.Join(words[start:end], " ")
		if strings.TrimSpace(text) == "" {
			continue
		}
		chunks = append(chunks, Chunk{
			ID:              stableChunkID(doc, index, text),
			SourceID:        doc.Source.ID,
			SourceTitle:     doc.Source.Title,
			SourceURL:       doc.Source.URL,
			Department:      doc.Source.Department,
			RiskLevel:       doc.Source.RiskLevel,
			FreshnessStatus: doc.Source.FreshnessStatus,
			ContentHash:     doc.ContentHash,
			Index:           index,
			Text:            text,
		})
	}
	return chunks, nil
}

func stableChunkID(doc Document, index int, text string) string {
	hash := sha1.Sum([]byte(fmt.Sprintf("%s:%s:%d:%s", doc.Source.ID, doc.ContentHash, index, text)))
	return doc.Source.ID + "-" + hex.EncodeToString(hash[:])[:12]
}

var whitespacePattern = regexp.MustCompile(`\s+`)

func normalizeText(text string) string {
	return strings.TrimSpace(whitespacePattern.ReplaceAllString(text, " "))
}

func wordCount(text string) int {
	return len(strings.Fields(text))
}
