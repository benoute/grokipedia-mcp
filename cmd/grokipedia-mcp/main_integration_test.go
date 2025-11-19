package main

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestIntegrationMCP(t *testing.T) {
	// // Build the server binary first
	// if err := exec.Command("go", "build", "-o", "grokipedia-mcp-test").Run(); err != nil {
	// 	t.Fatalf("Failed to build server: %v", err)
	// }
	// defer exec.Command("rm", "grokipedia-mcp-test").Run()

	// Create MCP client that connects to our server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)

	// Connect to the server using CommandTransport
	transport := &mcp.CommandTransport{
		// Command: exec.Command("./grokipedia-mcp-test"),
		Command: exec.Command("go", "run", "."),
	}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		t.Fatalf("Failed to connect to MCP server: %v", err)
	}
	defer session.Close()

	// Test tools/list
	t.Run("ListTools", func(t *testing.T) {
		toolsResult, err := session.ListTools(ctx, &mcp.ListToolsParams{})
		if err != nil {
			t.Fatalf("ListTools failed: %v", err)
		}

		tools := toolsResult.Tools
		if len(tools) != 2 {
			t.Errorf("Expected 2 tools, got %d", len(tools))
		}

		toolNames := make(map[string]bool)
		for _, tool := range tools {
			toolNames[tool.Name] = true
		}

		if !toolNames["search_grokipedia"] {
			t.Error("Expected search_grokipedia tool")
		}
		if !toolNames["get_grokipedia_page"] {
			t.Error("Expected get_grokipedia_page tool")
		}
	})

	// Test search tool (this will make real API calls)
	t.Run("SearchTool", func(t *testing.T) {
		params := &mcp.CallToolParams{
			Name: "search_grokipedia",
			Arguments: map[string]any{
				"query": "Grok",
			},
		}

		result, err := session.CallTool(ctx, params)
		if err != nil {
			t.Fatalf("CallTool search_grokipedia failed: %v", err)
		}

		if result.IsError {
			t.Logf("Search tool returned error (expected for real API): %v", result.Content)
			return // This is expected since we're hitting the real API
		}

		// If we get here, validate the response structure
		if len(result.Content) == 0 {
			t.Error("Expected at least one content item")
		}
	})

	// Test get page tool (this will make real API calls)
	t.Run("GetPageTool", func(t *testing.T) {
		params := &mcp.CallToolParams{
			Name: "get_grokipedia_page",
			Arguments: map[string]any{
				"slug": "Satoshi",
			},
		}

		result, err := session.CallTool(ctx, params)
		if err != nil {
			t.Fatalf("CallTool get_grokipedia_page failed: %v", err)
		}

		if result.IsError {
			t.Logf("Get page tool returned error (expected for real API): %v", result.Content)
			return // This is expected since we're hitting the real API
		}

		// If we get here, validate the response structure
		if len(result.Content) == 0 {
			t.Error("Expected at least one content item")
		}
	})
}
