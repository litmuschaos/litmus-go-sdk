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

// AuthClient defines the interface for authentication operations
type AuthClient interface {
	// GetToken returns the current authentication token
	GetToken() string

	// GetCredentials returns the current credentials
	GetCredentials() types.Credentials
}

// authClient implements the AuthClient interface
type authClient struct {
	credentials types.Credentials
}

// GetToken returns the current authentication token
func (c *authClient) GetToken() string {
	return c.credentials.Token
}

// GetCredentials returns the current credentials
func (c *authClient) GetCredentials() types.Credentials {
	return c.credentials
}
