package logic

import (
	"reflect"
	"testing"

	"github.com/grid-stream-org/api/internal/models"
)

// tests that GetUpdates correctly extracts non-empty fields from a struct for extracting fields from body of requests
func TestGenerateUpdatesMap(t *testing.T) {
	// Test Case 1: Partially filled Project struct
	project := models.Project{
		ID:        "proj-123",
		UtilityID: "util-456",
		Location:  "New York",
	}

	updates := ExtractBody(project)

	expected := map[string]any{
		"id":         "proj-123",
		"utility_id": "util-456",
		"location":   "New York",
	}

	if !reflect.DeepEqual(updates, expected) {
		t.Errorf("Expected %v, got %v", expected, updates)
	}

	// Test Case 2: Empty Project struct should return an empty map
	emptyProject := models.Project{}
	updates = ExtractBody(emptyProject)

	if len(updates) != 0 {
		t.Errorf("Expected empty map, got %v", updates)
	}

	// Test Case 3: Project with only UserID set
	partialProject := models.Project{
		UserID: "user-789",
	}

	expectedPartial := map[string]any{
		"user_id": "user-789",
	}

	updates = ExtractBody(partialProject)

	if !reflect.DeepEqual(updates, expectedPartial) {
		t.Errorf("Expected %v, got %v", expectedPartial, updates)
	}
}
