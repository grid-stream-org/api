package config

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
    "github.com/joho/godotenv"
)

type Config struct {
	Port	    int	            `envconfig:"port"` // port api will run on
	Timeout	    time.Duration   `envconfig:"timeout"` // timeout for requests
	Database	DatabaseConfig	 `envconfig:"database"`
    Logger      LoggerConfig     `envconfig:"LOG"`  
}

type DatabaseConfig struct {
	BigQuery	BigQueryConfig `envconfig:"bigquery"`
}

type BigQueryConfig struct {
	ProjectID	string `envconfig:"project_id"`
	CredsFile	string `envconfig:"credentials"`
}

type LoggerConfig struct {
	Level   string `envconfig:"LEVEL" default:"INFO"`
	Format  string `envconfig:"FORMAT" default:"text"`
	Output  string `envconfig:"OUTPUT" default:"stdout"`
}

func Load() (*Config, error) {
    err := godotenv.Load()
    if err != nil {
        return nil, errors.New("missing .env file")
    }


	port := func() int {
		if p, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
			return p
		}
		return 8080 // default port
	}()

	timeout := func() time.Duration {
		if t, err := time.ParseDuration(os.Getenv("TIMEOUT")); err == nil {
			return t
		}
		return 10 * time.Second
	}()

	bigQueryProjectId := os.Getenv("BIGQUERY_PROJECT_ID")
	if bigQueryProjectId == "" {
		return nil, errors.New("missing big query project id")
	}

	bigQueryCredentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if bigQueryCredentialsFile == "" {
		return nil, errors.New("missing big query creds")
	}

	return &Config{
		Port:	port,
		Timeout:	timeout,
		Database: DatabaseConfig{
			BigQuery: BigQueryConfig{
				ProjectID: bigQueryProjectId,
				CredsFile: bigQueryCredentialsFile,
			},
		},
        Logger: LoggerConfig{
            Level:  os.Getenv("LOG_LEVEL"),
            Format: os.Getenv("LOG_FORMAT"),
            Output: os.Getenv("LOG_OUTPUT"),
        },
	}, nil
}