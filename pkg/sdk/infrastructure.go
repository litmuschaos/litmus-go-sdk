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

	"github.com/litmuschaos/litmus-go-sdk/pkg/apis/infrastructure"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
)

// InfrastructureClient defines the interface for infrastructure operations
type InfrastructureClient interface {
	// List retrieves all infrastructure resources
	List() (interface{}, error)

	// Create creates a new infrastructure resource
	Create(name string, config map[string]interface{}) (interface{}, error)

	// Delete removes an infrastructure resource
	Delete(id string) error

	// Get retrieves infrastructure details
	Get(id string) (interface{}, error)

	// Disconnect terminates a connection to an infrastructure
	Disconnect(id string) error
}

// infrastructureClient implements the InfrastructureClient interface
type infrastructureClient struct {
	credentials types.Credentials
}

// List retrieves all infrastructure resources
func (c *infrastructureClient) List() (interface{}, error) {
	if c.credentials.ServerEndpoint == "" {
		return nil, fmt.Errorf("server endpoint not set in credentials")
	}
	
	if c.credentials.ProjectID == "" {
		return nil, fmt.Errorf("project ID not set in credentials")
	}

	request := models.ListInfraRequest{}
	
	response, err := infrastructure.GetInfraList(c.credentials, c.credentials.ProjectID, request)
	if err != nil {
		return nil, fmt.Errorf("failed to list infrastructure resources: %w", err)
	}

	return response.Data.ListInfraDetails, nil
}

// Create creates a new infrastructure resource
func (c *infrastructureClient) Create(name string, config map[string]interface{}) (interface{}, error) {
	if c.credentials.ServerEndpoint == "" {
		return nil, fmt.Errorf("server endpoint not set in credentials")
	}
	
	if c.credentials.ProjectID == "" {
		return nil, fmt.Errorf("project ID not set in credentials")
	}

	// Extract values from config or use defaults
	var (
		description    = getStringFromConfig(config, "description", fmt.Sprintf("Infrastructure created via Litmus SDK: %s", name))
		platformName   = getStringFromConfig(config, "platformName", "default-platform")
		environmentID  = getStringFromConfig(config, "environmentID", "")
		namespace      = getStringFromConfig(config, "namespace", "litmus")
		serviceAccount = getStringFromConfig(config, "serviceAccount", "litmus")
		nsExists       = getBoolFromConfig(config, "nsExists", false)
		saExists       = getBoolFromConfig(config, "saExists", false)
		skipSSL        = getBoolFromConfig(config, "skipSSL", false)
		nodeSelector   = getStringFromConfig(config, "nodeSelector", "")
		tolerations    = getStringFromConfig(config, "tolerations", "")
		mode           = getStringFromConfig(config, "mode", string(models.InfraScopeNamespace))
	)

	// Create infrastructure request
	infra := types.Infra{
		ProjectID:      c.credentials.ProjectID,
		InfraName:      name,
		Description:    description,
		PlatformName:   platformName,
		EnvironmentID:  environmentID,
		Namespace:      namespace,
		ServiceAccount: serviceAccount,
		NsExists:       nsExists,
		SAExists:       saExists,
		SkipSSL:        skipSSL,
		NodeSelector:   nodeSelector,
		Tolerations:    tolerations,
		Mode:           models.InfraScope(mode).String(),
	}

	response, err := infrastructure.ConnectInfra(infra, c.credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to create infrastructure: %w", err)
	}

	return response.Data.RegisterInfraDetails, nil
}

// Delete removes an infrastructure resource
func (c *infrastructureClient) Delete(id string)  error {
	return c.Disconnect(id)
}

// Get retrieves infrastructure details
func (c *infrastructureClient) Get(id string) (interface{}, error) {
	// Currently, there's no specific API for getting a single infrastructure
	// We'll get the list and filter for the requested ID
	if c.credentials.ServerEndpoint == "" {
		return nil, fmt.Errorf("server endpoint not set in credentials")
	}
	
	if c.credentials.ProjectID == "" {
		return nil, fmt.Errorf("project ID not set in credentials")
	}

	if id == "" {
		return nil, fmt.Errorf("infrastructure ID cannot be empty")
	}

	request := models.ListInfraRequest{
		InfraIDs: []string{id},
	}
	
	response, err := infrastructure.GetInfraList(c.credentials, c.credentials.ProjectID, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get infrastructure: %w", err)
	}

	if len(response.Data.ListInfraDetails.Infras) == 0 {
		return nil, fmt.Errorf("infrastructure not found with ID: %s", id)
	}

	return response.Data.ListInfraDetails.Infras[0], nil
}

// Disconnect terminates a connection to an infrastructure
func (c *infrastructureClient) Disconnect(id string)  error {
	if c.credentials.ServerEndpoint == "" {
		return fmt.Errorf("server endpoint not set in credentials")
	}
	
	if c.credentials.ProjectID == "" {
		return fmt.Errorf("project ID not set in credentials")
	}

	if id == "" {
		return fmt.Errorf("infrastructure ID cannot be empty")
	}

	_, err := infrastructure.DisconnectInfra(c.credentials.ProjectID, id, c.credentials)
	if err != nil {
		return fmt.Errorf("failed to disconnect infrastructure: %w", err)
	}

	return nil
}

// Helper functions for config extraction
func getStringFromConfig(config map[string]interface{}, key, defaultValue string) string {
	if val, ok := config[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

func getBoolFromConfig(config map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := config[key]; ok {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return defaultValue
}
