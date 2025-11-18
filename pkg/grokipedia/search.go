// Package grokipedia allows searching and retrieving pages from Grokipedia
package grokipedia

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SearchOption func(*searchOptions)

type searchOptions struct {
	Limit  int
	Offset int
}

func WithLimit(limit int) SearchOption {
	return func(opts *searchOptions) {
		opts.Limit = limit
	}
}

func WithOffset(offset int) SearchOption {
	return func(opts *searchOptions) {
		opts.Offset = offset
	}
}

type SearchResult struct {
	Title          string  `json:"title"`
	Slug           string  `json:"slug"`
	Snippet        string  `json:"snippet"`
	RelevanceScore float64 `json:"relevanceScore"`
}

type responseData struct {
	Results    []SearchResult `json:"results"`
	TotalCount int            `json:"total_count"`
}

func Search(ctx context.Context, query string, opts ...SearchOption) ([]SearchResult, error) {
	// Set defaults
	options := searchOptions{
		Limit:  10,
		Offset: 0,
	}
	for _, opt := range opts {
		opt(&options)
	}

	// Use Grokipedia's official API
	searchURL := fmt.Sprintf(
		"https://grokipedia.com/api/full-text-search?query=%s&limit=%d&offset=%d",
		url.QueryEscape(query),
		options.Limit,
		options.Offset,
	)

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching from Grokipedia API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: HTTP %d", resp.StatusCode)
	}

	var res responseData
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("error parsing API response: %w", err)
	}

	return res.Results, nil
}
