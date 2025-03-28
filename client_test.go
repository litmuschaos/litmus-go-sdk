package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/litmuschaos/litmus-go-sdk/pkg/apis"
	"github.com/litmuschaos/litmus-go-sdk/pkg/logger"
	"github.com/stretchr/testify/assert"
)

// Test configuration with defaults
var (
	testEndpoint = "http://127.0.0.1:39651"
	testUsername = "admin"
	testPassword = "LitmusChaos123@"
)

func init() {
	// Override defaults with environment variables if set
	if endpoint := os.Getenv("LITMUS_TEST_ENDPOINT"); endpoint != "" {
		testEndpoint = endpoint
	}
	if username := os.Getenv("LITMUS_TEST_USERNAME"); username != "" {
		testUsername = username
	}
	if password := os.Getenv("LITMUS_TEST_PASSWORD"); password != "" {
		testPassword = password
	}

	logger.Infof("Test configuration - Endpoint: %s, Username: %s", testEndpoint, testUsername)
}

// setupTestClient creates a test client using the configured credentials
func setupTestClient() (*LitmusClient, error) {
	return NewLitmusClient(testEndpoint, testUsername, testPassword)
}

func TestNewLitmusClient(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "successful client creation",
			endpoint: testEndpoint,
			username: testUsername,
			password: testPassword,
			wantErr:  false,
		},
		{
			name:     "invalid endpoint",
			endpoint: "invalid-url",
			username: testUsername,
			password: testPassword,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewLitmusClient(tt.endpoint, tt.username, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, client, "Client should not be nil")
		})
	}
}

func TestListProjects(t *testing.T) {
	tests := []struct {
		name        string
		setupClient func() (*LitmusClient, error)
		wantErr     bool
	}{
		{
			name: "successful projects listing",
			setupClient: func() (*LitmusClient, error) {
				return setupTestClient()
			},
			wantErr: false,
		},
		{
			name: "invalid auth token",
			setupClient: func() (*LitmusClient, error) {
				// Create client with invalid token for testing error case
				client, err := setupTestClient()
				if err != nil {
					return nil, err
				}
				// Invalidate the token
				client.credentials.Token = "invalid-token"
				return client, nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := tt.setupClient()
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			projects, err := client.ListProjects()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, projects, "Projects list should not be nil")
		})
	}
}

func TestCreateProject(t *testing.T) {
	tests := []struct {
		name        string
		setupClient func() (*LitmusClient, error)
		projectName string
		wantErr     bool
		validate    func(*testing.T, interface{}, error)
	}{
		{
			name: "successful project creation",
			setupClient: func() (*LitmusClient, error) {
				return setupTestClient()
			},
			projectName: fmt.Sprintf("test-project-%s", uuid.New().String()),
			wantErr:     false,
			validate: func(t *testing.T, result interface{}, err error) {
				project := result.(apis.CreateProjectResponse)
				assert.NotEmpty(t, project.Data.ID, "Project ID should not be empty")
				assert.NotEmpty(t, project.Data.Name, "Project name should not be empty")

				logger.InfoWithValues("Project created", map[string]interface{}{
					"name": project.Data.Name,
					"id":   project.Data.ID,
				})
			},
		},
		{
			name: "empty project name",
			setupClient: func() (*LitmusClient, error) {
				return setupTestClient()
			},
			projectName: "",
			wantErr:     true,
			validate:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := tt.setupClient()
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			project, err := client.CreateProject(tt.projectName)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, project, "Created project should not be nil")

			if tt.validate != nil {
				tt.validate(t, project, err)
			}
		})
	}
}
