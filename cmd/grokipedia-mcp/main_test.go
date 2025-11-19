package main

import (
	"flag"
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	// Save original args and flag command line
	originalArgs := os.Args
	originalFlagCommandLine := flag.CommandLine
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = originalFlagCommandLine
	}()

	tests := []struct {
		name              string
		args              []string
		expectedTransport string
		expectedPort      string
	}{
		{
			name:              "default values",
			args:              []string{"cmd"},
			expectedTransport: "stdio",
			expectedPort:      "8080",
		},
		{
			name:              "custom transport and port",
			args:              []string{"cmd", "-transport", "http", "-port", "9090"},
			expectedTransport: "http",
			expectedPort:      "9090",
		},
		{
			name:              "only transport",
			args:              []string{"cmd", "-transport", "http"},
			expectedTransport: "http",
			expectedPort:      "8080",
		},
		{
			name:              "only port",
			args:              []string{"cmd", "-port", "7070"},
			expectedTransport: "stdio",
			expectedPort:      "7070",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag command line for each test
			flag.CommandLine = flag.NewFlagSet(tt.args[0], flag.ContinueOnError)
			os.Args = tt.args

			transport, port := parseFlags()

			if transport != tt.expectedTransport {
				t.Errorf("Expected transport %s, got %s", tt.expectedTransport, transport)
			}
			if port != tt.expectedPort {
				t.Errorf("Expected port %s, got %s", tt.expectedPort, port)
			}
		})
	}
}
