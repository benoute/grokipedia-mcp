package grokipedia_test

import (
	"context"
	"strings"
	"testing"

	"github.com/benoute/grokipedia/pkg/grokipedia"
)

func TestE2EGetGrokPage(t *testing.T) {
	// Skip if running in CI or if explicitly disabled
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	ctx := context.Background()
	result, err := grokipedia.GetPage(ctx, "Grok")
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

func TestE2ESearchGrok(t *testing.T) {
	// Skip if running in CI or if explicitly disabled
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	ctx := context.Background()

	result, err := grokipedia.Search(ctx, "Grok", grokipedia.WithLimit(5), grokipedia.WithOffset(0))
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected at least one search result")
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

func TestE2ESearchAndGetPage(t *testing.T) {
	// Skip if running in CI or if explicitly disabled
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	ctx := context.Background()

	// First search for something
	searchResult, err := grokipedia.Search(ctx, "artificial intelligence", grokipedia.WithLimit(3), grokipedia.WithOffset(0))
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(searchResult) == 0 {
		t.Skip("No search results found, skipping page retrieval test")
	}

	// Try to get the first page
	firstSlug := searchResult[0].Slug
	pageResult, err := grokipedia.GetPage(ctx, firstSlug)
	if err != nil {
		t.Fatalf("GetPage failed for slug %s: %v", firstSlug, err)
	}

	if pageResult.Title == "" {
		t.Errorf("Page title is empty for slug %s", firstSlug)
	}

	if pageResult.Content == "" {
		t.Errorf("Page content is empty for slug %s", firstSlug)
	}

	t.Logf("Successfully retrieved page: %s", pageResult.Title)
}

func TestE2ESearchWithLimitOffset(t *testing.T) {
	// Skip if running in CI or if explicitly disabled
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	ctx := context.Background()

	// Test with different limits
	limits := []int{1, 3, 5}

	for _, limit := range limits {
		result, err := grokipedia.Search(ctx, "technology", grokipedia.WithLimit(limit), grokipedia.WithOffset(0))
		if err != nil {
			t.Fatalf("Search with limit %d failed: %v", limit, err)
		}

		if len(result) > limit {
			t.Errorf("Expected at most %d results, got %d", limit, len(result))
		}
	}
}
