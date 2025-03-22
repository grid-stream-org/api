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

type ProjectSummaryRepository interface {
	GetProjectSummary(ctx context.Context, utilityID string) ([]models.ProjectSummary, error)
}

type projectSummaryRepository struct {
	client bqclient.BQClient
	log    *slog.Logger
}

func NewProjectSummaryRepository(client bqclient.BQClient, log *slog.Logger) ProjectSummaryRepository {
	return &projectSummaryRepository{client: client, log: log}
}

func (r *projectSummaryRepository) GetProjectSummary(ctx context.Context, utilityID string) ([]models.ProjectSummary, error) {
	//straight claude 
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