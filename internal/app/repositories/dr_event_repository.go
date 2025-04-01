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

type DREventRepository interface {
	CreateDREvent(ctx context.Context, data *models.DREvents) error
	GetDREvent(ctx context.Context, id string) (*models.DREvents, error)
	UpdateDREvent(ctx context.Context, id string, data *models.DREvents) error
	DeleteDREvent(ctx context.Context, id string) error
	GetDREventsByProjectID(ctx context.Context, id string) ([]models.DREvents, error)
}

type drEventRepository struct {
	client bqclient.BQClient
}

func NewDREventRepository(client bqclient.BQClient, log *slog.Logger) DREventRepository {
	return &drEventRepository{client: client}
}

func (r *drEventRepository) CreateDREvent(ctx context.Context, data *models.DREvents) error {
	query := `
        DECLARE inserted BOOL DEFAULT FALSE;

        INSERT INTO gridstream_operations.dr_events (id, utility_id, start_time, end_time)
        SELECT 
            @id,
            @utility_id,
            TIMESTAMP(@start_time),
            TIMESTAMP(@end_time)
        FROM gridstream_operations.utilities p
        WHERE p.id = @utility_id;

        SET inserted = EXISTS(
            SELECT 1
            FROM gridstream_operations.dr_events c
            WHERE c.id = @id
        );
        SELECT inserted AS inserted;`

	params := []bigquery.QueryParameter{
		{Name: "id", Value: data.ID},
		{Name: "utility_id", Value: data.UtilityID},
		{Name: "start_time", Value: data.StartTime},
		{Name: "end_time", Value: data.EndTime},
	}

	it, err := r.client.Query(ctx, query, params)
	if err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to create demand response event", err)
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

	// If no row was inserted, return an error, likely incorrect  id
	if !inserted {
		return custom_error.New(http.StatusBadRequest, fmt.Sprintf("Failed to insert, please make sure your utility id is correct: %s", data.UtilityID), nil)
	}

	return nil
}

func (r *drEventRepository) GetDREvent(ctx context.Context, id string) (*models.DREvents, error) {
	var event models.DREvents
	if err := r.client.Get(ctx, "dr_events", id, &event); err != nil {
		if err == bqclient.ErrNotFound {
			return nil, custom_error.New(http.StatusNotFound, "demand response event not found", err)
		}
		return nil, custom_error.New(http.StatusInternalServerError, "Failed to retrieve demand response event", err)
	}

	return &event, nil
}

func (r *drEventRepository) UpdateDREvent(ctx context.Context, id string, data *models.DREvents) error {
	updates := logic.ExtractBody(data)

	if len(updates) == 0 {
		return custom_error.New(http.StatusBadRequest, "No fields to update", errors.New("No fields to update"))
	}

	if err := r.client.Update(ctx, "dr_events", id, updates); err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to update", err)
	}

	return nil
}

func (r *drEventRepository) DeleteDREvent(ctx context.Context, id string) error {
	if err := r.client.Delete(ctx, "dr_events", id); err != nil {
		if err == bqclient.ErrNotFound {
			return custom_error.New(http.StatusNotFound, "demand response event id not found", err)
		}
		return custom_error.New(http.StatusInternalServerError, "Failed to delete demand response event", err)
	}
	return nil
}

func (r *drEventRepository) GetDREventsByProjectID(ctx context.Context, id string) ([]models.DREvents, error) {
	query := `
		SELECT
			dr.id AS id,
			dr.start_time,
			dr.end_time,
			dr.utility_id,
			u.display_name AS utility_name
		FROM
			gridstream_operations.projects AS p
		JOIN
			gridstream_operations.dr_events AS dr
			ON p.utility_id = dr.utility_id
		JOIN
			gridstream_operations.utilities AS u
			ON dr.utility_id = u.id
		WHERE
			p.id = @project_id
		ORDER BY
			dr.start_time DESC;
	`

	params := []bigquery.QueryParameter{
		{Name: "project_id", Value: id},
	}

	it, err := r.client.Query(ctx, query, params)
	if err != nil {
		return nil, custom_error.New(http.StatusInternalServerError, "Failed to list der metadata", err)
	}

	drEvents := []models.DREvents{}
	for {
		var item models.DREvents
		err := it.Next(&item)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, custom_error.New(http.StatusInternalServerError, "Error reading der metadata", err)
		}
		drEvents = append(drEvents, item)
	}
	return drEvents, nil
}
