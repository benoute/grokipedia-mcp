package grokipedia

import (
	"context"
	"testing"
)

func TestIntegrationSearch(t *testing.T) {
	ctx := context.Background()

	result, err := Search(ctx, "Grok", WithLimit(5), WithOffset(0))

	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected at least one search result")
	}

	if len(result) > 5 {
		t.Error("Expected at max 5 results")
	}

	// Check that results contain slugs
	for i, res := range result {
		if res.Slug == "" {
			t.Errorf("Result %d has empty slug", i)
		}
		if len(res.Slug) < 1 {
			t.Errorf("Result %d slug is empty: %s", i, res.Slug)
		}
	}
}
