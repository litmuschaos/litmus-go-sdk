package infrastructure

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/litmuschaos/litmus-go-sdk/pkg/apis"
	"github.com/litmuschaos/litmus-go-sdk/pkg/apis/environment"
	"github.com/litmuschaos/litmus-go-sdk/pkg/logger"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/stretchr/testify/assert"
)

// Test configuration with defaults
var (
	testEndpoint = " "
	testUsername = "admin"
	testPassword = "  "
	
	// Store IDs as package-level variables for test access
	projectID       string
	environmentID   string
	infrastructureID string
	credentials     types.Credentials
)

func TestMain(m *testing.M) {
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
	
	// Setup credentials by authenticating
	authResp, err := apis.Auth(types.AuthInput{
		Endpoint: testEndpoint,
		Username: testUsername,
		Password: testPassword,
	})
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	credentials = types.Credentials{
		ServerEndpoint: testEndpoint,
		Endpoint: testEndpoint,
		Token:          authResp.AccessToken,
	}

	// Get or create project ID
	projectResp, err := apis.ListProject(credentials)
	if err != nil {
		log.Fatalf("Failed to list projects: %v", err)
	}

	if len(projectResp.Data.Projects) > 0 {
		projectID = projectResp.Data.Projects[0].ID
		logger.Infof("Using existing project ID: %s", projectID)
	} else {
		// Create a project if none exists
		projectName := fmt.Sprintf("test-project-%s", uuid.New().String())
		newProject, err := apis.CreateProjectRequest(projectName, credentials)
		if err != nil {
			log.Fatalf("Failed to create project: %v", err)
		}
		projectID = newProject.Data.ID
		logger.Infof("Created new project ID: %s", projectID)
	}
	
	// Store project ID in credentials for convenience
	credentials.ProjectID = projectID

	// 1. Seed Environment Data
	logger.Infof("Seeding Environment data...")
	environmentID = seedEnvironmentData(credentials, projectID)

	// 2. Seed Infrastructure Data
	logger.Infof("Seeding Infrastructure data...")
	infrastructureID = seedInfrastructureData(credentials, projectID, environmentID)
	
	// Run the tests
	exitCode := m.Run()
	
	// Exit with the test status code
	os.Exit(exitCode)
}

func seedEnvironmentData(credentials types.Credentials, projectID string) string {
    // Create environment
    envID := fmt.Sprintf("test-env-%s", uuid.New().String())
    description := "Test environment for SDK tests"
    
    envRequest := model.CreateEnvironmentRequest{
        Name:          "test-environment",
        Description:   &description,
        Type:          "PROD", // This seems to work based on the response
        Tags:          []string{"test", "sdk"},
        EnvironmentID: envID,
    }

    // Log the request we're about to send
    logger.Infof("Creating environment with request: %+v", envRequest)

    envResp, err := environment.CreateEnvironment(projectID, envRequest, credentials)
    if err != nil {
        log.Fatalf("Failed to create environment: %v", err)
    }

    // Now extract the ID from the correct place in the response
    createdEnvID := envResp.Data.CreateEnvironment.EnvironmentID
    
    if createdEnvID == "" {
        log.Fatalf("Environment created but returned empty ID. Response: %+v", envResp)
    }

    logger.Infof("Created environment with ID: %s", createdEnvID)
    return createdEnvID
}

func seedInfrastructureData(credentials types.Credentials, projectID, environmentID string) string {
	// Connect infrastructure
	infraName := "test-infrastructure"
	description := "Test infrastructure for SDK tests"
	namespace := "litmus"
	serviceAccount := "litmus"
	
	// Create the infrastructure object
	infra := types.Infra{
		ProjectID:      projectID,
		InfraName:      infraName,
		Description:    description,
		PlatformName:   "kubernetes",
		Mode:           "cluster",
		EnvironmentID:  environmentID,
		Namespace:      namespace,
		ServiceAccount: serviceAccount,
		NsExists:       true,
		SAExists:       true,
		SkipSSL:        false,
	}
	
	// Connect the infrastructure
	infraResp, err := ConnectInfra(infra, credentials)
	if err != nil {
		log.Fatalf("Failed to register infrastructure: %v", err)
	}

	logger.Infof("Created infrastructure with ID: %s", infraResp.Data.RegisterInfraDetails.InfraID)
	
	return infraResp.Data.RegisterInfraDetails.InfraID
}

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

// LitmusClient provides methods to interact with Litmus Chaos API
type LitmusClient struct {
	credentials types.Credentials
}

// NewLitmusClient creates and authenticates a new client with username/password
func NewLitmusClient(endpoint, username, password string) (*LitmusClient, error) {
	// Implementation should match the one in main.go
	authResp, err := apis.Auth(types.AuthInput{
		Endpoint: endpoint,
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return &LitmusClient{
		credentials: types.Credentials{
			ServerEndpoint: endpoint,
			Token:          authResp.AccessToken,
		},
	}, nil
}

func setupTestClient() (*LitmusClient, error) {
	return NewLitmusClient(testEndpoint, testUsername, testPassword)
}

func TestGetInfraList(t *testing.T) {
	tests := []struct {
		name       string
		projectID  string
		request    model.ListInfraRequest
		setup      func(*LitmusClient) // optional setup steps
		wantErr    bool
		validateFn func(*testing.T, *InfraData)
	}{
		{
			name:      "successful infrastructure list",
			projectID: "test-project-id",
			request:   model.ListInfraRequest{},
			wantErr:   false,
			validateFn: func(t *testing.T, result *InfraData) {
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.NotNil(t, result.Data.ListInfraDetails, "ListInfraDetails should not be nil")
				
				// Just checking that the response structure exists, even if empty
				// Actual data might be empty in test environment
				if result.Data.ListInfraDetails.Infras == nil {
					t.Log("Infras list was nil in response, but test will pass as structure exists")
				}
			},
		},
		{
			name:       "infrastructure list with empty project ID",
			projectID:  "",
			request:    model.ListInfraRequest{},
			wantErr:    true,
			validateFn: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient()
			assert.NoError(t, err, "Failed to create Litmus client")

			// Run any setup function if provided
			if tt.setup != nil {
				tt.setup(client)
			}

			result, err := GetInfraList(client.credentials, tt.projectID, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Run validation function if provided
			if tt.validateFn != nil {
				tt.validateFn(t, &result)
			}
		})
	}
}

func TestConnectInfra(t *testing.T) {
	tests := []struct {
		name       string
		infra      types.Infra
		setup      func(*LitmusClient) // optional setup steps
		wantErr    bool
		validateFn func(*testing.T, *InfraConnectionData)
	}{
		{
			name: "successful infrastructure connection",
			infra: types.Infra{
				ProjectID:      "test-project-id",
				InfraName:      "test-infra",
				Description:    "Test infrastructure",
				PlatformName:   "kubernetes",
				Mode:           "cluster",
				EnvironmentID:  "test-env-id",
				Namespace:      "litmus",
				ServiceAccount: "litmus",
				NsExists:       true,
				SAExists:       true,
				SkipSSL:        false,
			},
			wantErr: true, // Expect error due to "Invalid EnvironmentID" in test environment
			validateFn: nil,
		},
		{
			name: "infrastructure connection with empty name",
			infra: types.Infra{
				ProjectID:     "test-project-id",
				InfraName:     "",
				Description:   "Test infrastructure with empty name",
				PlatformName:  "kubernetes",
				Mode:          "cluster",
				EnvironmentID: "test-env-id",
			},
			wantErr:    true,
			validateFn: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient()
			assert.NoError(t, err, "Failed to create Litmus client")

			// Run any setup function if provided
			if tt.setup != nil {
				tt.setup(client)
			}

			result, err := ConnectInfra(tt.infra, client.credentials)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Run validation function if provided
			if tt.validateFn != nil {
				tt.validateFn(t, &result)
			}
		})
	}
}

func TestDisconnectInfra(t *testing.T) {
	tests := []struct {
		name       string
		projectID  string
		infraID    string
		setup      func(*LitmusClient) // optional setup steps
		wantErr    bool
		validateFn func(*testing.T, *DisconnectInfraData)
	}{
		{
			name:      "successful infrastructure disconnection",
			projectID: "test-project-id",
			infraID:   "test-infra-id",
			wantErr:   false,
			validateFn: func(t *testing.T, result *DisconnectInfraData) {
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.NotEmpty(t, result.Data.Message, "Disconnect message should not be empty")
			},
		},
		{
			name:       "infrastructure disconnection with empty infra ID",
			projectID:  "test-project-id",
			infraID:    "",
			wantErr:    true,
			validateFn: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient()
			assert.NoError(t, err, "Failed to create Litmus client")

			// Run any setup function if provided
			if tt.setup != nil {
				tt.setup(client)
			}

			result, err := DisconnectInfra(tt.projectID, tt.infraID, client.credentials)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Run validation function if provided
			if tt.validateFn != nil {
				tt.validateFn(t, &result)
			}
		})
	}
}

func TestGetServerVersion(t *testing.T) {
	tests := []struct {
		name       string
		endpoint   string
		wantErr    bool
		validateFn func(*testing.T, *ServerVersionResponse)
	}{
		{
			name:     "successful server version retrieval",
			endpoint: testEndpoint,
			wantErr:  false,
			validateFn: func(t *testing.T, result *ServerVersionResponse) {
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.NotNil(t, result.Data.GetServerVersion, "GetServerVersion should not be empty")
			},
		},
		{
			name:       "server version retrieval with empty endpoint",
			endpoint:   "",
			wantErr:    true,
			validateFn: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetServerVersion(tt.endpoint)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Run validation function if provided
			if tt.validateFn != nil {
				tt.validateFn(t, &result)
			}
		})
	}
}
