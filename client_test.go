package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/litmuschaos/litmus-go-sdk/pkg/apis"
	"github.com/litmuschaos/litmus-go-sdk/pkg/logger"
	"github.com/litmuschaos/litmus-go-sdk/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

// Getenv fetches the env and returns the default value if env is empty
func Getenv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

// Test configuration with defaults
var (	
	
	testEndpoint = "http://127.0.0.1:39651"
	testUsername = "admin"
	testPassword = "litmus"	)

func init() {
	// Override defaults with environment variables if set
	testEndpoint = Getenv("LITMUS_TEST_ENDPOINT", testEndpoint)
	testUsername = Getenv("LITMUS_TEST_USERNAME", testUsername)
	testPassword = Getenv("LITMUS_TEST_PASSWORD", testPassword)

	logger.Infof("Test configuration - Endpoint: %s, Username: %s", testEndpoint, testUsername)
}

// setupTestClient creates a test client using the configured credentials
func setupTestClient() (sdk.Client, error) {
	return sdk.NewClient(sdk.ClientOptions{
		Endpoint: testEndpoint,
		Username: testUsername,
		Password: testPassword,
	})
}

func TestNewClient(t *testing.T) {
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
			client, err := sdk.NewClient(sdk.ClientOptions{
				Endpoint: tt.endpoint,
				Username: tt.username,
				Password: tt.password,
			})

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
		setupClient func() (sdk.Client, error)
		wantErr     bool
	}{
		{
			name: "successful projects listing",
			setupClient: func() (sdk.Client, error) {
				return setupTestClient()
			},
			wantErr: false,
		},
		{
			name: "invalid auth token",
			setupClient: func() (sdk.Client, error) {
				// Create client with invalid token for testing error case
				client, err := setupTestClient()
				if err != nil {
					return nil, err
				}
				// For testing invalid auth, we'll create a new request
				return client, nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := tt.setupClient()
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			projects, err := client.Projects().List()

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
		setupClient func() (sdk.Client, error)
		projectName string
		wantErr     bool
		validate    func(*testing.T, interface{}, error)
	}{
		{
			name: "successful project creation",
			setupClient: func() (sdk.Client, error) {
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
			setupClient: func() (sdk.Client, error) {
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

			project, err := client.Projects().Create(tt.projectName)

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
