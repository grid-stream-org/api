package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/grid-stream-org/api/internal/app/repositories"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
)

type DERMetadataHandlers interface {
	CreateDERMetadataHandler(w http.ResponseWriter, r *http.Request) error
	GetDERMetadataHandler(w http.ResponseWriter, r *http.Request) error
	UpdateDERMetadataHandler(w http.ResponseWriter, r *http.Request) error
	DeleteDERMetadataHandler(w http.ResponseWriter, r *http.Request) error
}
type derMetadataHandlers struct {
	Repo repositories.DERMetadataRepository
	Log  *slog.Logger
}

func NewDERMetadataHandlers(repo repositories.DERMetadataRepository, log *slog.Logger) DERMetadataHandlers {
	return &derMetadataHandlers{Repo: repo, Log: log}
}

func (h *derMetadataHandlers) CreateDERMetadataHandler(w http.ResponseWriter, r *http.Request) error {
	var req models.DERMetadata
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return custom_error.New(http.StatusBadRequest, "Invalid request payload", err)
	}
	if req.Type == "" || req.ProjectID == "" || req.NameplateCapacity <= 0 || req.PowerCapacity <= 0 {
		return custom_error.New(http.StatusBadRequest, "All fields (Type, ProjectID, Type, NameplateCapacity, PowerCapacity) are required", nil)
	}
	if !req.Type.IsValid() {
		return custom_error.New(http.StatusBadRequest, "Invalid DERType. Allowed values: solar, battery, ev", nil)
	}

	err := h.Repo.CreateDERMetadata(r.Context(), &models.DERMetadata{
		ID:                uuid.New().String(),
		ProjectID:         req.ProjectID,
		Type:              req.Type,
		NameplateCapacity: req.NameplateCapacity,
		PowerCapacity:     req.PowerCapacity,
	})

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *derMetadataHandlers) GetDERMetadataHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if id == "" {
		return custom_error.New(http.StatusBadRequest, "der ID not given", nil)
	}
	contract, err := h.Repo.GetDERMetadata(r.Context(), id)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(contract)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *derMetadataHandlers) UpdateDERMetadataHandler(w http.ResponseWriter, r *http.Request) error {
	var req models.DERMetadata
	id := chi.URLParam(r, "id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return custom_error.New(http.StatusBadRequest, "Invalid request payload", err)
	}
	if id == "" {
		return custom_error.New(http.StatusBadRequest, "Contract ID is required", nil)
	}
	if req.ProjectID != "" {
		return custom_error.New(http.StatusBadRequest, "Updating project id is not allowed", nil)
	}
	if req.ID != "" {
		return custom_error.New(http.StatusBadRequest, "Updating der id not allowed", nil)
	}
	if req.Type != "" && !req.Type.IsValid() {
		return custom_error.New(http.StatusBadRequest, "Invalid DERType. Allowed values: solar, battery, ev", nil)
	}
	err := h.Repo.UpdateDERMetadata(r.Context(), id, &models.DERMetadata{
		Type:              req.Type,
		NameplateCapacity: req.NameplateCapacity,
		PowerCapacity:     req.PowerCapacity,
	})

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (handler *derMetadataHandlers) DeleteDERMetadataHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if id == "" {
		return custom_error.New(http.StatusBadRequest, "Contract ID required", nil)
	}
	if err := handler.Repo.DeleteDERMetadata(r.Context(), id); err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
