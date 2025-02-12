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

	// Check if running in production
	isProduction := os.Getenv("GO_ENV") == "production"

	// Load .env only in local development
	if !isProduction {
		_ = godotenv.Load(".env")
	}

	// Process environment variables into struct
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.WithStack(err)
	}

	// Ensure Firebase credentials file exists
	// also bypass this check if we are running unit tests
	if os.Getenv("TEST_ENV") != "true" {
		if _, err := os.Stat(cfg.Firebase.GoogleCredential); os.IsNotExist(err) {
			return nil, errors.WithStack(fmt.Errorf("Firebase credentials file not found: %s", cfg.Firebase.GoogleCredential))
		}
	}

	return &cfg, nil
}
