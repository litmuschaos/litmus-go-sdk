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

// ExperimentClient defines the interface for experiment operations
type ExperimentClient interface {
	// List retrieves all experiments
	List() (interface{}, error)

	// Create creates a new experiment
	Create(name string, config map[string]interface{}) (interface{}, error)

	// Delete removes an experiment
	Delete(id string) error

	// Update updates an experiment
	Update(id string, config map[string]interface{}) (interface{}, error)

	// Get retrieves experiment details
	Get(id string) (interface{}, error)

	// Run starts an experiment
	Run(id string) (interface{}, error)

	// Stop halts a running experiment
	Stop(id string) error
}

// experimentClient implements the ExperimentClient interface
type experimentClient struct {
	credentials types.Credentials
}

// List retrieves all experiments
func (c *experimentClient) List() (interface{}, error) {
	// TODO: Implement when experiment API is available
	return nil, nil
}

// Create creates a new experiment
func (c *experimentClient) Create(name string, config map[string]interface{}) (interface{}, error) {
	// TODO: Implement when experiment API is available
	return nil, nil
}

// Delete removes an experiment
func (c *experimentClient) Delete(id string) error {
	// TODO: Implement when experiment API is available
	return nil
}

// Update updates an experiment
func (c *experimentClient) Update(id string, config map[string]interface{}) (interface{}, error) {
	// TODO: Implement when experiment API is available
	return nil, nil
}

// Get retrieves experiment details
func (c *experimentClient) Get(id string) (interface{}, error) {
	// TODO: Implement when experiment API is available
	return nil, nil
}

// Run starts an experiment
func (c *experimentClient) Run(id string) (interface{}, error) {
	// TODO: Implement when experiment API is available
	return nil, nil
}

// Stop halts a running experiment
func (c *experimentClient) Stop(id string) error {
	// TODO: Implement when experiment API is available
	return nil
}
