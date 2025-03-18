package experiment

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/stretchr/testify/assert"
)

func setupTestCredentials() types.Credentials {
	return types.Credentials{
		ServerEndpoint: "http://127.0.0.1:35961",
		Token:          "test-token",
	}
}

func TestSaveExperiment(t *testing.T) {
	cred := setupTestCredentials()
	projectID := "test-project-id"

	request := model.SaveChaosExperimentRequest{
		ID:   "test-experiment-id",
		Name: "test-experiment",
	}

	result, err := SaveExperiment(projectID, request, cred)
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
	cred := setupTestCredentials()
	projectID := "test-project-id"
	experimentID := "test-experiment-id"

	result, err := RunExperiment(projectID, experimentID, cred)
	if err != nil {
		fmt.Printf("API call error: %v\n", err)
	} else {
		assert.NotNil(t, result, "Result should not be nil")
		assert.NotNil(t, result.Data, "Data should not be nil")
		assert.NotEmpty(t, result.Data.RunExperimentDetails.NotifyID, "NotifyID should not be empty")
	}
}

func TestGetExperimentList(t *testing.T) {
	cred := setupTestCredentials()
	projectID := "test-project-id"

	request := model.ListExperimentRequest{
		Pagination: &model.Pagination{
			Page:  1,
			Limit: 10,
		},
	}

	result, err := GetExperimentList(projectID, request, cred)
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
	cred := setupTestCredentials()
	projectID := "test-project-id"

	request := model.ListExperimentRunRequest{
		Pagination: &model.Pagination{
			Page:  1,
			Limit: 10,
		},
	}

	result, err := GetExperimentRunsList(projectID, request, cred)
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
	cred := setupTestCredentials()
	projectID := "test-project-id"
	experimentID := "test-experiment-id"

	result, err := DeleteChaosExperiment(projectID, &experimentID, cred)
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
	cred := setupTestCredentials()
	projectID := "test-project-id"

	request := model.SaveChaosExperimentRequest{
		ID:   "test-experiment-id",
		Name: "test-experiment",
	}

	result, err := CreateExperiment(projectID, request, cred)
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
