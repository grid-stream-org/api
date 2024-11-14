package database

import (
	"context"
	"log/slog"

	"cloud.google.com/go/bigquery"
	"github.com/grid-stream-org/api/internal/config"
)

var BigQueryClient *bigquery.Client

// InitializeBigQueryClient sets up the BigQuery client
func InitializeBigQueryClient(ctx context.Context, cfg *config.Config, log *slog.Logger) (*bigquery.Client, error) {
    // google api will automatically load credential file from .env file as long as it is named properly
    // all we need to pass in is the project id, unreal!
	client, err := bigquery.NewClient(ctx, cfg.Database.BigQuery.ProjectID)
	if err != nil {
		return nil, err
	}

	BigQueryClient = client
	log.Info("BigQuery client initialized successfully")
	return BigQueryClient, nil
}

// CloseBigQueryClient closes the BigQuery client
func CloseBigQueryClient(log *slog.Logger) {
	if BigQueryClient != nil {
		BigQueryClient.Close()
        log.Info("Closed client successfully")
	}
}
