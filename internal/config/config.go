package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	Port	    int	            `envconfig:"port"` // port api will run on
	Timeout	    time.Duration   `envconfig:"timeout"` // timeout for requests
	Database	DatabaseConfig	 `envconfig:"database"`
    Logger      LoggerConfig     `envconfig:"LOG"`  
    Firebase    FirebaseConfig   `envconfig:"auth"`
}

type DatabaseConfig struct {
	BigQuery	BigQueryConfig `envconfig:"bigquery"`
}

type BigQueryConfig struct {
	ProjectID	string `envconfig:"project_id"`
	CredsFile	string `envconfig:"credentials"`
}

type FirebaseConfig struct {
    ProjectID        string `envconfig:"firebase_project_id"`
    GoogleCredential string `envconfig:"firebase_google_credential"`
}

type LoggerConfig struct {
	Level   string `envconfig:"LEVEL" default:"INFO"`
	Format  string `envconfig:"FORMAT" default:"text"`
	Output  string `envconfig:"OUTPUT" default:"stdout"`
}

func Load() (*Config, error) {
    // Don't try to load .env file if we are in a test environment
    // TODO maybe there's a better way to get around this, should look into a better way of loading environment variables
    if os.Getenv("GO_ENV") != "test" {
        err := godotenv.Load()
        if err != nil {
            return nil, errors.New("missing .env file and required environment variables are not set")
        }
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

    firebaseProjectID := os.Getenv("FIREBASE_PROJECT_ID")
    if firebaseProjectID == "" {
        return nil, errors.New("missing Firebase project ID")
    }

    firebaseGoogleCredential := os.Getenv("GOOGLE_CLOUD_PROJECT")
    if firebaseGoogleCredential == "" {
        return nil, errors.New("missing Firebase Google credential file")
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
        Firebase: FirebaseConfig{
            ProjectID: firebaseProjectID,
            GoogleCredential: firebaseGoogleCredential,
        },
	}, nil
}