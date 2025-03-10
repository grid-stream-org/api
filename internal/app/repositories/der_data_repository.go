package repositories

import (
	"context"
	"log/slog"
	"net/http"

	"cloud.google.com/go/bigquery"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
	"github.com/grid-stream-org/go-commons/pkg/bqclient"
	"google.golang.org/api/iterator"
)

type DERRepository interface {
	GetDERDataByProjectID(ctx context.Context, id string) ([]models.DERData, error)
}

type derRepository struct {
	client bqclient.BQClient
}

func NewDERRepository(client bqclient.BQClient, log *slog.Logger) DERRepository {
	return &derRepository{client: client}
}

func (r *derRepository) GetDERDataByProjectID(ctx context.Context, id string) ([]models.DERData, error) {
	query := `
        SELECT 
            id, der_id, timestamp, current_output, units, project_id, baseline, is_online, is_standalone, connection_start_at, current_soc, power_meter_measurement, contract_threshold
        FROM gridstream_operations.der_data
        WHERE project_id = @project_id
        ORDER BY timestamp DESC;
    `

		params := []bigquery.QueryParameter{
        {Name: "project_id", Value: id},
    }

    it, err := r.client.Query(ctx, query, params) 
    if err != nil {
        return nil, custom_error.New(http.StatusInternalServerError, "Failed to list DER data", err)
    }

    derData := []models.DERData{}
    
    for {
        var item models.DERData
        err := it.Next(&item)
        
        if err == iterator.Done {
            break 
        }

        if err != nil {
            return nil, custom_error.New(http.StatusInternalServerError, "Error reading DER data", err)
        }
        derData = append(derData, item)
    }

    return derData, nil
}
