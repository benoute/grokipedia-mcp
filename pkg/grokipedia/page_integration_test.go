package grokipedia

import (
	"context"
	"strings"
	"testing"
)

func TestIntegrationGetPage(t *testing.T) {
	ctx := context.Background()
	result, err := GetPage(ctx, "Grok")
	if err != nil {
		t.Fatalf("GetPage failed: %v", err)
	}

	if result.Title == "" {
		t.Error("Expected non-empty title")
	}

	if result.Content == "" {
		t.Error("Expected non-empty content")
	}

	// Grok page should mention xAI or Elon Musk
	contentLower := strings.ToLower(result.Content)
	if !strings.Contains(contentLower, "xai") &&
		!strings.Contains(contentLower, "elon") &&
		!strings.Contains(contentLower, "musk") {
		t.Logf(
			"Warning: Grok page content doesn't mention xAI or Elon Musk. Content preview: %s...",
			result.Content[:200],
		)
	}

	// Check that citations are properly handled (even if empty)
	if result.Citations == nil {
		t.Error("Citations should not be nil")
	}
}
