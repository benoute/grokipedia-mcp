package grokipedia

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func Search(ctx context.Context, input SearchInput) (SearchOutput, error) {
	// Set defaults
	limit := input.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := input.Offset
	if offset < 0 {
		offset = 0
	}

	// Use Grokipedia's official API
	searchURL := fmt.Sprintf(
		"https://grokipedia.com/api/full-text-search?query=%s&limit=%d&offset=%d",
		url.QueryEscape(input.Query),
		limit,
		offset,
	)

	resp, err := http.Get(searchURL)
	if err != nil {
		return SearchOutput{}, fmt.Errorf("error fetching from Grokipedia API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SearchOutput{}, fmt.Errorf("API error: HTTP %d", resp.StatusCode)
	}

	var searchResp GrokipediaSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return SearchOutput{}, fmt.Errorf("error parsing API response: %w", err)
	}

	if len(searchResp.Results) == 0 {
		return SearchOutput{Results: []string{"No results found"}}, nil
	}

	var results []string

	for i, result := range searchResp.Results {
		if i >= limit { // Limit to the specified number of results
			break
		}

		results = append(results, result.Slug)
	}

	return SearchOutput{Results: results}, nil
}

func GetPage(ctx context.Context, slug string) (GetPageOutput, error) {
	// Use Grokipedia's official API to get page content
	pageURL := fmt.Sprintf(
		"https://grokipedia.com/api/page?slug=%s&includeContent=true&validateLinks=true",
		url.QueryEscape(slug),
	)

	resp, err := http.Get(pageURL)
	if err != nil {
		return GetPageOutput{}, fmt.Errorf("error fetching page from Grokipedia API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return GetPageOutput{}, fmt.Errorf("page not found: %s", slug)
		}
		return GetPageOutput{}, fmt.Errorf("API error: HTTP %d", resp.StatusCode)
	}

	var pageResp GrokipediaPageResponse
	if err := json.NewDecoder(resp.Body).Decode(&pageResp); err != nil {
		return GetPageOutput{}, fmt.Errorf("error parsing API response: %w", err)
	}

	if !pageResp.Found {
		return GetPageOutput{}, fmt.Errorf("page not found: %s", slug)
	}

	page := pageResp.Page

	// Ensure citations is never nil
	if page.Citations == nil {
		page.Citations = []Citation{}
	}

	output := GetPageOutput{
		Title:     page.Title,
		Content:   page.Content,
		Citations: page.Citations,
	}
	if output.Citations == nil {
		output.Citations = []Citation{}
	}

	return output, nil
}
