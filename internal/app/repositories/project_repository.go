package repositories

// handles database interactions for Projects
// for now we define a project as an instance of offloading power from a DER to the grid
// projectId	        STRING(REQUIRED)        - Home/Building ID
// utilityId            STRING(REQUIRED)        - Utility Company ID ex: NB Power vs SJ Energy
// connectionStartAt    TIMESTAMP(Required)     - Start of offloading

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/grid-stream-org/api/internal/app/logic"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
	"github.com/grid-stream-org/go-commons/pkg/bqclient"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, data *models.Project) error
	GetProject(ctx context.Context, id string) (*models.Project, error)
	UpdateProject(ctx context.Context, id string, data *models.Project) error
	DeleteProject(ctx context.Context, id string) error
}

type projectRepository struct {
	client bqclient.BQClient
}

func NewProjectRepository(client bqclient.BQClient, log *slog.Logger) ProjectRepository {
	return &projectRepository{client: client}
}

func (r *projectRepository) CreateProject(ctx context.Context, post *models.Project) error {

	if err := r.client.Put(ctx, "projects", post); err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to create project", err)
	}

	return nil
}

func (r *projectRepository) UpdateProject(ctx context.Context, id string, post *models.Project) error {

	updates := logic.ExtractBody(post)

	if len(updates) == 0 {
		return custom_error.New(http.StatusBadRequest, "No fields to update", nil)
	}

	if err := r.client.Update(ctx, "projects", id, updates); err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to update", err)
	}

	return nil
}

func (r *projectRepository) GetProject(ctx context.Context, id string) (*models.Project, error) {
	var proj models.Project
	if err := r.client.Get(ctx, "projects", id, &proj); err != nil {
		if err == bqclient.ErrNotFound {
			return nil, custom_error.New(http.StatusNotFound, "Project id not found", err)
		}
		return nil, custom_error.New(http.StatusInternalServerError, "Failed to get project", err)
	}
	return &proj, nil
}

func (r *projectRepository) DeleteProject(ctx context.Context, id string) error {
	if err := r.client.Delete(ctx, "projects", id); err != nil {
		if err == bqclient.ErrNotFound {
			return custom_error.New(http.StatusNotFound, "Project id not found", err)
		}
		return custom_error.New(http.StatusInternalServerError, "Failed to get project", err)
	}
	return nil
}
