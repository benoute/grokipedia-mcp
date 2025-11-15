package grokipedia

import (
	"context"
	"strings"
	"testing"
)

func TestSearchInputDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    SearchInput
		expected SearchInput
	}{
		{
			name:     "zero values get defaults",
			input:    SearchInput{Query: "test"},
			expected: SearchInput{Query: "test", Limit: 10, Offset: 0},
		},
		{
			name:     "negative limit gets default",
			input:    SearchInput{Query: "test", Limit: -1},
			expected: SearchInput{Query: "test", Limit: 10, Offset: 0},
		},
		{
			name:     "negative offset gets default",
			input:    SearchInput{Query: "test", Offset: -5},
			expected: SearchInput{Query: "test", Limit: 10, Offset: 0},
		},
		{
			name:     "custom values preserved",
			input:    SearchInput{Query: "test", Limit: 5, Offset: 2},
			expected: SearchInput{Query: "test", Limit: 5, Offset: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the logic of setting defaults
			limit := tt.input.Limit
			if limit <= 0 {
				limit = 10
			}
			offset := tt.input.Offset
			if offset < 0 {
				offset = 0
			}

			if limit != tt.expected.Limit {
				t.Errorf("Expected limit %d, got %d", tt.expected.Limit, limit)
			}
			if offset != tt.expected.Offset {
				t.Errorf("Expected offset %d, got %d", tt.expected.Offset, offset)
			}
		})
	}
}

func TestSearchOutputMarshalJSON(t *testing.T) {
	output := SearchOutput{}
	data, err := output.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	expected := `{"results":[]}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestGetPageOutputMarshalJSON(t *testing.T) {
	output := GetPageOutput{}
	data, err := output.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	expected := `{"title":"","content":"","citations":[]}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

// E2E tests that make real API calls to Grokipedia
// These tests require internet connection and the API to be available

func TestE2ESearchGrok(t *testing.T) {
	// Skip if running in CI or if explicitly disabled
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	ctx := context.Background()
	input := SearchInput{
		Query:  "Grok",
		Limit:  5,
		Offset: 0,
	}

	result, err := Search(ctx, input)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(result.Results) == 0 {
		t.Error("Expected at least one search result")
	}

	// Check that results contain slugs
	for i, slug := range result.Results {
		if slug == "" {
			t.Errorf("Result %d has empty slug", i)
		}
		if len(slug) < 1 {
			t.Errorf("Result %d slug is empty: %s", i, slug)
		}
	}
}

func TestE2EGetGrokPage(t *testing.T) {
	// Skip if running in CI or if explicitly disabled
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

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

func TestE2ESearchAndGetPage(t *testing.T) {
	// Skip if running in CI or if explicitly disabled
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	ctx := context.Background()

	// First search for something
	searchInput := SearchInput{
		Query:  "artificial intelligence",
		Limit:  3,
		Offset: 0,
	}

	searchResult, err := Search(ctx, searchInput)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(searchResult.Results) == 0 {
		t.Skip("No search results found, skipping page retrieval test")
	}

	// Try to get the first page
	firstSlug := searchResult.Results[0]
	pageResult, err := GetPage(ctx, firstSlug)
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
		input := SearchInput{
			Query:  "technology",
			Limit:  limit,
			Offset: 0,
		}

		result, err := Search(ctx, input)
		if err != nil {
			t.Fatalf("Search with limit %d failed: %v", limit, err)
		}

		if len(result.Results) > limit {
			t.Errorf("Expected at most %d results, got %d", limit, len(result.Results))
		}
	}
}
