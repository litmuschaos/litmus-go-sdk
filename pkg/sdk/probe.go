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

// ProbeClient defines the interface for probe operations
type ProbeClient interface {
	// List retrieves all probes
	List() (interface{}, error)

	// Create creates a new probe
	Create(name string, config map[string]interface{}) (interface{}, error)

	// Delete removes a probe
	Delete(id string) error

	// Update updates a probe
	Update(id string, config map[string]interface{}) (interface{}, error)

	// Get retrieves probe details
	Get(id string) (interface{}, error)

	// Execute runs a probe
	Execute(id string, params map[string]string) (interface{}, error)
}

// probeClient implements the ProbeClient interface
type probeClient struct {
	credentials types.Credentials
}

// List retrieves all probes
func (c *probeClient) List() (interface{}, error) {
	// TODO: Implement when probe API is available
	return nil, nil
}

// Create creates a new probe
func (c *probeClient) Create(name string, config map[string]interface{}) (interface{}, error) {
	// TODO: Implement when probe API is available
	return nil, nil
}

// Delete removes a probe
func (c *probeClient) Delete(id string) error {
	// TODO: Implement when probe API is available
	return nil
}

// Update updates a probe
func (c *probeClient) Update(id string, config map[string]interface{}) (interface{}, error) {
	// TODO: Implement when probe API is available
	return nil, nil
}

// Get retrieves probe details
func (c *probeClient) Get(id string) (interface{}, error) {
	// TODO: Implement when probe API is available
	return nil, nil
}

// Execute runs a probe
func (c *probeClient) Execute(id string, params map[string]string) (interface{}, error) {
	// TODO: Implement when probe API is available
	return nil, nil
}
