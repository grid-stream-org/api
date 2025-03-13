package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/grid-stream-org/api/internal/app/repositories"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
)

type NotificationHandler interface {
	NotifyUserHandler(w http.ResponseWriter, r *http.Request) error
}

type notificationtHandler struct {
	repo repositories.NotificationRepository
	log  *slog.Logger
}

func NewNotificationHandler(r repositories.NotificationRepository, log *slog.Logger) NotificationHandler {
	return &notificationtHandler{
		repo: r,
		log:  log,
	}
}

func (h *notificationtHandler) NotifyUserHandler(w http.ResponseWriter, r *http.Request) error {
	var req models.FaultNotification
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return custom_error.New(http.StatusBadRequest, "Invalid request payload", err)
	}
	if req.ProjectID == "" {
		return custom_error.New(http.StatusBadRequest, "Project ID must not be empty", nil)
	}
    
	// TODO: validate dates and maybe messages
	if err := h.repo.NotifyUser(r.Context(), &req); err != nil {
		return err
	}
    w.WriteHeader(http.StatusOK)
	return nil
}
