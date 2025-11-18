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

type pageOptions struct {
	includeContent bool
}

type PageOption func(*pageOptions)

func WithoutContent() PageOption {
	return func(opts *pageOptions) {
		opts.includeContent = false
	}
}

type GetPageOutput struct {
	Title     string     `json:"title" jsonschema:"The page title"`
	Content   string     `json:"content" jsonschema:"The full page content"`
	Citations []citation `json:"citations" jsonschema:"List of citations"`
}

type ResponseData struct {
	Page  Page `json:"page"`
	Found bool `json:"found"`
}

type Page struct {
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Citations []citation `json:"citations"`
	Images    []image    `json:"images"`
	Metadata  metadata   `json:"metadata"`
	Slug      string     `json:"slug"`
}

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

var ErrNotFound = errors.New("page not found")

func GetPage(ctx context.Context, slug string, opts ...PageOption) (*Page, error) {
	options := pageOptions{
		includeContent: true,
	}

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

	req, err := http.NewRequestWithContext(timeoutCtx, "GET", pageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching page from Grokipedia API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(resp.Body)
		if err != nil && len(msg) > 0 {
			return nil, fmt.Errorf("API error: HTTP %d: %s", resp.StatusCode, msg)
		}
		return nil, fmt.Errorf("API error: HTTP %d", resp.StatusCode)
	}

	var res ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("error parsing API response: %w", err)
	}

	if !res.Found {
		return nil, ErrNotFound
	}

	return &res.Page, nil
}
