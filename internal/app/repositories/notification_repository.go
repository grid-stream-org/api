package repositories

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
	"github.com/grid-stream-org/api/pkg/firebase"
)

type NotificationRepository interface {
	NotifyUser(ctx context.Context, data *models.FaultNotification) error
}

type notificationRepository struct {
	fb  firebase.FirebaseClient
	log *slog.Logger
}

func NewNotificationRepository(fb firebase.FirebaseClient, log *slog.Logger) NotificationRepository {
	return &notificationRepository{
		fb:  fb,
		log: log,
	}
}

func (r *notificationRepository) NotifyUser(ctx context.Context, data *models.FaultNotification) error {
	firestore := r.fb.Firestore()

	_, _, err := firestore.Collection("notifications").Add(ctx, data)
	if err != nil {
		return custom_error.New(http.StatusInternalServerError, "Failed to send notification", err)
	}

	return nil
}
