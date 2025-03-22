package firebase

import (
	"context"
	"log/slog"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type FirebaseConfig struct {
	ProjectID        string `envconfig:"FIREBASE_PROJECT_ID" koanf:"project_id"`
	GoogleCredential string `envconfig:"FIREBASE_GOOGLE_CREDENTIAL" koanf:"google_credential"`
}

type FirebaseClient interface {
	Auth() *auth.Client
	Firestore() *firestore.Client
	Close() error
}

type firebaseClient struct {
	authClient      *auth.Client
	firestoreClient *firestore.Client
}

func (f *firebaseClient) Auth() *auth.Client {
	return f.authClient
}

func (f *firebaseClient) Firestore() *firestore.Client {
	return f.firestoreClient
}

func (f *firebaseClient) Close() error {
	return f.firestoreClient.Close()
}

// NewFirebaseClient initializes Firebase Auth & Firestore clients
func NewFirebaseClient(ctx context.Context, cfg *FirebaseConfig, log *slog.Logger) (FirebaseClient, error) {
	opt := option.WithCredentialsFile(cfg.GoogleCredential)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &firebaseClient{
		authClient:      authClient,
		firestoreClient: firestoreClient,
	}, nil
}
