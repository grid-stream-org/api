package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/grid-stream-org/api/internal/app/repositories"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
)

type ProjectAverageHandler interface {
	CreateProjectAverageHandler(w http.ResponseWriter, r *http.Request) error
	GetProjectAveragesHandler(w http.ResponseWriter, r *http.Request) error
}

type projectAverageHandler struct {
	repo repositories.ProjectAverageRepository
	log  *slog.Logger
}

func NewProjectAverageHandlers(repo repositories.ProjectAverageRepository, log *slog.Logger) ProjectAverageHandler {
	return &projectAverageHandler{repo: repo, log: log}
}

func (h *projectAverageHandler) CreateProjectAverageHandler(w http.ResponseWriter, r *http.Request) error {
	var req models.ProjectAverage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return custom_error.New(http.StatusBadRequest, "Invalid request payload", err)
	}

	// Validate required fields
	if req.ProjectID == "" || req.StartTime.IsZero() || req.EndTime.IsZero() || req.Baseline <= 0 || req.ContractThreshold <= 0 {
		return custom_error.New(http.StatusBadRequest, "All fields (project_id, start_time, end_time, baseline, contract_threshold, average_output) are required", nil)
	}

	// Validate time range
	if req.EndTime.Before(req.StartTime) {
		return custom_error.New(http.StatusBadRequest, "End time must be after start time", nil)
	}

	err := h.repo.CreateProjectAverage(r.Context(), &req)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *projectAverageHandler) GetProjectAveragesHandler(w http.ResponseWriter, r *http.Request) error {
	// Get query parameters
	projectID := r.URL.Query().Get("project_id")
	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")
	
	// Project ID is required
	if projectID == "" {
		return custom_error.New(http.StatusBadRequest, "project_id query parameter is required", nil)
	}

	// If both time parameters are provided, filter by date range
	if startTimeStr != "" && endTimeStr != "" {
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return custom_error.New(http.StatusBadRequest, "Invalid start_time format. Use RFC3339 format (e.g., 2006-01-02T15:04:05Z)", err)
		}

		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return custom_error.New(http.StatusBadRequest, "Invalid end_time format. Use RFC3339 format (e.g., 2006-01-02T15:04:05Z)", err)
		}

		if endTime.Before(startTime) {
			return custom_error.New(http.StatusBadRequest, "End time must be after start time", nil)
		}

		averages, err := h.repo.GetProjectAveragesByDateRange(r.Context(), projectID, startTime, endTime)
		if err != nil {
			return err
		}

		if averages == nil {
			averages = []models.ProjectAverage{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(averages); err != nil {
			return err
		}
		return nil
	}

	// If time parameters are not provided, get all averages for the project
	averages, err := h.repo.GetProjectAveragesByProjectID(r.Context(), projectID)
	if err != nil {
		return err
	}

	if averages == nil {
		averages = []models.ProjectAverage{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(averages); err != nil {
		return err
	}
	return nil
}