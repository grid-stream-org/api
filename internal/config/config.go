package config

import (
	"fmt"
	"os"

	"github.com/grid-stream-org/batcher/pkg/bqclient"
	"github.com/grid-stream-org/batcher/pkg/logger"
	"github.com/joho/godotenv" // Still needed for local
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Port           int      `envconfig:"PORT"` // API port
	AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS"`
	Database       *bqclient.Config
	Logger         *logger.Config
	Firebase       *FirebaseConfig
}

type FirebaseConfig struct {
	ProjectID        string `envconfig:"FIREBASE_PROJECT_ID"`
	GoogleCredential string `envconfig:"FIREBASE_GOOGLE_CREDENTIAL"`
}

func Load() (*Config, error) {
	var cfg Config

	// Load .env only in local development (not in Cloud Run)
	if os.Getenv("GO_ENV") != "production" {
		_ = godotenv.Load(".env") // Ignore error if .env is missing
	}

	// Process environment variables into struct
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.WithStack(err)
	}

	// Ensure Firebase credentials file exists
	if _, err := os.Stat(cfg.Firebase.GoogleCredential); os.IsNotExist(err) {
		return nil, errors.WithStack(fmt.Errorf("Firebase credentials file not found: %s", cfg.Firebase.GoogleCredential))
	}

	return &cfg, nil
}
