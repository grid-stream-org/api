package repositories

import (
	"context"
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
	BatchCreateDERMetadata(ctx context.Context, data []models.DERMetadata) error
	GetDERMetadata(ctx context.Context, id string) (*models.DERMetadata, error)
	ListDERMetadataByProject(ctx context.Context, id string) ([]models.DERMetadata, error)
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
	if err := r.client.Put(ctx, "der_metadata", data); err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to create der metadata", err)
	}
	return nil
}

func (r *derMetadataRepository) BatchCreateDERMetadata(ctx context.Context, data []models.DERMetadata) error {
	if err := r.client.Put(ctx, "der_metadata", data); err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to create all der metadata", err)
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

func (r *derMetadataRepository) ListDERMetadataByProject(ctx context.Context, id string) ([]models.DERMetadata, error) {
	query := `
        SELECT *
        FROM gridstream_operations.der_metadata
        WHERE project_id = @project_id`
	params := []bigquery.QueryParameter{
		{Name: "project_id", Value: id},
	}
	it, err := r.client.Query(ctx, query, params)
	if err != nil {
		return nil, custom_error.New(http.StatusInternalServerError, "Failed to list der metadata", err)
	}

	derMetadata := []models.DERMetadata{}
	for {
		var item models.DERMetadata
		err := it.Next(&item)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, custom_error.New(http.StatusInternalServerError, "Error reading der metadata", err)
		}
		derMetadata = append(derMetadata, item)
	}

	return derMetadata, nil
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
