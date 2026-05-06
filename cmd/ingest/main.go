package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"askoc-mvp/internal/rag"
)

func main() {
	sourcePath := flag.String("sources", "data/seed-sources.json", "path to approved source allowlist")
	outPath := flag.String("out", "data/rag-chunks.json", "path to write local RAG chunks")
	maxWords := flag.Int("max-words", rag.DefaultMaxChunkWords, "maximum words per chunk")
	overlapWords := flag.Int("overlap-words", rag.DefaultChunkOverlapWords, "overlap words between chunks")
	flag.Parse()

	ctx := context.Background()
	allowlist, err := rag.LoadAllowlist(ctx, *sourcePath)
	if err != nil {
		slog.Error("load source allowlist", "error", err)
		os.Exit(1)
	}

	fetcher := rag.NewFetcher(allowlist, &http.Client{Timeout: 15 * time.Second})
	var chunks []rag.Chunk
	for _, source := range allowlist.Sources {
		doc, err := fetcher.Fetch(ctx, source.URL)
		if err != nil {
			slog.Error("fetch allowlisted source", "source_id", source.ID, "error", err)
			os.Exit(1)
		}
		sourceChunks, err := rag.ChunkDocument(doc, rag.ChunkOptions{MaxWords: *maxWords, OverlapWords: *overlapWords})
		if err != nil {
			slog.Error("chunk source", "source_id", source.ID, "error", err)
			os.Exit(1)
		}
		chunks = append(chunks, sourceChunks...)
	}

	body, err := json.MarshalIndent(chunks, "", "  ")
	if err != nil {
		slog.Error("encode chunks", "error", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*outPath, append(body, '\n'), 0o644); err != nil {
		slog.Error("write chunks", "path", *outPath, "error", err)
		os.Exit(1)
	}
	fmt.Printf("wrote %d chunks to %s\n", len(chunks), *outPath)
}
