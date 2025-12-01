package grokipedia

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type citation struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type image struct {
	URL         string `json:"url"`
	Caption     string `json:"caption"`
	Description string `json:"description"`
}

type metadata struct {
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Views     int    `json:"views"`
}

type Page struct {
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Citations []citation `json:"citations"`
	Images    []image    `json:"images"`
	Metadata  metadata   `json:"metadata"`
	Slug      string     `json:"slug"`
}

var ErrNotFound = errors.New("page not found")

type pageOptions struct {
	includeContent bool
}

var defaultPageOptions = pageOptions{
	includeContent: true,
}

type PageOption func(*pageOptions)

func WithoutContent() PageOption {
	return func(opts *pageOptions) {
		opts.includeContent = false
	}
}

func GetPage(ctx context.Context, slug string, opts ...PageOption) (*Page, error) {
	// Default options
	options := defaultPageOptions

	// Apply options
	for _, opt := range opts {
		opt(&options)
	}

	// Use Grokipedia's official API to get page content
	pageURL := fmt.Sprintf(
		"https://grokipedia.com/api/page?slug=%s&includeContent=%t&validateLinks=true",
		url.QueryEscape(slug),
		options.includeContent,
	)

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Make request
	req, err := http.NewRequestWithContext(timeoutCtx, "GET", pageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching page from Grokipedia API: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 status
	if resp.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(resp.Body)
		if err != nil && len(msg) > 0 {
			return nil, fmt.Errorf("API error: HTTP %d: %s", resp.StatusCode, msg)
		}
		return nil, fmt.Errorf("API error: HTTP %d", resp.StatusCode)
	}

	// Decode JSON response
	var res struct {
		Page  Page `json:"page"`
		Found bool `json:"found"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("error parsing API response: %w", err)
	}

	if !res.Found {
		return nil, ErrNotFound
	}

	return &res.Page, nil
}
