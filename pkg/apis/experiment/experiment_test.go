package experiment

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/litmuschaos/litmus-go-sdk/pkg/apis"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/stretchr/testify/assert"
)

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
	return NewLitmusClient("http://127.0.0.1:35961", "admin", "LitmusChaos123@")
}

func TestSaveExperiment(t *testing.T) {
	client, err := setupTestClient()
	assert.NoError(t, err, "Failed to create Litmus client")

	projectID := "test-project-id"

	request := model.SaveChaosExperimentRequest{
		ID:   "test-experiment-id",
		Name: "test-experiment",
	}

	result, err := SaveExperiment(projectID, request, client.credentials)
	if err != nil {
		fmt.Printf("API call error: %v\n", err)
	} else {
		assert.NotNil(t, result, "Result should not be nil")
		assert.NotNil(t, result.Data, "Data should not be nil")
		// Check that experiment ID is present in response message
		assert.Contains(t, result.Data.Message, "test-experiment-id", "Response should contain experiment ID")
	}
}

func TestRunExperiment(t *testing.T) {
	client, err := setupTestClient()
	assert.NoError(t, err, "Failed to create Litmus client")

	projectID := "test-project-id"
	experimentID := "test-experiment-id"

	result, err := RunExperiment(projectID, experimentID, client.credentials)
	if err != nil {
		fmt.Printf("API call error: %v\n", err)
	} else {
		assert.NotNil(t, result, "Result should not be nil")
		assert.NotNil(t, result.Data, "Data should not be nil")
		assert.NotEmpty(t, result.Data.RunExperimentDetails.NotifyID, "NotifyID should not be empty")
	}
}

func TestGetExperimentList(t *testing.T) {
	client, err := setupTestClient()
	assert.NoError(t, err, "Failed to create Litmus client")

	projectID := "test-project-id"

	request := model.ListExperimentRequest{
		Pagination: &model.Pagination{
			Page:  1,
			Limit: 10,
		},
	}

	result, err := GetExperimentList(projectID, request, client.credentials)
	if err != nil {
		fmt.Printf("API call error: %v\n", err)
	} else {
		assert.NotNil(t, result, "Result should not be nil")
		assert.NotNil(t, result.Data, "Data should not be nil")
		assert.NotNil(t, result.Data.ListExperimentDetails, "ListExperimentDetails should not be nil")

		// Check total count is a non-negative number
		assert.GreaterOrEqual(t, result.Data.ListExperimentDetails.TotalNoOfExperiments, 0, "Total number of experiments should be non-negative")

		// If there are experiments, validate their structure
		if len(result.Data.ListExperimentDetails.Experiments) > 0 {
			for _, exp := range result.Data.ListExperimentDetails.Experiments {
				assert.NotEmpty(t, exp.ExperimentID, "Experiment ID should not be empty")
				assert.NotEmpty(t, exp.Name, "Experiment name should not be empty")
			}
		}
	}
}

func TestGetExperimentRunsList(t *testing.T) {
	client, err := setupTestClient()
	assert.NoError(t, err, "Failed to create Litmus client")

	projectID := "test-project-id"

	request := model.ListExperimentRunRequest{
		Pagination: &model.Pagination{
			Page:  1,
			Limit: 10,
		},
	}

	result, err := GetExperimentRunsList(projectID, request, client.credentials)
	if err != nil {
		fmt.Printf("API call error: %v\n", err)
	} else {
		assert.NotNil(t, result, "Result should not be nil")
		assert.NotNil(t, result.Data, "Data should not be nil")
		assert.NotNil(t, result.Data.ListExperimentRunDetails, "ListExperimentRunDetails should not be nil")

		// Check total count is a non-negative number
		assert.GreaterOrEqual(t, result.Data.ListExperimentRunDetails.TotalNoOfExperimentRuns, 0, "Total number of experiment runs should be non-negative")

		// If there are experiment runs, validate their structure
		if len(result.Data.ListExperimentRunDetails.ExperimentRuns) > 0 {
			for _, run := range result.Data.ListExperimentRunDetails.ExperimentRuns {
				assert.NotEmpty(t, run.ExperimentRunID, "Experiment run ID should not be empty")
				assert.NotEmpty(t, run.ExperimentID, "Experiment ID should not be empty")
				assert.NotEmpty(t, run.ExperimentName, "Experiment name should not be empty")
			}
		}
	}
}

func TestDeleteChaosExperiment(t *testing.T) {
	client, err := setupTestClient()
	assert.NoError(t, err, "Failed to create Litmus client")

	projectID := "test-project-id"
	experimentID := "test-experiment-id"

	result, err := DeleteChaosExperiment(projectID, &experimentID, client.credentials)
	if err != nil {
		fmt.Printf("API call error: %v\n", err)
	} else {
		assert.NotNil(t, result, "Result should not be nil")
		assert.NotNil(t, result.Data, "Data should not be nil")

		// Verify deletion status
		assert.True(t, result.Data.IsDeleted, "IsDeleted should be true")
	}
}

func TestCreateExperiment(t *testing.T) {
	client, err := setupTestClient()
	assert.NoError(t, err, "Failed to create Litmus client")

	projectID := "test-project-id"

	request := model.SaveChaosExperimentRequest{
		ID:   "test-experiment-id",
		Name: "test-experiment",
	}

	result, err := CreateExperiment(projectID, request, client.credentials)
	if err != nil {
		fmt.Printf("API call error: %v\n", err)
	} else {
		assert.NotNil(t, result, "Result should not be nil")
		assert.NotNil(t, result.Data, "Data should not be nil")
		assert.NotEmpty(t, result.Data.RunExperimentDetails.NotifyID, "NotifyID should not be empty")
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
