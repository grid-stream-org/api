package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/grid-stream-org/api/internal/app/repositories"
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
func (h *ProjectHandlers) GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // Get the project ID from the URL

	project, err := h.Repo.GetProject(context.Background(), id)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Return the project as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(project)
    if err != nil {
        http.Error(w, "Encoding failed", http.StatusNotFound)
		return
    }
}

// TODO: CreateProjectHandler handles creating a new project
func (h *ProjectHandlers) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	// logic for creating a project, here we would use our models/project.go to verify that we have a valid project
	// don't forget to add import 	"github.com/grid-stream-org/api/internal/models"
}

// TODO: UpdateProjectHandler handles updating a project
func (h *ProjectHandlers) UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
