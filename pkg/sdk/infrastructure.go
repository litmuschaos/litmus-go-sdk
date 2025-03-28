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
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
)

// InfrastructureClient defines the interface for infrastructure operations
type InfrastructureClient interface {
	// List retrieves all infrastructure resources
	List() (interface{}, error)

	// Create creates a new infrastructure resource
	Create(name string, config map[string]interface{}) (interface{}, error)

	// Delete removes an infrastructure resource
	Delete(id string) error

	// Update updates an infrastructure resource
	Update(id string, config map[string]interface{}) (interface{}, error)

	// Get retrieves infrastructure details
	Get(id string) (interface{}, error)

	// Connect establishes a connection to an infrastructure
	Connect(id string, params map[string]string) error

	// Disconnect terminates a connection to an infrastructure
	Disconnect(id string) error
}

// infrastructureClient implements the InfrastructureClient interface
type infrastructureClient struct {
	credentials types.Credentials
}

// List retrieves all infrastructure resources
func (c *infrastructureClient) List() (interface{}, error) {
	// TODO: Implement when infrastructure API is available
	return nil, nil
}

// Create creates a new infrastructure resource
func (c *infrastructureClient) Create(name string, config map[string]interface{}) (interface{}, error) {
	// TODO: Implement when infrastructure API is available
	return nil, nil
}

// Delete removes an infrastructure resource
func (c *infrastructureClient) Delete(id string) error {
	// TODO: Implement when infrastructure API is available
	return nil
}

// Update updates an infrastructure resource
func (c *infrastructureClient) Update(id string, config map[string]interface{}) (interface{}, error) {
	// TODO: Implement when infrastructure API is available
	return nil, nil
}

// Get retrieves infrastructure details
func (c *infrastructureClient) Get(id string) (interface{}, error) {
	// TODO: Implement when infrastructure API is available
	return nil, nil
}

// Connect establishes a connection to an infrastructure
func (c *infrastructureClient) Connect(id string, params map[string]string) error {
	// TODO: Implement when infrastructure API is available
	return nil
}

// Disconnect terminates a connection to an infrastructure
func (c *infrastructureClient) Disconnect(id string) error {
	// TODO: Implement when infrastructure API is available
	return nil
}
