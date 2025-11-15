.PHONY: build test clean run fmt vet

build:
	go build -o grokipedia-mcp ./cmd/grokipedia-mcp

test:
	go test -v ./...

clean:
	rm -f grokipedia-mcp

run: build
	./grokipedia-mcp

fmt:
	go fmt ./...

vet:
	go vet ./...
