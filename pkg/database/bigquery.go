package database

import (
	"context"
	"log/slog"

	"cloud.google.com/go/bigquery"
	"github.com/grid-stream-org/api/internal/config"
)

var BigQueryClient *bigquery.Client

// InitializeBigQueryClient sets up the BigQuery client
func InitializeBigQueryClient(ctx context.Context, cfg *config.Config, log *slog.Logger) error {
	client, err := bigquery.NewClient(ctx, cfg.Database.BigQuery.ProjectID)
	if err != nil {
		return err
	}

	BigQueryClient = client
	log.Info("BigQuery client initialized successfully")
	return nil
}

// CloseBigQueryClient closes the BigQuery client
func CloseBigQueryClient(log *slog.Logger) {
	if BigQueryClient != nil {
		BigQueryClient.Close()
        log.Info("Closed client successfully")
	}
}
