package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/benoute/grokipedia-client-go/pkg/grokipedia"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func SearchGrokipedia(ctx context.Context, req *mcp.CallToolRequest, input grokipedia.SearchInput) (
	*mcp.CallToolResult,
	grokipedia.SearchOutput,
	error,
) {
	output, err := grokipedia.Search(ctx, input)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
			IsError: true,
		}, grokipedia.SearchOutput{}, nil
	}

	// Format the display text (simplified, as full formatting was moved to library but we have structured data)
	var contentLines []string
	contentLines = append(contentLines, fmt.Sprintf("Search results for '%s':", input.Query))
	contentLines = append(contentLines, "")
	for i, slug := range output.Results {
		contentLines = append(contentLines, fmt.Sprintf("%d. %s", i+1, slug))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: strings.Join(contentLines, "\n")},
		},
	}, output, nil
}

func GetGrokipediaPage(ctx context.Context, req *mcp.CallToolRequest, input grokipedia.GetPageInput) (
	*mcp.CallToolResult,
	grokipedia.GetPageOutput,
	error,
) {
	output, err := grokipedia.GetPage(ctx, input.Slug)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
			IsError: true,
		}, grokipedia.GetPageOutput{}, nil
	}

	// Format the display text
	var contentLines []string
	contentLines = append(contentLines, fmt.Sprintf("# %s", output.Title))
	contentLines = append(contentLines, "")
	contentLines = append(contentLines, output.Content)

	if len(output.Citations) > 0 {
		contentLines = append(contentLines, "")
		contentLines = append(contentLines, "## Citations")
		for _, citation := range output.Citations {
			contentLines = append(contentLines, fmt.Sprintf("[%s] %s - %s", citation.ID, citation.Title, citation.URL))
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: strings.Join(contentLines, "\n")},
		},
	}, output, nil
}

func main() {
	// Create a server with search and page retrieval tools
	server := mcp.NewServer(&mcp.Implementation{Name: "grokipedia-mcp", Version: "v1.0.0"}, nil)

	// Add search tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_grokipedia",
		Description: "Search an AI-generated online encyclopedia for articles and information on various topics, providing titles, snippets, and metadata",
	}, SearchGrokipedia)

	// Add page retrieval tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_grokipedia_page",
		Description: "Retrieve the full content of a specific encyclopedia page by its identifier, including title, content, and citations",
	}, GetGrokipediaPage)

	// Run the server over stdin/stdout
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
