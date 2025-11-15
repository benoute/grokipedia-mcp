package grokipedia

import "encoding/json"

type SearchInput struct {
	Query  string `json:"query" jsonschema:"The search query for Grokipedia"`
	Limit  int    `json:"limit,omitempty" jsonschema:"Maximum number of results to return (default: 10)"`
	Offset int    `json:"offset,omitempty" jsonschema:"Number of results to skip (default: 0)"`
}

type SearchOutput struct {
	Results []string `json:"results" jsonschema:"List of search results or article snippets"`
}

// MarshalJSON ensures Results is never null in JSON output
func (s SearchOutput) MarshalJSON() ([]byte, error) {
	type Alias SearchOutput
	alias := Alias(s)
	if alias.Results == nil {
		alias.Results = []string{}
	}
	return json.Marshal(alias)
}

type GetPageInput struct {
	Slug string `json:"slug" jsonschema:"The page slug to retrieve (e.g., 'United_Petroleum')"`
}

type GetPageOutput struct {
	Title     string     `json:"title" jsonschema:"The page title"`
	Content   string     `json:"content" jsonschema:"The full page content"`
	Citations []Citation `json:"citations" jsonschema:"List of citations"`
}

// MarshalJSON ensures Citations is never null in JSON output
func (g GetPageOutput) MarshalJSON() ([]byte, error) {
	type Alias GetPageOutput
	alias := Alias(g)
	if alias.Citations == nil {
		alias.Citations = []Citation{}
	}
	return json.Marshal(alias)
}

type GrokipediaSearchResponse struct {
	Results    []SearchResult `json:"results"`
	TotalCount int            `json:"total_count,omitempty"`
}

type SearchResult struct {
	Title          string      `json:"title"`
	Slug           string      `json:"slug"`
	Snippet        string      `json:"snippet,omitempty"`
	RelevanceScore float64     `json:"relevanceScore,omitempty"`
	ViewCount      json.Number `json:"viewCount,omitempty"`
}

type GrokipediaPageResponse struct {
	Page  PageData `json:"page"`
	Found bool     `json:"found"`
}

type PageData struct {
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Citations []Citation `json:"citations"`
	Images    []Image    `json:"images"`
	Metadata  Metadata   `json:"metadata"`
	Slug      string     `json:"slug"`
}

type Citation struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Image struct {
	URL         string `json:"url"`
	Caption     string `json:"caption"`
	Description string `json:"description"`
}

type Metadata struct {
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Views     int    `json:"views"`
}
