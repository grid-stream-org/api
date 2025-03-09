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

type ContractHandler interface {
	CreateContractHandler(w http.ResponseWriter, r *http.Request) error
	GetContractHandler(w http.ResponseWriter, r *http.Request) error
	UpdateContractHandler(w http.ResponseWriter, r *http.Request) error
	DeleteContractHandler(w http.ResponseWriter, r *http.Request) error
    GetContractsByProjectIDHandler(w http.ResponseWriter, r *http.Request) error
}

type contractHandler struct {
	Repo repositories.ContractRepository
	Log  *slog.Logger
}

func NewContractHandlers(repo repositories.ContractRepository, log *slog.Logger) ContractHandler {
	return &contractHandler{Repo: repo, Log: log}
}

func (h *contractHandler) CreateContractHandler(w http.ResponseWriter, r *http.Request) error {
	var req models.Contract
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return custom_error.New(http.StatusBadRequest, "Invalid request payload", err)
	}
	// TODO: add better start date and end date validation
	if req.ContractThreshold <= 0 || req.ProjectID == "" || !req.EndDate.Valid || !req.StartDate.Valid || req.Status == "" {
		return custom_error.New(http.StatusBadRequest, "All fields (contract threshold, start date, end date, projectID, status) are required", nil)
	}
	if !req.Status.IsValid() {
		return custom_error.New(http.StatusBadRequest, "Invalid status must be active, inactive or pending", nil)
	}
	err := h.Repo.CreateContract(r.Context(), &models.Contract{
		ID:                uuid.New().String(),
		ContractThreshold: req.ContractThreshold,
		StartDate:         req.StartDate,
		EndDate:           req.EndDate,
		ProjectID:         req.ProjectID,
		Status:            req.Status,
	})
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *contractHandler) GetContractHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if id == "" {
		return custom_error.New(http.StatusBadRequest, "Contract ID not given", nil)
	}
	contract, err := h.Repo.GetContract(r.Context(), id)
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

func (h *contractHandler) UpdateContractHandler(w http.ResponseWriter, r *http.Request) error {
	var req models.Contract
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
		return custom_error.New(http.StatusBadRequest, "Updating contract id not allowed", nil)
	}
	if req.Status != "" && !req.Status.IsValid() {
		return custom_error.New(http.StatusBadRequest, "Invalid status must be active, inactive or pending", nil)
	}
	err := h.Repo.UpdateContract(r.Context(), id, &models.Contract{
		ContractThreshold: req.ContractThreshold,
		StartDate:         req.StartDate,
		EndDate:           req.EndDate,
		Status:            req.Status,
	})

	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *contractHandler) DeleteContractHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if id == "" {
		return custom_error.New(http.StatusBadRequest, "Contract ID required", nil)
	}
	if err := h.Repo.DeleteContract(r.Context(), id); err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *contractHandler) GetContractsByProjectIDHandler(w http.ResponseWriter, r *http.Request) error {
    id := chi.URLParam(r, "projectId")
    if id == "" {
		return custom_error.New(http.StatusBadRequest, "Project ID required", nil)
	}

    contracts, err := h.Repo.GetContractsByProjectID(r.Context(), id)
	if err != nil {
		return err
	}

	if contracts == nil {
		contracts = []models.Contract{}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(contracts); err != nil {
        return err
    }
    w.WriteHeader(http.StatusOK)

    return nil
}