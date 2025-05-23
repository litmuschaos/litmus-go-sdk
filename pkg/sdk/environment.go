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

	"github.com/litmuschaos/litmus-go-sdk/pkg/apis/environment"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
)

// EnvironmentClient defines the interface for environment operations
type EnvironmentClient interface {
	// List retrieves all environments
	List() (models.ListEnvironmentResponse, error)

	// Create creates a new environment
	Create(name string, request models.CreateEnvironmentRequest) (models.Environment, error)

	// Delete removes an environment
	Delete(id string) error

	// Get retrieves environment details
	Get(id string) (models.Environment, error)
}

// environmentClient implements the EnvironmentClient interface
type environmentClient struct {
	credentials types.Credentials
}

// List retrieves all environments
func (c *environmentClient) List() (models.ListEnvironmentResponse, error) {
	if c.credentials.Endpoint == "" {
		return models.ListEnvironmentResponse{}, fmt.Errorf("endpoint not set in credentials")
	}

	response, err := environment.ListChaosEnvironments(c.credentials.ProjectID, c.credentials)
	if err != nil {
		return models.ListEnvironmentResponse{}, fmt.Errorf("failed to list environments: %w", err)
	}

	return response.ListEnvironments, nil
}

// Create creates a new environment
func (c *environmentClient) Create(name string, request models.CreateEnvironmentRequest) (models.Environment, error) {
	if c.credentials.Endpoint == "" {
		return models.Environment{}, fmt.Errorf("endpoint not set in credentials")
	}

	if c.credentials.ProjectID == "" {
		return models.Environment{}, fmt.Errorf("project ID not set in credentials")
	}

	// Set name if not already set
	if request.Name == "" {
		request.Name = name
	}
	
	// Set default tags if not provided
	if len(request.Tags) == 0 {
		request.Tags = []string{"litmus-sdk"}
	}

	response, err := environment.CreateEnvironment(c.credentials.ProjectID, request, c.credentials)
	if err != nil {
		return models.Environment{}, fmt.Errorf("failed to create environment: %w", err)
	}

	return response.CreateEnvironment, nil
}

// Delete removes an environment
func (c *environmentClient) Delete(id string) error {
	if c.credentials.Endpoint == "" {
		return fmt.Errorf("endpoint not set in credentials")
	}

	if c.credentials.ProjectID == "" {
		return fmt.Errorf("project ID not set in credentials")
	}

	if id == "" {
		return fmt.Errorf("environment ID cannot be empty")
	}

	_, err := environment.DeleteEnvironment(c.credentials.ProjectID, id, c.credentials)
	if err != nil {
		return fmt.Errorf("failed to delete environment: %w", err)
	}

	return nil
}

// Get retrieves environment details
func (c *environmentClient) Get(id string) (models.Environment, error) {
	if c.credentials.Endpoint == "" {
		return models.Environment{}, fmt.Errorf("endpoint not set in credentials")
	}

	if c.credentials.ProjectID == "" {
		return models.Environment{}, fmt.Errorf("project ID not set in credentials")
	}

	if id == "" {
		return models.Environment{}, fmt.Errorf("environment ID cannot be empty")
	}

	response, err := environment.GetChaosEnvironment(c.credentials.ProjectID, id, c.credentials)
	if err != nil {
		return models.Environment{}, fmt.Errorf("failed to get environment: %w", err)
	}

	return response.GetEnvironment, nil
}
