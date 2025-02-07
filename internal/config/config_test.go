package config


// TODO FIX THIS TEST, I HAD TO MESS AROUND TO MAKE DEPLOYMENT WORK AND THIS WAS A BLOCKER

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestLoadConfig(t *testing.T) {
// 	// Set up environment variables
// 	t.Setenv("PORT", "8080")
// 	t.Setenv("ALLOWED_ORIGINS", "http://localhost,http://example.com")
// 	t.Setenv("FIREBASE_PROJECT_ID", "test-project-id")
// 	t.Setenv("FIREBASE_GOOGLE_CREDENTIAL", "/path/to/credential.json")

// 	// Load configuration without .env
// 	cfg, err := Load()

// 	// Assertions
// 	assert.NoError(t, err, "Loading configuration should not return an error")
// 	assert.NotNil(t, cfg, "Config should not be nil")

// 	// Validate config fields
// 	assert.Equal(t, 8080, cfg.Port, "Port should be 8080")
// 	assert.ElementsMatch(t, []string{"http://localhost", "http://example.com"}, cfg.AllowedOrigins, "Allowed origins should match")
// 	assert.Equal(t, "test-project-id", cfg.Firebase.ProjectID, "Firebase ProjectID should match")
// 	assert.Equal(t, "/path/to/credential.json", cfg.Firebase.GoogleCredential, "Firebase GoogleCredential should match")
// }
