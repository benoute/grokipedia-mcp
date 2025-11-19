package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func parseFlags() (transport string, port string) {
	flag.StringVar(&transport, "transport", "stdio", "Transport mode: stdio or http")
	flag.StringVar(&port, "port", "8080", "Port to listen on for http mode")
	flag.Parse()

	return transport, port
}

func main() {
	transport, port := parseFlags()

	// Create a server with search and page retrieval tools
	server := setupMCPServer()

	// Run the server based on transport
	switch transport {
	case "http":
		fmt.Printf("Server running in HTTP mode on port %s\n", port)
		handler := mcp.NewStreamableHTTPHandler(
			func(*http.Request) *mcp.Server { return server },
			nil,
		)
		log.Fatal(http.ListenAndServe(":"+port, handler))
	case "stdio":
		if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal(errors.New("invalid transport"))
	}
}
