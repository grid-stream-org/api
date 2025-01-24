package repositories

// handles database interactions for Projects
// for now we define a project as an instance of offloading power from a DER to the grid
// projectId	        STRING(REQUIRED)        - Home/Building ID
// utilityId            STRING(REQUIRED)        - Utility Company ID ex: NB Power vs SJ Energy
// connectionStartAt    TIMESTAMP(Required)     - Start of offloading

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
	"github.com/grid-stream-org/batcher/pkg/bqclient"
)

type ProjectRepository struct {
	client bqclient.BQClient
}

func NewProjectRepository(client bqclient.BQClient, log *slog.Logger) *ProjectRepository {
	return &ProjectRepository{client: client}
}

// create new instance of oflloading
func (r *ProjectRepository) CreateProject(ctx context.Context, post *models.Project) error {
	// Use BigQuery client to insert a new project
	// Example: Use the client to run a query or insert data
	// INSERT INTO A1.Project (projectId, utilityId, connectionStartAt) VALUES ('projId','utilId', '2021-01-26 16:50:03' )

	// we can verify that the utilityId and userId is valid before inserting, this ensures that we only call bq once
	// query := `
	// 	DECLARE utilityExists BOOL;

	// 	SET utilityExists = (
	// 	    SELECT COUNT(1) > 0
	// 	    FROM ` + "`grid-stream.gridstream_operations.utilities`" + `
	// 	    WHERE utility_id = @utilityId
	// 	);

	// 	IF NOT utilityExists THEN
	//         RAISE_ERROR('Invalid utilityId: ' || @utilityId);
	//     IF NOT userExists THEN
	//         RAISE_ERROR('Invalid userId: ' || @userId);
	// 	IF utilityExistsTHEN
	// 	    INSERT INTO ` + "`grid-stream.gridstream_operations.projects`" + `
	// 	    (project_id, utility_id, user_id, location)
	// 	    VALUES (@projectId, @utilityId, @userId, @location);
	// 	END IF;
	// `

	// query := `
	//     INSERT INTO gridstream_operations.projects
	//         VALUES(@id, @utility_id, @user_id, @location);
	// `

	if err := r.client.Put(ctx, "projects", post); err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to create project", err)
	}

	return nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, post *models.Project) error {
	
    updates := make(map[string]any)

	if post.UtilityID != "" {
		updates["utility_id"] = post.UtilityID
	}

	if post.UserID != "" {
		updates["user_id"] = post.UserID
	}

	if post.Location != "" {
		updates["location"] = post.Location
	}

	if len(updates) == 0 {
		return custom_error.New(http.StatusBadRequest, "No fields to update", errors.New("No fields to update"))
	}

    if err := r.client.Update(ctx, "projects", post.ID, updates); err != nil {
        return custom_error.New(http.StatusInternalServerError, "Failed to update", err)
    }

	return nil
}

func (r *ProjectRepository) GetProject(ctx context.Context, id string) (*models.Project, error) {
	var proj models.Project
	if err := r.client.Get(ctx, "projects", id, &proj); err != nil {
		if err == bqclient.ErrNotFound {
			return nil, custom_error.New(http.StatusNotFound, "Project id not found", err)
		}
		return nil, custom_error.New(http.StatusInternalServerError, "Failed to get project", err)
	}
	return &proj, nil
}
