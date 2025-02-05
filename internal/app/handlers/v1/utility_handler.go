package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/grid-stream-org/api/internal/app/repositories"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
)

type UtilityHandler interface {
	CreateUtilityHandler(w http.ResponseWriter, r *http.Request) error
	GetUtilityHandler(w http.ResponseWriter, r *http.Request) error
	UpdateUtilityHandler(w http.ResponseWriter, r *http.Request) error
	DeleteUtilityHandler(w http.ResponseWriter, r *http.Request) error
}

type utilityHandler struct {
	Repo   repositories.UtilityRepository
	logger *slog.Logger
}

func NewUtilityRepository(repo repositories.UtilityRepository, logger *slog.Logger) UtilityHandler {
	return &utilityHandler{
		Repo:   repo,
		logger: logger,
	}
}

func (handler *utilityHandler) GetUtilityHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if id == "" {
		return custom_error.New(http.StatusBadRequest, "ID is required", errors.New("ID is required"))
	}

	var util *models.Utility

	util, err := handler.Repo.GetUtility(r.Context(), id)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(util)
	if err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to encode utility into json", err)
	}
	w.WriteHeader(http.StatusOK)

	return nil
}

func (handler *utilityHandler) CreateUtilityHandler(w http.ResponseWriter, r *http.Request) error {
	var req models.Utility
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return custom_error.New(http.StatusBadRequest, "Invalid request payload", err)
	}
	if req.DisplayName == "" {
		return custom_error.New(http.StatusBadRequest, "Display name required", errors.New("Display name not provided"))
	}

	if err := handler.Repo.CreateUtility(r.Context(), &models.Utility{DisplayName: req.DisplayName}); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func (handler *utilityHandler) UpdateUtilityHandler(w http.ResponseWriter, r *http.Request) error {
	var req models.Utility
	id := chi.URLParam(r, "id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return custom_error.New(http.StatusBadRequest, "Invalid request payload", err)
	}
	if req.DisplayName == "" {
		return custom_error.New(http.StatusBadRequest, "No fields to update", nil)
	}
	if req.ID != "" {
		return custom_error.New(http.StatusBadRequest, "Not permitted to update Utility ID", nil)
	}
	if id == "" {
		return custom_error.New(http.StatusBadRequest, "Utlity ID required", nil)
	}
	if err := handler.Repo.UpdateUtility(r.Context(), id, &models.Utility{DisplayName: req.DisplayName}); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

func (handler *utilityHandler) DeleteUtilityHandler(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if id == "" {
		return custom_error.New(http.StatusBadRequest, "Utlity ID required", errors.New("Utlity ID required"))
	}
	if err := handler.Repo.DeleteUtility(r.Context(), id); err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
