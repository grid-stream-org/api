package repositories

import (
	"context"
	"log/slog"
	"net/http"

	"cloud.google.com/go/bigquery"
	"github.com/google/uuid"
	"github.com/grid-stream-org/api/internal/app/logic"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
	"github.com/grid-stream-org/go-commons/pkg/bqclient"
	"google.golang.org/api/iterator"
)

type UtilityRepository interface {
	CreateUtility(ctx context.Context, data *models.Utility) error
	GetUtility(ctx context.Context, id string) (*models.Utility, error)
	UpdateUtility(ctx context.Context, id string, data *models.Utility) error
	DeleteUtility(ctx context.Context, id string) error
	GetProjectSummary(ctx context.Context, utilityID string) ([]models.ProjectSummary, error)
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

func (r *utilityRepository) GetProjectSummary(ctx context.Context, utilityID string) ([]models.ProjectSummary, error) {
	query := `
-- Total Active Contracts
WITH active_contracts AS (
  SELECT COUNT(*) as total_active
  FROM gridstream_operations.contracts
  WHERE status = 'active'
  AND project_id IN (SELECT id FROM gridstream_operations.projects WHERE utility_id = @utility_id)
),

-- Total Pending Contracts
pending_contracts AS (
  SELECT COUNT(*) as total_pending
  FROM gridstream_operations.contracts
  WHERE status = 'pending'
  AND project_id IN (SELECT id FROM gridstream_operations.projects WHERE utility_id = @utility_id)
),

-- Total Sum of Contract Thresholds
contract_threshold_sum AS (
  SELECT SUM(contract_threshold) as total_threshold
  FROM gridstream_operations.contracts
  WHERE project_id IN (SELECT id FROM gridstream_operations.projects WHERE utility_id = @utility_id)
),

-- Next DR Event
next_dr_event AS (
  SELECT id, start_time, end_time
  FROM gridstream_operations.dr_events
  WHERE start_time > CURRENT_TIMESTAMP()
  AND utility_id = @utility_id
  ORDER BY start_time ASC
  LIMIT 1
),

-- Most Recent DR Event
recent_dr_event AS (
  SELECT id, start_time, end_time
  FROM gridstream_operations.dr_events
  WHERE end_time < CURRENT_TIMESTAMP()
  AND utility_id = @utility_id
  ORDER BY end_time DESC
  LIMIT 1
)

-- Combine all results
SELECT
  a.total_active,
  p.total_pending,
  c.total_threshold,
  n.id as next_event_id,
  n.start_time as next_event_start,
  n.end_time as next_event_end,
  r.id as recent_event_id,
  r.start_time as recent_event_start,
  r.end_time as recent_event_end
FROM
  active_contracts a,
  pending_contracts p,
  contract_threshold_sum c,
  next_dr_event n,
  recent_dr_event r
`
	params := []bigquery.QueryParameter{
		{Name: "utility_id", Value: utilityID},
	}

	it, err := r.client.Query(ctx, query, params)
	if err != nil {
		return nil, custom_error.New(http.StatusInternalServerError, "Failed to fetch project summary", err)
	}

	summaries := []models.ProjectSummary{}
	for {
		var summary models.ProjectSummary
		err := it.Next(&summary)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, custom_error.New(http.StatusInternalServerError, "Error reading project summary data", err)
		}
		summaries = append(summaries, summary)
	}

	// This should only return one row
	return summaries, nil
}