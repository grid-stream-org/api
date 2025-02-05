package repositories

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"cloud.google.com/go/bigquery"
	"github.com/grid-stream-org/api/internal/app/logic"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
	"github.com/grid-stream-org/batcher/pkg/bqclient"
	"google.golang.org/api/iterator"
)

type DERMetadataRepository interface {
	CreateDERMetadata(ctx context.Context, data *models.DERMetadata) error
	GetDERMetadata(ctx context.Context, id string) (*models.DERMetadata, error)
	UpdateDERMetadata(ctx context.Context, id string, data *models.DERMetadata) error
	DeleteDERMetadata(ctx context.Context, id string) error
}
type derMetadataRepository struct {
	client bqclient.BQClient
}

func NewDERMetadataRepository(client bqclient.BQClient, log *slog.Logger) DERMetadataRepository {
	return &derMetadataRepository{client: client}
}

func (r *derMetadataRepository) CreateDERMetadata(ctx context.Context, data *models.DERMetadata) error {
	query := `
        DECLARE inserted BOOL DEFAULT FALSE;

        INSERT INTO gridstream_operations.der_metadata (id, type, nameplate_capacity, power_capacity, project_id)
        SELECT 
            @id,
            @type,
            @nameplate_capacity,
            @power_capacity,
            @project_id
        FROM gridstream_operations.projects p
        WHERE p.id = @project_id;

        SET inserted = EXISTS(
            SELECT 1
            FROM gridstream_operations.der_metadata c
            WHERE c.id = @id
        );
        SELECT inserted AS inserted;`

	params := []bigquery.QueryParameter{
		{Name: "id", Value: data.ID},
		{Name: "type", Value: data.Type},
		{Name: "nameplate_capacity", Value: data.NameplateCapacity},
		{Name: "power_capacity", Value: data.PowerCapacity},
		{Name: "project_id", Value: data.ProjectID},
	}

	it, err := r.client.Query(ctx, query, params)
	if err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to create der metadata", err)
	}

	// Check to see if we inserted
	var inserted bool
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return custom_error.New(http.StatusInternalServerError, "Error reading insertion result", err)
		}

		// Extract the boolean value
		if len(row) > 0 {
			inserted, _ = row[0].(bool)
		}
	}

	// If no row was inserted, return an error, likely incorrect project id
	if !inserted {
		return custom_error.New(http.StatusBadRequest, fmt.Sprintf("Failed to insert, please make sure your project id was correct: %s", data.ProjectID), nil)
	}

	return nil
}

func (r *derMetadataRepository) GetDERMetadata(ctx context.Context, id string) (*models.DERMetadata, error) {
	var derMetadata models.DERMetadata
	if err := r.client.Get(ctx, "der_metadata", id, &derMetadata); err != nil {
		if err == bqclient.ErrNotFound {
			return nil, custom_error.New(http.StatusNotFound, "der not found", err)
		}
		return nil, custom_error.New(http.StatusInternalServerError, "Failed to retrieve der", err)
	}

	return &derMetadata, nil
}

func (r *derMetadataRepository) UpdateDERMetadata(ctx context.Context, id string, data *models.DERMetadata) error {
	updates := logic.ExtractBody(data)

	if len(updates) == 0 {
		return custom_error.New(http.StatusBadRequest, "No fields to update", nil)
	}

	if err := r.client.Update(ctx, "der_metadata", id, updates); err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to update", err)
	}

	return nil
}

func (r *derMetadataRepository) DeleteDERMetadata(ctx context.Context, id string) error {
	if err := r.client.Delete(ctx, "der_metadata", id); err != nil {
		if err == bqclient.ErrNotFound {
			return custom_error.New(http.StatusNotFound, "der id not found", err)
		}
		return custom_error.New(http.StatusInternalServerError, "Failed to delete contract", err)
	}
	return nil
}
