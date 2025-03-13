package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"cloud.google.com/go/bigquery"
	"github.com/grid-stream-org/api/internal/app/logic"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
	"github.com/grid-stream-org/go-commons/pkg/bqclient"
	"google.golang.org/api/iterator"
)

type ContractRepository interface {
	CreateContract(ctx context.Context, data *models.Contract) error
	GetContract(ctx context.Context, id string) (*models.Contract, error)
	UpdateContract(ctx context.Context, id string, data *models.Contract) error
	DeleteContract(ctx context.Context, id string) error
	GetContractsByProjectID(ctx context.Context, id string) ([]models.Contract, error)
}

type contractRepository struct {
	client bqclient.BQClient
    log *slog.Logger
}

func NewContractRepository(client bqclient.BQClient, log *slog.Logger) ContractRepository {
	return &contractRepository{client: client, log: log}
}

func (r *contractRepository) CreateContract(ctx context.Context, data *models.Contract) error {
	query := `
        DECLARE inserted BOOL DEFAULT FALSE;

        INSERT INTO gridstream_operations.contracts (id, contract_threshold, start_date, end_date, status, project_id)
        SELECT 
            @id,
            @contract_threshold,
            DATE(@start_date),
            DATE(@end_date),
            @status,
            @project_id
        FROM gridstream_operations.projects p
        WHERE p.id = @project_id;

        SET inserted = EXISTS(
            SELECT 1
            FROM gridstream_operations.contracts c
            WHERE c.id = @id
        );
        SELECT inserted AS inserted;`

	params := []bigquery.QueryParameter{
		{Name: "id", Value: data.ID},
		{Name: "contract_threshold", Value: data.ContractThreshold},
		{Name: "start_date", Value: data.StartDate},
		{Name: "end_date", Value: data.EndDate},
		{Name: "status", Value: data.Status},
		{Name: "project_id", Value: data.ProjectID},
	}

	it, err := r.client.Query(ctx, query, params)
	if err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to create contract", err)
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
		return custom_error.New(http.StatusBadRequest, fmt.Sprintf("Failed to insert, please make sure your project id was correct: %s", data.ProjectID), errors.New("Failed to insert, please make sure your project id was correct"))
	}

	return nil
}

func (r *contractRepository) GetContract(ctx context.Context, id string) (*models.Contract, error) {
	var contract models.Contract
	if err := r.client.Get(ctx, "contracts", id, &contract); err != nil {
		if err == bqclient.ErrNotFound {
			return nil, custom_error.New(http.StatusNotFound, "Contract not found", err)
		}
		return nil, custom_error.New(http.StatusInternalServerError, "Failed to retrieve contract", err)
	}

	return &contract, nil
}

func (r *contractRepository) UpdateContract(ctx context.Context, id string, data *models.Contract) error {
	updates := logic.ExtractBody(data)

	if len(updates) == 0 {
		return custom_error.New(http.StatusBadRequest, "No fields to update", errors.New("No fields to update"))
	}

	if err := r.client.Update(ctx, "contracts", id, updates); err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to update", err)
	}

	return nil
}

func (r *contractRepository) DeleteContract(ctx context.Context, id string) error {
	if err := r.client.Delete(ctx, "contracts", id); err != nil {
		if err == bqclient.ErrNotFound {
			return custom_error.New(http.StatusNotFound, "contract id not found", err)
		}
		return custom_error.New(http.StatusInternalServerError, "Failed to delete contract", err)
	}
	return nil
}

func (r *contractRepository) GetContractsByProjectID(ctx context.Context, id string) ([]models.Contract, error) {

	query := `
        SELECT 
            c.id AS id,
            c.contract_threshold,
            c.start_date,
            c.end_date,
            c.status,
            c.project_id
        FROM 
            gridstream_operations.contracts AS c
        WHERE 
            c.project_id = @project_id
        ORDER BY 
            c.start_date DESC;
    `

	params := []bigquery.QueryParameter{
		{Name: "project_id", Value: id},
	}

	it, err := r.client.Query(ctx, query, params)
	if err != nil {
		return nil, custom_error.New(http.StatusInternalServerError, "Failed to list contracts", err)
	}

	contracts := []models.Contract{}
	for {
		var item models.Contract
		err := it.Next(&item)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, custom_error.New(http.StatusInternalServerError, "Error reading contract data", err)
		}
		contracts = append(contracts, item)
	}
	return contracts, nil
}
