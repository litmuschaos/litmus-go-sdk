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

// EnvironmentClient defines the interface for environment operations
type EnvironmentClient interface {
	// List retrieves all environments
	List() (interface{}, error)

	// Create creates a new environment
	Create(name string, config map[string]interface{}) (interface{}, error)

	// Delete removes an environment
	Delete(id string) error

	// Update updates an environment
	Update(id string, config map[string]interface{}) (interface{}, error)

	// Get retrieves environment details
	Get(id string) (interface{}, error)
}

// environmentClient implements the EnvironmentClient interface
type environmentClient struct {
	credentials types.Credentials
}

// List retrieves all environments
func (c *environmentClient) List() (interface{}, error) {
	// TODO: Implement when environment API is available
	return nil, nil
}

// Create creates a new environment
func (c *environmentClient) Create(name string, config map[string]interface{}) (interface{}, error) {
	// TODO: Implement when environment API is available
	return nil, nil
}

// Delete removes an environment
func (c *environmentClient) Delete(id string) error {
	// TODO: Implement when environment API is available
	return nil
}

// Update updates an environment
func (c *environmentClient) Update(id string, config map[string]interface{}) (interface{}, error) {
	// TODO: Implement when environment API is available
	return nil, nil
}

// Get retrieves environment details
func (c *environmentClient) Get(id string) (interface{}, error) {
	// TODO: Implement when environment API is available
	return nil, nil
}
