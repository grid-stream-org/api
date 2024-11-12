package repositories

// handles database interactions with projects
// for now we define a project as an instance of offloading power from a DER to the grid
// projectId	        STRING(REQUIRED)        - Home/Building ID	
// utilityId            STRING(REQUIRED)        - Utility Company ID ex: NB Power vs SJ Energy
// connectionStartAt    TIMESTAMP(Required)     - Start of offloading


import (
	"context"
	"cloud.google.com/go/bigquery"
	
   // "github.com/grid-stream-org/api/pkg/database"
	"github.com/grid-stream-org/api/internal/models"
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
	return nil
}

func (r *ProjectRepository) GetProject(ctx context.Context, id string) (*models.Project, error) {
	// Use BigQuery client to retrieve a project by ID
	return nil, nil
}

