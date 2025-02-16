package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/grid-stream-org/api/internal/app/repositories"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
)

type DERMetadataHandlers interface {
	CreateDERMetadataHandler(w http.ResponseWriter, r *http.Request) error
	BatchCreateDERMetadataHandler(w http.ResponseWriter, r *http.Request) error
	GetDERMetadataHandler(w http.ResponseWriter, r *http.Request) error
	ListDERMetadataByProjectHandler(w http.ResponseWriter, r *http.Request) error
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

	// Power capacity can be entered by the user at a later date as it changes
	if req.ID == "" || req.Type == "" || req.ProjectID == "" || req.NameplateCapacity <= 0 {
		return custom_error.New(http.StatusBadRequest, "All fields (ID, Type, ProjectID, Type, NameplateCapacity) are required", nil)
	}
	if !req.Type.IsValid() {
		return custom_error.New(http.StatusBadRequest, "Invalid DERType. Allowed values: solar, battery, ev", nil)
	}

	err := h.Repo.CreateDERMetadata(r.Context(), &req)

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *derMetadataHandlers) BatchCreateDERMetadataHandler(w http.ResponseWriter, r *http.Request) error {
	var req []models.DERMetadata
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return custom_error.New(http.StatusBadRequest, "Invalid request payload", err)
	}

	for _, der := range req {
		// Power capacity can be entered by the user at a later date as it changes
		if der.Type == "" || der.ProjectID == "" || der.NameplateCapacity <= 0 {
			return custom_error.New(http.StatusBadRequest, "All fields (Type, ProjectID, Type, NameplateCapacity) are required", nil)
		}
		if !der.Type.IsValid() {
			return custom_error.New(http.StatusBadRequest, "Invalid DERType. Allowed values: solar, battery, ev", nil)
		}
	}

	err := h.Repo.BatchCreateDERMetadata(r.Context(), req)

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

func (h *derMetadataHandlers) ListDERMetadataByProjectHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("project_id")
	if id == "" {
		return custom_error.New(http.StatusBadRequest, "project ID request", nil)
	}

	ders, err := h.Repo.ListDERMetadataByProject(r.Context(), id)
	if err != nil {
		return err
	}

	if ders == nil {
		ders = []models.DERMetadata{}
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(ders)
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
