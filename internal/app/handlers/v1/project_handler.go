package handlers

import (
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grid-stream-org/api/internal/app/repositories"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
)

// ProjectHandlers contains the repository and logger
type ProjectHandlers struct {
	Repo *repositories.ProjectRepository
	Log  *slog.Logger
}

// NewProjectHandlers creates a new instance of ProjectHandlers
func NewProjectHandlers(repo *repositories.ProjectRepository, log *slog.Logger) *ProjectHandlers {
	return &ProjectHandlers{Repo: repo, Log: log}
}

// GetProjectHandler handles retrieving a project by ID
func (h *ProjectHandlers) GetProjectHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id") // Get the project ID from the URL

	project, err := h.Repo.GetProject(r.Context(), id)
	if err != nil {
		return err
	}

	// Return the project as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(project)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *ProjectHandlers) CreateProjectHandler(w http.ResponseWriter, r *http.Request) error {

	var req models.Project
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return custom_error.New(http.StatusBadRequest, "Invalid request payload", err)
	}
	if req.UtilityID == "" || req.Location == "" {
		return custom_error.New(http.StatusBadRequest, "All fields (utilityId, userId, location) are required", nil)
	}
	err := h.Repo.CreateProject(r.Context(), &models.Project{
		ID:        uuid.New().String(),
		UtilityID: req.UtilityID,
		UserID:    "",
		Location:  req.Location,
	})
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *ProjectHandlers) UpdateProjectHandler(w http.ResponseWriter, r *http.Request) error {
	var req models.Project
	id := chi.URLParam(r, "id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return custom_error.New(http.StatusBadRequest, "Invalid request payload", err)
	}
	if id == "" {
		return custom_error.New(http.StatusBadRequest, "Project ID is required", nil)
	}

	if req.ID != "" {
		return custom_error.New(http.StatusBadRequest, "Updating project id is not allowed", nil)
	}
	err := h.Repo.UpdateProject(r.Context(), id, &models.Project{
		UtilityID: req.UtilityID,
		UserID:    req.UserID,
		Location:  req.Location,
	})

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *ProjectHandlers) DeleteProjectHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if id == "" {
		return custom_error.New(http.StatusBadRequest, "Project ID is required", nil)
	}
	err := h.Repo.DeleteProject(r.Context(), id)

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
