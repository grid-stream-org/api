package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/grid-stream-org/api/internal/app/repositories"
)

type DERHandler interface {
	GetDERDataByProjectIDHandler(w http.ResponseWriter, r *http.Request) error
}

type derHandler struct {
	Repo repositories.DERRepository
	Log  *slog.Logger
}

func NewDERHandlers(repo repositories.DERRepository, log *slog.Logger) DERHandler {
	return &derHandler{Repo: repo, Log: log}
}

func (h *derHandler) GetDERDataByProjectIDHandler(w http.ResponseWriter, r *http.Request) error {
	projectID := chi.URLParam(r, "project_id")
	derData, err := h.Repo.GetDERDataByProjectID(r.Context(), projectID)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(derData)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}