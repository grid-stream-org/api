package config

import (
	"github.com/grid-stream-org/batcher/pkg/bqclient"
	"github.com/grid-stream-org/batcher/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Port           int      `envconfig:"PORT"` // port api will run on
	AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS"`
	Database       *bqclient.Config
	Logger         *logger.Config
	Firebase       *FirebaseConfig
}

type FirebaseConfig struct {
	ProjectID        string `envconfig:"firebase_project_id"`
	GoogleCredential string `envconfig:"firebase_google_credential"`
}

func Load(skipDotenv bool) (*Config, error) {
	var cfg Config

	if !skipDotenv {
		if err := godotenv.Load(); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.WithStack(err)
	}

	return &cfg, nil
}
