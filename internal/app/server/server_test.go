package server

import (
	"log/slog"
	"os"
	"testing"

	"github.com/grid-stream-org/api/internal/config"
)

// MockBigQueryClient is a mock implementation of *bigquery.Client
type MockBigQueryClient struct{}

// MockAuthClient is a mock implementation of *auth.Client
type MockAuthClient struct{}

type SlogConf struct{}

// TODO: Make NewServer take in interfaces instead of concrete implementation of bigquery client and firebase client
func TestNewServer(t *testing.T) {
	// Set up mock dependencies
	mockBQClient := &MockBigQueryClient{}
	mockAuthClient := &MockAuthClient{}
    mockSlog := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	mockLogger := slog.New(mockSlog)

	// Create a dummy config
	cfg := &config.Config{
		Port: 8080,
	}

	// Call the NewServer function with mocked dependencies
	srv := NewServer(cfg, mockBQClient, mockAuthClient, mockLogger)

	// Validate that the server is created correctly
	if srv == nil {
		t.Fatal("Expected server to be created, but got nil")
	}

	if srv.Addr != ":8080" {
		t.Errorf("Expected server address to be :8080, but got %s", srv.Addr)
	}

	// Check timeouts
	if srv.ReadTimeout != 10*time.Second {
		t.Errorf("Expected ReadTimeout to be 10s, but got %v", srv.ReadTimeout)
	}
	if srv.WriteTimeout != 10*time.Second {
		t.Errorf("Expected WriteTimeout to be 10s, but got %v", srv.WriteTimeout)
	}

	// Check if the handler is set
	if srv.Handler == nil {
		t.Fatal("Expected Handler to be set, but got nil")
	}
}