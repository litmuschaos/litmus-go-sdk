/*
Copyright © 2025 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package sdk

import (
	"fmt"

	"github.com/litmuschaos/litmus-go-sdk/pkg/apis/experiment"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
)

// ExperimentClient defines the interface for experiment operations
type ExperimentClient interface {
	// List retrieves all experiments
	List() (models.ListExperimentResponse, error)

	// Create creates a new experiment
	Create(name string, config map[string]interface{}) (string, error)

	// Delete removes an experiment
	Delete(id string) error

	// Update updates an experiment
	Update(id string, config map[string]interface{}) (string, error)

	// Get retrieves experiment details
	Get(id string) (models.ExperimentRun, error)

	// Run starts an experiment
	Run(id string) (string, error)

	// GetRunPhase retrieves just the status/phase of a specific experiment run
	GetRunPhase(runID string) (string, error)
}

// experimentClient implements the ExperimentClient interface
type experimentClient struct {
	credentials types.Credentials
}

// List retrieves all experiments
func (c *experimentClient) List() (models.ListExperimentResponse, error) {
	if c.credentials.ServerEndpoint == "" {
		return models.ListExperimentResponse{}, fmt.Errorf("server endpoint not set in credentials")
	}
	
	if c.credentials.ProjectID == "" {
		return models.ListExperimentResponse{}, fmt.Errorf("project ID not set in credentials")
	}

	request := models.ListExperimentRequest{}
	
	response, err := experiment.GetExperimentList(c.credentials.ProjectID, request, c.credentials)
	if err != nil {
		return models.ListExperimentResponse{}, fmt.Errorf("failed to list experiments: %w", err)
	}

	return response.ListExperimentDetails, nil
}

// Create creates a new experiment
func (c *experimentClient) Create(name string, config map[string]interface{}) (string, error) {
	if c.credentials.ServerEndpoint == "" {
		return "", fmt.Errorf("server endpoint not set in credentials")
	}
	
	if c.credentials.ProjectID == "" {
		return "", fmt.Errorf("project ID not set in credentials")
	}

	// Build the request with basic details
	request := models.SaveChaosExperimentRequest{
		Name:        name,
		Description: fmt.Sprintf("Experiment created via Litmus SDK: %s", name),
	}

	// Save the experiment first
	saveResp, err := experiment.SaveExperiment(c.credentials.ProjectID, request, c.credentials)
	if err != nil {
		return "", fmt.Errorf("failed to create experiment: %w", err)
	}

	return saveResp.Message, nil
}

// Delete removes an experiment
func (c *experimentClient) Delete(id string) error {
	if c.credentials.ServerEndpoint == "" {
		return fmt.Errorf("server endpoint not set in credentials")
	}
	
	if c.credentials.ProjectID == "" {
		return fmt.Errorf("project ID not set in credentials")
	}

	if id == "" {
		return fmt.Errorf("experiment ID cannot be empty")
	}

	response, err := experiment.DeleteChaosExperiment(c.credentials.ProjectID, &id, c.credentials)
	if err != nil {
		return fmt.Errorf("failed to delete experiment: %w", err)
	}

	if !response.IsDeleted {
		return fmt.Errorf("experiment deletion was not successful")
	}

	return nil
}

// Update updates an experiment
func (c *experimentClient) Update(id string, config map[string]interface{}) (string, error) {
	if c.credentials.ServerEndpoint == "" {
		return "", fmt.Errorf("server endpoint not set in credentials")
	}
	
	if c.credentials.ProjectID == "" {
		return "", fmt.Errorf("project ID not set in credentials")
	}

	if id == "" {
		return "", fmt.Errorf("experiment ID cannot be empty")
	}

	// Build the request with updated details
	request := models.SaveChaosExperimentRequest{
		ID: id,
	}
	

	saveResp, err := experiment.SaveExperiment(c.credentials.ProjectID, request, c.credentials)
	if err != nil {
		return "", fmt.Errorf("failed to update experiment: %w", err)
	}

	return saveResp.Message, nil
}

func (c *experimentClient) Get(runID string) (models.ExperimentRun, error) {
	if c.credentials.ServerEndpoint == "" {
		return models.ExperimentRun{}, fmt.Errorf("server endpoint not set in credentials")
	}
	
	if c.credentials.ProjectID == "" {
		return models.ExperimentRun{}, fmt.Errorf("project ID not set in credentials")
	}

	if runID == "" {
		return models.ExperimentRun{}, fmt.Errorf("experiment run ID cannot be empty")
	}

	response, err := experiment.GetExperimentRun(c.credentials.ProjectID, runID, c.credentials)
	if err != nil {
		return models.ExperimentRun{}, fmt.Errorf("failed to get experiment run: %w", err)
	}

	// Return the full experiment run data
	return response.ExperimentRun, nil
}

// Run starts an experiment
func (c *experimentClient) Run(id string) (string, error) {
	if c.credentials.ServerEndpoint == "" {
		return "", fmt.Errorf("server endpoint not set in credentials")
	}
	
	if c.credentials.ProjectID == "" {
		return "", fmt.Errorf("project ID not set in credentials")
	}

	if id == "" {
		return "", fmt.Errorf("experiment ID cannot be empty")
	}

	response, err := experiment.RunExperiment(c.credentials.ProjectID, id, c.credentials)
	if err != nil {
		return "", fmt.Errorf("failed to run experiment: %w", err)
	}

	return response.RunChaosExperiment.NotifyID, nil
}

// GetRunPhase retrieves just the status/phase of a specific experiment run
func (c *experimentClient) GetRunPhase(runID string) (string, error) {
	experimentRun, err := c.Get(runID)
	if err != nil {
		return "", fmt.Errorf("failed to get experiment run phase: %w", err)
	}
	
	return string(experimentRun.Phase), nil
}

