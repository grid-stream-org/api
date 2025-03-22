package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/grid-stream-org/api/internal/app/repositories"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
)

type ProjectSummaryHandler interface {
	GetProjectSummaryHandler(w http.ResponseWriter, r *http.Request) error
}

type projectSummaryHandler struct {
	repo repositories.ProjectSummaryRepository
}

func NewProjectSummaryHandler(repo repositories.ProjectSummaryRepository, log *slog.Logger) ProjectSummaryHandler {
	return &projectSummaryHandler{repo: repo}
}

func (h *projectSummaryHandler) GetProjectSummaryHandler(w http.ResponseWriter, r *http.Request) error {
	utilityID := r.URL.Query().Get("utility_id")
	if utilityID == "" {
		return custom_error.New(http.StatusBadRequest, "utility_id query parameter is required", nil)
	}

	summaries, err := h.repo.GetProjectSummary(r.Context(), utilityID)
	if err != nil {
		return err
	}

	// Create a default summary if no data was found
	var summary models.ProjectSummary
	if len(summaries) > 0 {
		summary = summaries[0]
	} else {
		summary = models.ProjectSummary{
			TotalActive:    0,
			TotalPending:   0,
			TotalThreshold: 0,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(summary); err != nil {
		return custom_error.New(http.StatusInternalServerError, "Error encoding response", err)
	}
	return nil
}