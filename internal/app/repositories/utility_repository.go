package repositories

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/grid-stream-org/api/internal/app/logic"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
	"github.com/grid-stream-org/batcher/pkg/bqclient"
)


type UtilityRepository interface {
	CreateUtility(ctx context.Context, data *models.Utility) error
	GetUtility(ctx context.Context, id string) (*models.Utility, error)
	UpdateUtility(ctx context.Context, id string, data *models.Utility) error
	DeleteUtility(ctx context.Context, id string) error
}
type utilityRepository struct {
	client bqclient.BQClient
	logger *slog.Logger
}

func NewUtilityRepository(client bqclient.BQClient, logger *slog.Logger) UtilityRepository {
	return &utilityRepository{
		client: client,
		logger: logger,
	}
}

func (r *utilityRepository) CreateUtility(ctx context.Context, data *models.Utility) error {
	data.ID = uuid.New().String()

	if err := r.client.Put(ctx, "utilities", data); err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to create utility", err)
	}

	return nil
}

func (r *utilityRepository) GetUtility(ctx context.Context, id string) (*models.Utility, error) {
	var util models.Utility
	if err := r.client.Get(ctx, "utilities", id, &util); err != nil {
		if err == bqclient.ErrNotFound {
			return nil, custom_error.New(http.StatusNotFound, "Utility id not found", err)
		}
		return nil, custom_error.New(http.StatusInternalServerError, "Failed to fetch utility", err)
	}
	return &util, nil
}

func (r *utilityRepository) UpdateUtility(ctx context.Context, id string, data *models.Utility) error {

	// // we already validated display name is not empty in the handler
	updates := logic.ExtractBody(data)

	if err := r.client.Update(ctx, "utilities", id, updates); err != nil {
		if err == bqclient.ErrNotFound {
			return custom_error.New(http.StatusNotFound, "Utility id not found", err)
		}
		return custom_error.New(http.StatusInternalServerError, "Failed updating utility", err)
	}

	return nil
}

func (r *utilityRepository) DeleteUtility(ctx context.Context, id string) error {
	if err := r.client.Delete(ctx, "utilities", id); err != nil {
		if err == bqclient.ErrNotFound {
			return custom_error.New(http.StatusNotFound, "Utility id not found", err)
		}
		return custom_error.New(http.StatusInternalServerError, "Failed to delete utility id", err)
	}
	return nil
}
