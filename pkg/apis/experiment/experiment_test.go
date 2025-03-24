package experiment

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/litmuschaos/litmus-go-sdk/pkg/apis"
	"github.com/litmuschaos/litmus-go-sdk/pkg/logger"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/stretchr/testify/assert"
)

// Test configuration with defaults
var (
	testEndpoint = "http://127.0.0.1:35961"
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

// LitmusClient provides methods to interact with Litmus Chaos API
type LitmusClient struct {
	credentials types.Credentials
}

func setupTestClient() (*LitmusClient, error) {
	return NewLitmusClient(testEndpoint, testUsername, testPassword)
}

func TestSaveExperiment(t *testing.T) {
	tests := []struct {
		name       string
		projectID  string
		request    model.SaveChaosExperimentRequest
		setup      func(*LitmusClient) // optional setup steps
		wantErr    bool
		validateFn func(*testing.T, *SaveExperimentData)
	}{
		{
			name:      "successful experiment save",
			projectID: "test-project-id",
			request: model.SaveChaosExperimentRequest{
				ID:   "test-experiment-id",
				Name: "test-experiment",
			},
			wantErr: false,
			validateFn: func(t *testing.T, result *SaveExperimentData) {
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				// Check that experiment ID is present in response message
				assert.Contains(t, result.Data.Message, "test-experiment-id",
					"Response should contain experiment ID")
			},
		},
		{
			name:      "save experiment with empty ID",
			projectID: "test-project-id",
			request: model.SaveChaosExperimentRequest{
				ID:   "",
				Name: "test-experiment",
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

			result, err := SaveExperiment(tt.projectID, tt.request, client.credentials)

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

func TestRunExperiment(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		experimentID string
		setup        func(*LitmusClient) // optional setup steps
		wantErr      bool
		validateFn   func(*testing.T, *RunExperimentResponse)
	}{
		{
			name:         "successful experiment run",
			projectID:    "test-project-id",
			experimentID: "test-experiment-id",
			wantErr:      false,
			validateFn: func(t *testing.T, result *RunExperimentResponse) {
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.NotEmpty(t, result.Data.RunExperimentDetails.NotifyID, "NotifyID should not be empty")
			},
		},
		{
			name:         "experiment run with empty ID",
			projectID:    "test-project-id",
			experimentID: "",
			wantErr:      true,
			validateFn:   nil,
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

			result, err := RunExperiment(tt.projectID, tt.experimentID, client.credentials)

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

func TestGetExperimentList(t *testing.T) {
	tests := []struct {
		name       string
		projectID  string
		request    model.ListExperimentRequest
		setup      func(*LitmusClient) // optional setup steps
		wantErr    bool
		validateFn func(*testing.T, *ExperimentListData)
	}{
		{
			name:      "successful experiment list fetch",
			projectID: "test-project-id",
			request: model.ListExperimentRequest{
				Pagination: &model.Pagination{
					Page:  1,
					Limit: 10,
				},
			},
			wantErr: false,
			validateFn: func(t *testing.T, result *ExperimentListData) {
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.NotNil(t, result.Data.ListExperimentDetails, "ListExperimentDetails should not be nil")

				// Check total count is a non-negative number
				assert.GreaterOrEqual(t, result.Data.ListExperimentDetails.TotalNoOfExperiments, 0,
					"Total number of experiments should be non-negative")

				// If there are experiments, validate their structure
				if len(result.Data.ListExperimentDetails.Experiments) > 0 {
					for _, exp := range result.Data.ListExperimentDetails.Experiments {
						assert.NotEmpty(t, exp.ExperimentID, "Experiment ID should not be empty")
						assert.NotEmpty(t, exp.Name, "Experiment name should not be empty")
					}
				}
			},
		},
		{
			name:      "experiment list with pagination",
			projectID: "test-project-id",
			request: model.ListExperimentRequest{
				Pagination: &model.Pagination{
					Page:  1,
					Limit: 5,
				},
			},
			wantErr: false,
			validateFn: func(t *testing.T, result *ExperimentListData) {
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.NotNil(t, result.Data.ListExperimentDetails, "ListExperimentDetails should not be nil")

				// Verify pagination works by checking max results
				if len(result.Data.ListExperimentDetails.Experiments) > 0 {
					assert.LessOrEqual(t,
						len(result.Data.ListExperimentDetails.Experiments),
						5,
						"Should return 5 or fewer results with limit=5")
				}
			},
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

			result, err := GetExperimentList(tt.projectID, tt.request, client.credentials)

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

func TestGetExperimentRunsList(t *testing.T) {
	tests := []struct {
		name       string
		projectID  string
		request    model.ListExperimentRunRequest
		setup      func(*LitmusClient) // optional setup steps
		wantErr    bool
		validateFn func(*testing.T, *ExperimentRunListData) // fixed type
	}{
		{
			name:      "successful experiment runs list fetch",
			projectID: "test-project-id",
			request: model.ListExperimentRunRequest{
				Pagination: &model.Pagination{
					Page:  1,
					Limit: 10,
				},
			},
			wantErr: false,
			validateFn: func(t *testing.T, result *ExperimentRunListData) { // fixed type
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.NotNil(t, result.Data.ListExperimentRunDetails, "ListExperimentRunDetails should not be nil")
				assert.GreaterOrEqual(t, result.Data.ListExperimentRunDetails.TotalNoOfExperimentRuns, 0,
					"Total number of experiment runs should be non-negative")
			},
		},
		{
			name:      "experiment runs list with pagination",
			projectID: "test-project-id",
			request: model.ListExperimentRunRequest{
				Pagination: &model.Pagination{
					Page:  1,
					Limit: 5,
				},
				// Removed filter to avoid the linter error
			},
			wantErr: false,
			validateFn: func(t *testing.T, result *ExperimentRunListData) { // fixed type
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.NotNil(t, result.Data.ListExperimentRunDetails, "ListExperimentRunDetails should not be nil")

				// Verify pagination works by checking max results
				if len(result.Data.ListExperimentRunDetails.ExperimentRuns) > 0 {
					assert.LessOrEqual(t,
						len(result.Data.ListExperimentRunDetails.ExperimentRuns),
						5,
						"Should return 5 or fewer results with limit=5")
				}
			},
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

			result, err := GetExperimentRunsList(tt.projectID, tt.request, client.credentials)

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

func TestDeleteChaosExperiment(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		experimentID string
		wantErr      bool
		validateFn   func(*testing.T, *DeleteChaosExperimentData)
	}{
		{
			name:         "successful experiment deletion",
			projectID:    "test-project-id",
			experimentID: "test-experiment-id",
			wantErr:      false,
			validateFn: func(t *testing.T, result *DeleteChaosExperimentData) {
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.True(t, result.Data.IsDeleted, "IsDeleted should be true")
			},
		},
		{
			name:         "experiment deletion with empty ID",
			projectID:    "test-project-id",
			experimentID: "",
			wantErr:      true,
			validateFn:   nil, // We expect an error, so no validation needed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient()
			assert.NoError(t, err, "Failed to create Litmus client")

			experimentID := tt.experimentID
			result, err := DeleteChaosExperiment(tt.projectID, &experimentID, client.credentials)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.validateFn != nil {
				tt.validateFn(t, &result)
			}
		})
	}
}

func TestCreateExperiment(t *testing.T) {
	tests := []struct {
		name       string
		projectID  string
		request    model.SaveChaosExperimentRequest
		setup      func(*LitmusClient) // optional setup steps
		wantErr    bool
		validateFn func(*testing.T, *RunExperimentResponse)
	}{
		{
			name:      "successful experiment creation and run",
			projectID: "test-project-id",
			request: model.SaveChaosExperimentRequest{
				ID:   "test-experiment-id",
				Name: "test-experiment",
			},
			wantErr: false,
			validateFn: func(t *testing.T, result *RunExperimentResponse) {
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.NotEmpty(t, result.Data.RunExperimentDetails.NotifyID, "NotifyID should not be empty")
			},
		},
		{
			name:      "experiment creation with empty ID",
			projectID: "test-project-id",
			request: model.SaveChaosExperimentRequest{
				ID:   "",
				Name: "test-experiment",
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

			result, err := CreateExperiment(tt.projectID, tt.request, client.credentials)

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

func TestResponseStructureMarshaling(t *testing.T) {
	t.Run("RunExperimentResponse", func(t *testing.T) {
		jsonData := `{
			"data": {
				"runChaosExperiment": {
					"notifyID": "test-notify-id"
				}
			}
		}`

		var response RunExperimentResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		assert.NoError(t, err, "Unmarshaling should succeed")
		assert.Equal(t, "test-notify-id", response.Data.RunExperimentDetails.NotifyID)
	})

	t.Run("SaveExperimentData", func(t *testing.T) {
		jsonData := `{
			"data": {
				"saveChaosExperiment": "Experiment saved successfully"
			}
		}`

		var response SaveExperimentData
		err := json.Unmarshal([]byte(jsonData), &response)
		assert.NoError(t, err, "Unmarshaling should succeed")
		assert.Equal(t, "Experiment saved successfully", response.Data.Message)
	})

	t.Run("DeleteChaosExperimentData", func(t *testing.T) {
		jsonData := `{
			"data": {
				"deleteChaosExperiment": true
			}
		}`

		var response DeleteChaosExperimentData
		err := json.Unmarshal([]byte(jsonData), &response)
		assert.NoError(t, err, "Unmarshaling should succeed")
		assert.True(t, response.Data.IsDeleted)
	})
}
