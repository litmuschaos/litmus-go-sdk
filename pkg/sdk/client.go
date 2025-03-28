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

	"github.com/litmuschaos/litmus-go-sdk/pkg/apis"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
)

// Client is the interface for the Litmus API client
type Client interface {
	// Project operations
	Projects() ProjectClient

	// Authentication operations
	Auth() AuthClient

	// Environment operations
	Environments() EnvironmentClient

	// Experiment operations
	Experiments() ExperimentClient

	// Infrastructure operations
	Infrastructure() InfrastructureClient

	// Probe operations
	Probes() ProbeClient
}

// ClientOptions contains configuration for the API client
type ClientOptions struct {
	Endpoint string
	Username string
	Password string
}

// LitmusClient implements the Client interface
type LitmusClient struct {
	credentials          types.Credentials
	projectClient        ProjectClient
	authClient           AuthClient
	environmentClient    EnvironmentClient
	experimentClient     ExperimentClient
	infrastructureClient InfrastructureClient
	probeClient          ProbeClient
}

// NewClient creates a new Litmus API client
func NewClient(options ClientOptions) (Client, error) {
	authResp, err := apis.Auth(types.AuthInput{
		Endpoint: options.Endpoint,
		Username: options.Username,
		Password: options.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	credentials := types.Credentials{
		Endpoint: options.Endpoint,
		Token:    authResp.AccessToken,
		Username: options.Username,
	}

	client := &LitmusClient{
		credentials: credentials,
	}

	client.projectClient = &projectClient{credentials: credentials}
	client.authClient = &authClient{credentials: credentials}
	client.environmentClient = &environmentClient{credentials: credentials}
	client.experimentClient = &experimentClient{credentials: credentials}
	client.infrastructureClient = &infrastructureClient{credentials: credentials}
	client.probeClient = &probeClient{credentials: credentials}

	return client, nil
}

// Projects returns a ProjectClient for project operations
func (c *LitmusClient) Projects() ProjectClient {
	return c.projectClient
}

// Auth returns an AuthClient for authentication operations
func (c *LitmusClient) Auth() AuthClient {
	return c.authClient
}

// Environments returns an EnvironmentClient for environment operations
func (c *LitmusClient) Environments() EnvironmentClient {
	return c.environmentClient
}

// Experiments returns an ExperimentClient for experiment operations
func (c *LitmusClient) Experiments() ExperimentClient {
	return c.experimentClient
}

// Infrastructure returns an InfrastructureClient for infrastructure operations
func (c *LitmusClient) Infrastructure() InfrastructureClient {
	return c.infrastructureClient
}

// Probes returns a ProbeClient for probe operations
func (c *LitmusClient) Probes() ProbeClient {
	return c.probeClient
}
