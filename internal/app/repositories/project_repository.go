package repositories

// handles database interactions for Projects
// for now we define a project as an instance of offloading power from a DER to the grid
// projectId	        STRING(REQUIRED)        - Home/Building ID
// utilityId            STRING(REQUIRED)        - Utility Company ID ex: NB Power vs SJ Energy
// connectionStartAt    TIMESTAMP(Required)     - Start of offloading

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/grid-stream-org/api/internal/errors"
	"github.com/grid-stream-org/api/internal/models"
	"google.golang.org/api/iterator"
)

type ProjectRepository struct {
	client *bigquery.Client
}

func NewProjectRepository(client *bigquery.Client) *ProjectRepository {
	return &ProjectRepository{client: client}
}

// create new instance of oflloading
func (r *ProjectRepository) CreateProject(ctx context.Context, post *models.Project) error {
	// Use BigQuery client to insert a new project
	// Example: Use the client to run a query or insert data
	// INSERT INTO A1.Project (projectId, utilityId, connectionStartAt) VALUES ('projId','utilId', '2021-01-26 16:50:03' )
    return errors.New(http.StatusNotImplemented, "not yet implemented")
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, post *models.Project) error {
	// Use BigQuery client to insert a new project
	// Example: Use the client to run a query or insert data
	// INSERT INTO A1.Project (projectId, utilityId, connectionStartAt) VALUES ('projId','utilId', '2021-01-26 16:50:03' )
    return errors.New(http.StatusNotImplemented, "not yet implemented")
}

func (r *ProjectRepository) GetProject(ctx context.Context, id string) (*models.Project, error) {
	// Use BigQuery client to retrieve a project by ID
	query := `
     SELECT projectId, utilityId, connectionStartAt
     FROM A1.Project
     WHERE projectId = @projectId`

	// Create a query and set parameters
	q := r.client.Query(query)
	q.Parameters = []bigquery.QueryParameter{
		{
			Name:  "projectId",
			Value: id,
		},
	}

	// Run the query
	it, err := q.Read(ctx)
	if err != nil {
		return nil, errors.New(http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve project with ID %s", id))
	}

	var project models.Project
	// We only expect one row, so we can use Next once
	var row map[string]bigquery.Value
	err = it.Next(&row)
	if err == iterator.Done {
		return nil, errors.New(http.StatusNotFound, fmt.Sprintf("Project ID %s Not Found", id))
	}

	// Map the result to the project struct
	project = models.Project{
		ProjectID:         row["projectId"].(string),
		UtilityID:         row["utilityId"].(string),
		ConnectionStartAt: row["connectionStartAt"].(time.Time),
	}
	return &project, nil
}
