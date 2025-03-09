package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/grid-stream-org/api/internal/app/handlers/v1"
	"github.com/grid-stream-org/api/internal/custom_error"
	"github.com/grid-stream-org/api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type MockContractRepository struct {
	mock.Mock
}

func (m *MockContractRepository) CreateContract(ctx context.Context, contract *models.Contract) error {
	args := m.Called(ctx, contract)
	return args.Error(0)
}

func (m *MockContractRepository) GetContract(ctx context.Context, id string) (*models.Contract, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Contract), args.Error(1)
}

func (m *MockContractRepository) UpdateContract(ctx context.Context, id string, contract *models.Contract) error {
	args := m.Called(ctx, id, contract)
	return args.Error(0)
}

func (m *MockContractRepository) DeleteContract(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockContractRepository) GetContractsByProjectID(ctx context.Context, id string) ([]models.Contract, error) {
	args := m.Called(ctx, id)
	var c []models.Contract
	return c, args.Error(1)
}

// convert time.Time to bigquery.NullDate
func toNullDate(t time.Time) bigquery.NullDate {
	return bigquery.NullDate{
		Valid: true,
		Date:  civil.Date{Year: t.Year(), Month: t.Month(), Day: t.Day()},
	}
}

func TestCreateContractHandler(t *testing.T) {
	mockRepo := new(MockContractRepository)
	handler := handlers.NewContractHandlers(mockRepo, nil)

	startDate := toNullDate(time.Now())
	endDate := toNullDate(time.Now().AddDate(1, 0, 0)) // 1 year later

	// Test cases
	tests := []struct {
		name           string
		requestBody    models.Contract
		mockReturnErr  error
		expectedStatus int
		expectRepoCall bool
	}{
		{
			name: "Success - Valid Input",
			requestBody: models.Contract{
				ContractThreshold: 100,
				ProjectID:         "proj-123",
				StartDate:         startDate,
				EndDate:           endDate,
				Status:            "active",
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusCreated,
			expectRepoCall: true,
		},
		{
			name: "Fail - Invalid Status",
			requestBody: models.Contract{
				ContractThreshold: 100,
				ProjectID:         "proj-123",
				StartDate:         startDate,
				EndDate:           endDate,
				Status:            "invalid_status",
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
			expectRepoCall: false,
		},
		{
			name: "Fail - Missing Required Fields",
			requestBody: models.Contract{
				ProjectID: "proj-123",
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
			expectRepoCall: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.requestBody)
			req, err := http.NewRequest("POST", "/contracts", bytes.NewBuffer(body))
			assert.NoError(t, err)
			rec := httptest.NewRecorder()

			// Mock only if repo should be called
			if tc.expectRepoCall {
				mockRepo.On("CreateContract", mock.Anything, mock.Anything).Return(tc.mockReturnErr).Once()
			}

			err = handler.CreateContractHandler(rec, req)

			// Check expected status
			if err != nil {
				customErr, ok := err.(*custom_error.CustomError)
				if ok {
					assert.Equal(t, tc.expectedStatus, customErr.Code)
				} else {
					t.Fatalf("expected custom error, got: %v", err)
				}
			} else {
				assert.Equal(t, tc.expectedStatus, rec.Code)
			}

			// Ensure CreateContract is only called when expected
			if tc.expectRepoCall {
				mockRepo.AssertCalled(t, "CreateContract", mock.Anything, mock.Anything)
			} else {
				mockRepo.AssertNotCalled(t, "CreateContract")
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
