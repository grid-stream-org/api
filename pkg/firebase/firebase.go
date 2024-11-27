package firebase

import (
	"context"
	"log/slog"

	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/grid-stream-org/api/internal/config"
	"google.golang.org/api/option"
)

var FirebaseAuthClient *auth.Client

// InitializeFirebaseClient sets up the Firebase Auth client
func InitializeFirebaseClient(ctx context.Context, cfg *config.Config, log *slog.Logger) (*auth.Client, error) {
	// Create an option with the credentials file from your config
	opt := option.WithCredentialsFile(cfg.Firebase.GoogleCredential)

	// Initialize the Firebase App
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	// Initialize the Firebase Auth client
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
    
	FirebaseAuthClient = client
	log.Info("Firebase Auth client initialized successfully")
	return FirebaseAuthClient, nil
}
