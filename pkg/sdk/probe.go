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

	"github.com/litmuschaos/litmus-go-sdk/pkg/apis/probe"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
)

// ProbeClient defines the interface for probe operations
type ProbeClient interface {
	// List retrieves all probes
	List(projectID string) (interface{}, error)

	// Delete removes a probe
	Delete(projectID string, id string) error

	// Get retrieves probe details
	Get(projectID string, id string) (interface{}, error)

	// GetProbeYAML retrieves the YAML configuration for a probe
	GetProbeYAML(projectID string, id string, params map[string]string) (interface{}, error)
}

// probeClient implements the ProbeClient interface
type probeClient struct {
	credentials types.Credentials
}

// List retrieves all probes
func (c *probeClient) List(projectID string) (interface{}, error) {
	if c.credentials.ServerEndpoint == "" {
		return nil, fmt.Errorf("server endpoint not set in credentials")
	}

	if projectID == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}

	// Use all probe types as no specific type is requested
	var probeTypes []*models.ProbeType
	
	response, err := probe.ListProbeRequest(projectID, probeTypes, c.credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to list probes: %w", err)
	}

	return response.Data.Probes, nil
}


// Delete removes a probe
func (c *probeClient) Delete(projectID string, id string) error {
	if c.credentials.ServerEndpoint == "" {
		return fmt.Errorf("server endpoint not set in credentials")
	}

	if projectID == "" {
		return fmt.Errorf("project ID cannot be empty")
	}

	if id == "" {
		return fmt.Errorf("probe ID cannot be empty")
	}

	response, err := probe.DeleteProbeRequest(projectID, id, c.credentials)
	if err != nil {
		return fmt.Errorf("failed to delete probe: %w", err)
	}

	if !response.Data.DeleteProbe {
		return fmt.Errorf("probe deletion was not successful")
	}

	return nil
}


// Get retrieves probe details
func (c *probeClient) Get(projectID string, id string) (interface{}, error) {
	if c.credentials.ServerEndpoint == "" {
		return nil, fmt.Errorf("server endpoint not set in credentials")
	}

	if projectID == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}

	if id == "" {
		return nil, fmt.Errorf("probe ID cannot be empty")
	}

	response, err := probe.GetProbeRequest(projectID, id, c.credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to get probe: %w", err)
	}

	return response.Data.GetProbe, nil
}

// Execute runs a probe
func (c *probeClient) GetProbeYAML(projectID string, id string, params map[string]string) (interface{}, error) {
	if c.credentials.ServerEndpoint == "" {
		return nil, fmt.Errorf("server endpoint not set in credentials")
	}

	if projectID == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}

	if id == "" {
		return nil, fmt.Errorf("probe ID cannot be empty")
	}

	// Create a request to get probe YAML
	request := models.GetProbeYAMLRequest{
		ProbeName: id,
		Mode: models.Mode(params["mode"]),
	}

	response, err := probe.GetProbeYAMLRequest(projectID, request, c.credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to execute probe: %w", err)
	}

	return response.Data.GetProbeYAML, nil
}
