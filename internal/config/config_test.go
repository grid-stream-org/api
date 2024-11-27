package config

import (
	"testing"
	"time"
)

func TestPopulatedConfig(t *testing.T) {
    // Set all necessary environment variables
    t.Setenv("GO_ENV", "test")
    t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./fbq-file.json")
    t.Setenv("BIGQUERY_PROJECT_ID", "test-id")
    t.Setenv("GOOGLE_CLOUD_PROJECT", "./fb-file.json")
    t.Setenv("FIREBASE_PROJECT_ID", "fb-test-id")
    t.Setenv("LOG_LEVEL", "INFO")
    t.Setenv("LOG_FORMAT", "text")
    t.Setenv("LOG_OUTPUT", "stdout")
    t.Setenv("PORT", "8080")
    t.Setenv("TIMEOUT", "10s")

    // expected values
    port := 8080
    timeout := 10 * time.Second
    bqFile := "./fbq-file.json"
    logLevel := "INFO"
    logFormat := "text"
    logOutput := "stdout"
    fbId := "fb-test-id"

    cfg, err := Load()
    if err != nil {
        t.Fatalf("Error loading config: %v", err)
    }

    // Perform assertions
    if cfg.Port != port {
        t.Fatalf("Port mismatch: expected %d, got %d", port, cfg.Port)
    }

    if cfg.Timeout != timeout {
        t.Fatalf("Timeout mismatch: expected %v, got %v", timeout, cfg.Timeout)
    }

    if cfg.Database.BigQuery.CredsFile != bqFile {
        t.Fatalf("BigQuery creds file mismatch: expected %s, got %s", bqFile, cfg.Database.BigQuery.CredsFile)
    }

    if cfg.Logger.Level != logLevel {
        t.Fatalf("Log level mismatch: expected %s, got %s", logLevel, cfg.Logger.Level)
    }

    if cfg.Logger.Format != logFormat {
        t.Fatalf("Log format mismatch: expected %s, got %s", logFormat, cfg.Logger.Format)
    }

    if cfg.Logger.Output != logOutput {
        t.Fatalf("Log output mismatch: expected %s, got %s", logOutput, cfg.Logger.Output)
    }

    if cfg.Firebase.ProjectID != fbId {
        t.Fatalf("Firebase project ID mismatch: expected %s, got %s", fbId, cfg.Firebase.ProjectID)
    }
}

func TestMissingCredential(t *testing.T) {
    t.Setenv("GO_ENV", "test")
    t.Setenv("BIGQUERY_PROJECT_ID", "test-id")
    t.Setenv("GOOGLE_CLOUD_PROJECT", "./fb-file.json")
    t.Setenv("FIREBASE_PROJECT_ID", "fb-test-id")
    t.Setenv("LOG_LEVEL", "INFO")
    t.Setenv("LOG_FORMAT", "text")
    t.Setenv("LOG_OUTPUT", "stdout")
    t.Setenv("PORT", "8080")
    t.Setenv("TIMEOUT", "10s")

    expect := "missing big query creds"

    _ , err := Load()
    if err == nil {
        t.Fatalf(`Error message should be "%s" but is empty`, expect)
    }
}