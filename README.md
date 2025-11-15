# Grokipedia Client

A Go library and Model Context Protocol (MCP) server for accessing Grokipedia.

## Components

This project provides two ways to interact with Grokipedia:

### 📚 Go Library (`pkg/grokipedia`)
A standalone Go package for direct API access to Grokipedia's search and page retrieval functionality.

### 🤖 MCP Server (`cmd/grokipedia-mcp`)
An MCP-compatible server that exposes Grokipedia functionality as tools for AI assistants like Claude.

## What is Grokipedia?

Grokipedia is an AI-generated online encyclopedia launched by xAI, providing articles created primarily by Grok, xAI's large language model. It's positioned as an alternative to traditional encyclopedias with a focus on comprehensive, unbiased knowledge.

## Features

Both the library and MCP server provide access to:

1. **🔍 Full-text search**: Find articles by querying Grokipedia's search API
2. **📄 Page retrieval**: Get complete article content, titles, and citations
3. **⚙️ Configurable parameters**: Set custom limits and offsets for search results

### MCP Tools

The MCP server exposes these as tools for AI assistants:

- **`search_grokipedia`**: Search for articles with configurable limit/offset
- **`get_grokipedia_page`**: Retrieve full page content by slug

## Prerequisites

### For the Go Library
- Go 1.19 or later

### For the MCP Server
- Go 1.19 or later
- Claude Desktop or another MCP-compatible client

## Installation

1. Clone or download this repository
2. Navigate to the project directory
3. Install dependencies:
   ```bash
   go mod tidy
   ```

### Building the Components

**MCP Server:**
```bash
make build
# or
go build -o grokipedia-mcp ./cmd/grokipedia-mcp
```

**Go Library:** Add to your project with:
```bash
go get github.com/benoute/grokipedia-client-go/pkg/grokipedia
```

## Usage

### 📚 Using the Go Library

Add the package to your Go project:

```bash
go get github.com/benoute/grokipedia-client-go/pkg/grokipedia
```

```go
import "github.com/benoute/grokipedia-client-go/pkg/grokipedia"

// Search with default parameters (limit: 10, offset: 0)
output, err := grokipedia.Search(context.Background(), grokipedia.SearchInput{Query: "quantum computing"})
if err != nil {
    // handle error
}
// output.Results contains slugs of matching pages

// Search with custom limit and offset
output, err := grokipedia.Search(context.Background(), grokipedia.SearchInput{
    Query:  "artificial intelligence",
    Limit:  20,
    Offset: 10,
})

// Get full page content
page, err := grokipedia.GetPage(context.Background(), "Quantum_computing")
if err != nil {
    // handle error
}
// page contains Title, Content, Citations
```

### 🤖 Using the MCP Server

#### Installation & Setup

1. Build the server:
```bash
make build
# or
go build -o grokipedia-mcp ./cmd/grokipedia-mcp
```

2. Configure Claude Desktop by adding to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "grokipedia": {
      "command": "/absolute/path/to/grokipedia-mcp",
      "args": [],
      "env": {}
    }
  }
}
```

3. Restart Claude Desktop

#### Available Tools

Once configured, Claude can use:
- `search_grokipedia` - Search with optional limit/offset parameters
- `get_grokipedia_page` - Retrieve full article content

Example queries:
- "Search Grokipedia for information about artificial intelligence"
- "Get the full content of the Grok page"
- "Find articles about quantum computing with 5 results"

### Integrating with Claude Desktop

To use this server with Claude Desktop:

1. Open your Claude Desktop configuration file:
   - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Windows: `%APPDATA%\Claude\claude_desktop_config.json`

2. Add the server configuration:

```json
{
  "mcpServers": {
    "grokipedia": {
      "command": "/absolute/path/to/grokipedia-mcp",
      "args": [],
      "env": {}
    }
  }
}
```

Replace `/absolute/path/to/grokipedia-mcp` with the actual path to the built binary.

3. Restart Claude Desktop

4. The `search_grokipedia` tool should now be available in Claude Desktop

### Using the tools

Once configured, you can ask Claude questions like:
- "Search for information about artificial intelligence"
- "What is quantum computing?"
- "Find articles on climate change"
- "Get the full content of the United_Petroleum page"
- "Read the article about Python programming"

## How it works

### Go Library
The `pkg/grokipedia` package provides direct HTTP client access to Grokipedia's API endpoints:

- **Search API**: `https://grokipedia.com/api/full-text-search?query={query}&limit={limit}&offset={offset}`
- **Page API**: `https://grokipedia.com/api/page?slug={slug}&includeContent=true&validateLinks=true`

The library handles JSON marshaling, HTTP requests, and response parsing with proper error handling.

### MCP Server
The MCP server (`cmd/grokipedia-mcp`) wraps the Go library and exposes it via the Model Context Protocol:

- Implements MCP protocol using the official Go SDK
- Registers tools that map to library functions
- Handles tool calls and formats responses for AI assistants
- Includes proper error handling and structured output

