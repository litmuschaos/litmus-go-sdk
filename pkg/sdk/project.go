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
	"github.com/litmuschaos/litmus-go-sdk/pkg/apis"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
)

// ProjectClient defines the interface for project operations
type ProjectClient interface {
	// List retrieves all projects
	List() (apis.ListProjectResponse, error)

	// Create creates a new project with the given name
	Create(name string) (apis.CreateProjectResponse, error)

	// GetDetails retrieves detailed information about projects
	GetDetails() (apis.ProjectDetails, error)
}

// projectClient implements the ProjectClient interface
type projectClient struct {
	credentials types.Credentials
}

// List retrieves all projects
func (c *projectClient) List() (apis.ListProjectResponse, error) {
	return apis.ListProject(c.credentials)
}

// Create creates a new project with the given name
func (c *projectClient) Create(name string) (apis.CreateProjectResponse, error) {
	return apis.CreateProjectRequest(name, c.credentials)
}

// GetDetails retrieves detailed information about projects
func (c *projectClient) GetDetails() (apis.ProjectDetails, error) {
	return apis.GetProjectDetails(c.credentials)
}
