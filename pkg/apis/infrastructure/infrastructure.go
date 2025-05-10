/*
Copyright © 2025 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a1 copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package infrastructure

import (
	"encoding/json"
	"fmt"

	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	"github.com/litmuschaos/litmus-go-sdk/pkg/utils"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
)

// GetInfraList lists the Chaos Infrastructure connected to the specified project
func GetInfraList(c types.Credentials, pid string, request models.ListInfraRequest) (InfraData, error) {
	if pid == "" {
		return InfraData{}, fmt.Errorf("project ID cannot be empty")
	}
	
	return utils.SendGraphQLRequest[InfraData](
		fmt.Sprintf("%s%s", c.ServerEndpoint, utils.GQLAPIPath),
		c.Token,
		ListInfraQuery,
		struct {
			ProjectID        string                  `json:"projectID"`
			ListInfraRequest models.ListInfraRequest `json:"request"`
		}{
			ProjectID:        pid,
			ListInfraRequest: request,
		},
		"Error in Getting Chaos Infrastructure List",
	)
}

// ConnectInfra connects the Infra with the given details
func ConnectInfra(infra types.Infra, cred types.Credentials) (InfraConnectionData, error) {
	registerRequest := CreateRegisterInfraRequest(infra)

	// Add node selector if provided
	if infra.NodeSelector != "" {
		registerRequest.NodeSelector = &infra.NodeSelector
	}

	// Add tolerations if provided
	if infra.Tolerations != "" {
		var toleration []*models.Toleration
		err := json.Unmarshal([]byte(infra.Tolerations), &toleration)
		if err != nil {
			utils.LogError("Error unmarshaling tolerations", err)
			return InfraConnectionData{}, fmt.Errorf("error unmarshaling tolerations: %v", err)
		}
		registerRequest.Tolerations = toleration
	}

	return utils.SendGraphQLRequest[InfraConnectionData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		RegisterInfraQuery,
		struct {
			ProjectID           string                      `json:"projectID"`
			RegisterInfraRequest models.RegisterInfraRequest `json:"request"`
		}{
			ProjectID:           infra.ProjectID,
			RegisterInfraRequest: registerRequest,
		},
		"Error in registering Chaos Infrastructure",
	)

}
func CreateRegisterInfraRequest(infra types.Infra) models.RegisterInfraRequest {
	return models.RegisterInfraRequest{
		Name:               infra.InfraName,
		InfraScope:         infra.Mode,
		Description:        &infra.Description,
		PlatformName:       infra.PlatformName,
		EnvironmentID:      infra.EnvironmentID,
		InfrastructureType: models.InfrastructureTypeKubernetes,
		InfraNamespace:     &infra.Namespace,
		ServiceAccount:     &infra.ServiceAccount,
		InfraNsExists:      &infra.NsExists,
		InfraSaExists:      &infra.SAExists,
		SkipSsl:            &infra.SkipSSL,
	}
}

// DisconnectInfra sends GraphQL API request for disconnecting Chaos Infra(s).
func DisconnectInfra(projectID string, infraID string, cred types.Credentials) (DisconnectInfraData, error) {
	return utils.SendGraphQLRequest[DisconnectInfraData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		DisconnectInfraQuery,
		struct {
			ProjectID string `json:"projectID"`
			InfraID   string `json:"infraID"`
		}{
			ProjectID: projectID,
			InfraID:   infraID,
		},
		"Error in disconnecting Chaos Infrastructure",
	)
}

func GetServerVersion(endpoint string) (ServerVersionResponse, error) {
	return utils.SendGraphQLRequest[ServerVersionResponse](
		fmt.Sprintf("%s%s", endpoint, utils.GQLAPIPath),
		"", // No token required for version check
		ServerVersionQuery,
		struct{}{}, // No variables needed
		"Error getting server version",
	)
}
