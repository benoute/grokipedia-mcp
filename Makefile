.PHONY: build build-all test clean run fmt vet

build:
	go build -o grokipedia-mcp ./cmd/grokipedia-mcp

build-all:
	@mkdir -p bin
	@for os in linux darwin; do \
		for arch in amd64 arm64; do \
			echo "Building for $$os/$$arch"; \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch \
				go build -o bin/grokipedia-mcp-$$os-$$arch ./cmd/grokipedia-mcp; \
		done; \
	done

test:
	go test -v ./...

clean:
	rm -f grokipedia-mcp
	rm -rf bin

run: build
	./grokipedia-mcp

fmt:
	go fmt ./...

vet:
	go vet ./...
