# Grokipedia MCP Server

A Model Context Protocol (MCP) server for accessing Grokipedia, the online encyclopedia by xAI.

## Overview

This project provides:
- An **MCP server** that exposes Grokipedia functionality as tools for AI assistants
- A **Go library** (`pkg/grokipedia`) for direct API access in your own applications

## Installation

### Prerequisites

- Go 1.23 or later

### Building from Source

```bash
git clone https://github.com/benoute/grokipedia-mcp.git
cd grokipedia-mcp
make build
```

The binary will be created at `./bin/grokipedia-mcp`.

### Using Go Install

```bash
go install github.com/benoute/grokipedia-mcp/cmd/grokipedia-mcp@latest
```

## Usage

### Claude Desktop (stdio mode)

1. Add the following to your `claude_desktop_config.json`:

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

2. Restart Claude Desktop

### HTTP Mode

Run the server in HTTP mode for use with web-based MCP clients:

```bash
# Default port 8080
./grokipedia-mcp -http

# Custom port
./grokipedia-mcp -http -port 3000
```

### Command-Line Options

| Option  | Default | Description                                      |
|---------|---------|--------------------------------------------------|
| `-http` | `false` | Run as HTTP server instead of stdio              |
| `-port` | `8080`  | Port to listen on when using HTTP mode           |

## Tools Reference

### search_grokipedia

Search Grokipedia for articles and information on various topics.

**Parameters:**

| Parameter | Type   | Required | Default | Description                          |
|-----------|--------|----------|---------|--------------------------------------|
| `query`   | string | Yes      | -       | Search query                         |
| `limit`   | int    | No       | 10      | Maximum number of results to return  |
| `offset`  | int    | No       | 0       | Number of results to skip            |

**Returns:** Array of search results containing:
- `title` - Article title
- `slug` - Article identifier (use with `get_grokipedia_page`)
- `snippet` - Text excerpt with search matches
- `relevanceScore` - Relevance ranking score

### get_grokipedia_page

Retrieve the full content of a Grokipedia page.

**Parameters:**

| Parameter | Type   | Required | Description                                    |
|-----------|--------|----------|------------------------------------------------|
| `slug`    | string | Yes      | Page identifier (from search results or URL)   |

**Returns:** Full article content including:
- Title and content (markdown)
- Citations with URLs
- Images with captions
- Metadata (created/updated dates, view count)

## Go Library

The `pkg/grokipedia` package can be used directly in Go applications:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/benoute/grokipedia-mcp/pkg/grokipedia"
)

func main() {
    ctx := context.Background()

    // Search for articles
    results, err := grokipedia.Search(ctx, "artificial intelligence",
        grokipedia.WithLimit(5),
        grokipedia.WithOffset(0),
    )
    if err != nil {
        log.Fatal(err)
    }

    for _, result := range results {
        fmt.Printf("%s (score: %.2f)\n", result.Title, result.RelevanceScore)
    }

    // Get a specific page
    page, err := grokipedia.GetPage(ctx, "Artificial_intelligence")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Title: %s\n", page.Title)
    fmt.Printf("Content: %s\n", page.Content[:200])
}
```

### Search Options

- `grokipedia.WithLimit(n)` - Limit results (default: 10)
- `grokipedia.WithOffset(n)` - Skip results for pagination (default: 0)

### Page Options

- `grokipedia.WithoutContent()` - Retrieve page metadata without full content

## License

MIT License
