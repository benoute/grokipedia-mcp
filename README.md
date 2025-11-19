# Grokipedia MCP Server

A Model Context Protocol (MCP) server for accessing Grokipedia, an online encyclopedia by xAI.

## What is Grokipedia?

Grokipedia is an AI-generated online encyclopedia launched by xAI, providing articles created primarily
by Grok, xAI's large language model. It's positioned as an alternative to traditional encyclopedias
with a focus on comprehensive, unbiased knowledge.

## Overview

This project provides an MCP-compatible server that exposes Grokipedia functionality as tools for AI
assistants like Claude. It also includes a Go library for direct API access.

## Features

The MCP server provides access to Grokipedia through the following tools:

- **`search_grokipedia`**: Full-text search for articles with relevance scores, snippets, and
  configurable limit/offset
- **`get_grokipedia_page`**: Retrieve complete article content, titles, citations, images, and
  metadata

Key capabilities:
1. **🔍 Full-text search**: Query Grokipedia's search API with relevance scores and snippets
2. **📄 Page retrieval**: Get complete article content, titles, citations, images, and metadata
3. **⚙️ Configurable parameters**: Set custom limits and offsets for search results


### Building the MCP Server

```bash
make build
```

## Usage

### Configuration

1. Configure Claude Desktop by adding to `claude_desktop_config.json`:

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

### Command-Line Options

The MCP server supports the following command-line options:

- `-transport`: Transport mode for the server
  - `stdio` (default): Use standard input/output for communication (recommended for Claude Desktop)
  - `http`: Run as an HTTP server
- `-port`: Port to listen on when using HTTP transport (default: 8080)

#### Available Tools

Once configured, Claude can use:
- `search_grokipedia` - Search with optional limit/offset parameters, returns titles, snippets, and
  relevance scores
- `get_grokipedia_page` - Retrieve full article content, citations, images, and metadata

Example queries:
- "Search for information about artificial intelligence"
- "What is quantum computing?"
- "Find articles on nuclear energy"
- "Get the full content of the United_Petroleum page"
- "Read the article about Go programming"

