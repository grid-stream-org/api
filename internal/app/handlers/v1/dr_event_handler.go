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

type DREventHandlers interface {
	CreateDREventHandler(w http.ResponseWriter, r *http.Request) error
	GetDREventHandler(w http.ResponseWriter, r *http.Request) error
	UpdateDREventHandler(w http.ResponseWriter, r *http.Request) error
	DeleteDREventHandler(w http.ResponseWriter, r *http.Request) error
	GetDREventsByProjectIDHandler(w http.ResponseWriter, r *http.Request) error
}

type drEventHandlers struct {
	Repo repositories.DREventRepository
	Log  *slog.Logger
}

func NewDREventHandlers(repo repositories.DREventRepository, log *slog.Logger) DREventHandlers {
	return &drEventHandlers{Repo: repo, Log: log}
}

func (h *drEventHandlers) GetDREventHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	drEvent, err := h.Repo.GetDREvent(r.Context(), id)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(drEvent)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *drEventHandlers) CreateDREventHandler(w http.ResponseWriter, r *http.Request) error {
	var req models.DREvents
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return custom_error.New(http.StatusBadRequest, "Invalid request payload", err)
	}
	if req.UtilityID == "" || req.StartTime.String() == "" || req.EndTime.String() == "" {
		return custom_error.New(http.StatusBadRequest, "All fields (utilityId, userId, location) are required", nil)
	}
	req.ID = uuid.New().String()

	if err := h.Repo.CreateDREvent(r.Context(), &req); err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(req); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *drEventHandlers) UpdateDREventHandler(w http.ResponseWriter, r *http.Request) error {
	var req models.DREvents
	id := chi.URLParam(r, "id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return custom_error.New(http.StatusBadRequest, "Invalid request payload", err)
	}
	if id == "" {
		return custom_error.New(http.StatusBadRequest, "Demand Response Event ID is required", nil)
	}

	if req.ID != "" {
		return custom_error.New(http.StatusBadRequest, "Updating Demand Response Event id is not allowed", nil)
	}

	if req.ID != "" {
		return custom_error.New(http.StatusBadRequest, "Updating utility ID is not allowed", nil)
	}
	err := h.Repo.UpdateDREvent(r.Context(), id, &req)

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *drEventHandlers) DeleteDREventHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if id == "" {
		return custom_error.New(http.StatusBadRequest, "demand response ID is required", nil)
	}
	err := h.Repo.DeleteDREvent(r.Context(), id)

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *drEventHandlers) GetDREventsByProjectIDHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "projectId")
	if id == "" {
		return custom_error.New(http.StatusBadRequest, "project ID is required", nil)
	}
	events, err := h.Repo.GetDREventsByProjectID(r.Context(), id)
	if err != nil {
		return err
	}

	if events == nil {
		events = []models.DREvents{}
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(events)
}
