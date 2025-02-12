package firebase

import (
	"context"
	"log/slog"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/grid-stream-org/api/internal/config"
	"google.golang.org/api/option"
)

var FirebaseAuthClient *auth.Client

type FirebaseClient interface {
	Auth() *auth.Client
	Firestore() *firestore.Client
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

// InitializeFirebaseClient sets up the Firebase Auth client and Firestore client
func NewFirebaseClient(ctx context.Context, cfg *config.Config, log *slog.Logger) (FirebaseClient, error) {
	opt := option.WithCredentialsFile(cfg.Firebase.GoogleCredential)
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

	// Construct and return the interface
	return &firebaseClient{
		authClient:      authClient,
		firestoreClient: firestoreClient,
	}, nil
}
