package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/benoute/grokipedia/pkg/grokipedia"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchToolInput struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

type SearchToolOutput struct {
	Results []grokipedia.SearchResult `json:"results"`
}

type GetPageToolInput struct {
	Slug string `json:"slug"`
}

func SearchGrokipedia(ctx context.Context, req *mcp.CallToolRequest, input SearchToolInput) (
	*mcp.CallToolResult,
	*SearchToolOutput,
	error,
) {
	var opts []grokipedia.SearchOption
	if input.Limit > 0 {
		opts = append(opts, grokipedia.WithLimit(input.Limit))
	}
	if input.Offset > 0 {
		opts = append(opts, grokipedia.WithOffset(input.Offset))
	}

	results, err := grokipedia.Search(ctx, input.Query, opts...)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
			IsError: true,
		}, nil, nil
	}

	// Format the display text
	var contentLines []string
	contentLines = append(contentLines, fmt.Sprintf("Search results for '%s':", input.Query))
	contentLines = append(contentLines, "")
	for i, result := range results {
		contentLines = append(
			contentLines,
			fmt.Sprintf("%d. %s (slug: %s)", i+1, result.Title, result.Slug),
		)
		if len(result.Snippet) > 0 {
			contentLines = append(contentLines, fmt.Sprintf("   %s", result.Snippet))
		}
		contentLines = append(contentLines, "")
	}

	searchToolOutput := SearchToolOutput{
		Results: results,
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: strings.Join(contentLines, "\n")},
		},
	}, &searchToolOutput, nil
}

func GetGrokipediaPage(ctx context.Context, req *mcp.CallToolRequest, input GetPageToolInput) (
	*mcp.CallToolResult,
	*grokipedia.Page,
	error,
) {
	page, err := grokipedia.GetPage(ctx, input.Slug)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
			IsError: true,
		}, nil, nil
	}

	// Format the display text
	var contentLines []string
	contentLines = append(contentLines, fmt.Sprintf("# %s", page.Title))
	contentLines = append(contentLines, "")
	contentLines = append(contentLines, page.Content)

	if len(page.Citations) > 0 {
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, "## Citations")
		for _, citation := range page.Citations {
			contentLines = append(
				contentLines,
				fmt.Sprintf("[%s] %s - %s", citation.ID, citation.Title, citation.URL),
			)
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: strings.Join(contentLines, "\n")},
		},
	}, page, nil
}

func main() {
	// Create a server with search and page retrieval tools
	server := mcp.NewServer(&mcp.Implementation{Name: "grokipedia-mcp", Version: "v1.0.0"}, nil)

	// Add search tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "search_grokipedia",
		Description: "Search Grokipedia online encyclopedia for articles and information " +
			"on various topics, providing titles, snippets, and metadata",
	}, SearchGrokipedia)

	// Add page retrieval tool
	mcp.AddTool(server, &mcp.Tool{
		Name: "get_grokipedia_page",
		Description: "Retrieve the full content of a specific Grokipedia encyclopedia page by " +
			"its identifier, including title, content, and citations",
	}, GetGrokipediaPage)

	// Run the server over stdin/stdout
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
