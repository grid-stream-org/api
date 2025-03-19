package repositories

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
	"github.com/grid-stream-org/go-commons/pkg/bqclient"
	"google.golang.org/api/iterator"
)

type ProjectAverageRepository interface {
	CreateProjectAverage(ctx context.Context, data *models.ProjectAverage) error
	GetProjectAveragesByProjectID(ctx context.Context, projectID string) ([]models.ProjectAverage, error)
	GetProjectAveragesByDateRange(ctx context.Context, projectID string, startTime, endTime time.Time) ([]models.ProjectAverage, error)
}

type projectAverageRepository struct {
	client bqclient.BQClient
	log    *slog.Logger
}

func NewProjectAverageRepository(client bqclient.BQClient, log *slog.Logger) ProjectAverageRepository {
	return &projectAverageRepository{client: client, log: log}
}

func (r *projectAverageRepository) CreateProjectAverage(ctx context.Context, data *models.ProjectAverage) error {
	query := `
		DECLARE inserted BOOL DEFAULT FALSE;

		INSERT INTO gridstream_operations.project_averages (project_id, start_time, end_time, baseline, contract_threshold, average_output)
		SELECT 
			@project_id,
			TIMESTAMP(@start_time),
			TIMESTAMP(@end_time),
			@baseline,
			@contract_threshold,
			@average_output
		FROM gridstream_operations.projects p
		WHERE p.id = @project_id;

		SET inserted = EXISTS(
			SELECT 1
			FROM gridstream_operations.project_averages pa
			WHERE pa.project_id = @project_id 
			AND pa.start_time = TIMESTAMP(@start_time)
		);
		SELECT inserted AS inserted;
	`

	params := []bigquery.QueryParameter{
		{Name: "project_id", Value: data.ProjectID},
		{Name: "start_time", Value: data.StartTime.Format(time.RFC3339)},
		{Name: "end_time", Value: data.EndTime.Format(time.RFC3339)},
		{Name: "baseline", Value: data.Baseline},
		{Name: "contract_threshold", Value: data.ContractThreshold},
		{Name: "average_output", Value: data.AverageOutput},
	}

	it, err := r.client.Query(ctx, query, params)
	if err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to create project average", err)
	}

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

		if len(row) > 0 {
			inserted, _ = row[0].(bool)
		}
	}

	if !inserted {
		return custom_error.New(http.StatusBadRequest, 
			"Failed to insert, please make sure your project id is correct: "+data.ProjectID, 
			errors.New("Failed to insert, project ID not found"))
	}

	return nil
}

func (r *projectAverageRepository) GetProjectAveragesByProjectID(ctx context.Context, projectID string) ([]models.ProjectAverage, error) {
	query := `
		SELECT 
			project_id,
			start_time,
			end_time,
			baseline,
			contract_threshold,
			average_output
		FROM 
			gridstream_operations.project_averages
		WHERE 
			project_id = @project_id
		ORDER BY 
			start_time DESC;
	`

	params := []bigquery.QueryParameter{
		{Name: "project_id", Value: projectID},
	}

	it, err := r.client.Query(ctx, query, params)
	if err != nil {
		return nil, custom_error.New(http.StatusInternalServerError, "Failed to fetch project averages", err)
	}

	averages := []models.ProjectAverage{}
	for {
		var item models.ProjectAverage
		err := it.Next(&item)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, custom_error.New(http.StatusInternalServerError, "Error reading project average data", err)
		}
		averages = append(averages, item)
	}

	return averages, nil
}

func (r *projectAverageRepository) GetProjectAveragesByDateRange(ctx context.Context, projectID string, startTime, endTime time.Time) ([]models.ProjectAverage, error) {
	query := `
		SELECT 
			project_id,
			start_time,
			end_time,
			baseline,
			contract_threshold,
			average_output
		FROM 
			gridstream_operations.project_averages
		WHERE 
			project_id = @project_id
			AND start_time >= TIMESTAMP(@start_time)
			AND end_time <= TIMESTAMP(@end_time)
		ORDER BY 
			start_time ASC;
	`

	params := []bigquery.QueryParameter{
		{Name: "project_id", Value: projectID},
		{Name: "start_time", Value: startTime.Format(time.RFC3339)},
		{Name: "end_time", Value: endTime.Format(time.RFC3339)},
	}

	it, err := r.client.Query(ctx, query, params)
	if err != nil {
		return nil, custom_error.New(http.StatusInternalServerError, "Failed to fetch project averages in date range", err)
	}

	averages := []models.ProjectAverage{}
	for {
		var item models.ProjectAverage
		err := it.Next(&item)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, custom_error.New(http.StatusInternalServerError, "Error reading project average data", err)
		}
		averages = append(averages, item)
	}

	return averages, nil
}