package environment

import (
	"fmt"

	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	"github.com/litmuschaos/litmus-go-sdk/pkg/utils"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
)

// CreateEnvironment connects the Infra with the given details
func CreateEnvironment(pid string, request models.CreateEnvironmentRequest, cred types.Credentials) (CreateEnvironmentData, error) {
	return utils.SendGraphQLRequest[CreateEnvironmentData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		CreateEnvironmentQuery,
		struct {
			ProjectID string                          `json:"projectID"`
			Request   models.CreateEnvironmentRequest `json:"request"`
		}{
			ProjectID: pid,
			Request:   request,
		},
		"Error in Creating Chaos Environment",
	)
}

func ListChaosEnvironments(pid string, cred types.Credentials) (ListEnvironmentData, error) {
	if pid == "" {
		return ListEnvironmentData{}, fmt.Errorf("project ID cannot be empty")
	}
	
	return utils.SendGraphQLRequest[ListEnvironmentData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		ListEnvironmentQuery,
		struct {
			ProjectID string                        `json:"projectID"`
			Request   models.ListEnvironmentRequest `json:"request"`
		}{
			ProjectID: pid,
			Request:   models.ListEnvironmentRequest{},
		},
		"Error in Getting Chaos Environment List",
	)
}

func GetChaosEnvironment(pid string, envid string, cred types.Credentials) (GetEnvironmentData, error) {
	return utils.SendGraphQLRequest[GetEnvironmentData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		GetEnvironmentQuery,
		struct {
			ProjectID     string `json:"projectID"`
			EnvironmentID string `json:"environmentID"`
		}{
			ProjectID:     pid,
			EnvironmentID: envid,
		},
		"Error in Getting Chaos Environment",
	)
}

func DeleteEnvironment(pid string, envid string, cred types.Credentials) (DeleteChaosEnvironmentData, error) {
	return utils.SendGraphQLRequest[DeleteChaosEnvironmentData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		DeleteEnvironmentQuery,
		struct {
			ProjectID     string `json:"projectID"`
			EnvironmentID string `json:"environmentID"`
		}{
			ProjectID:     pid,
			EnvironmentID: envid,
		},
		"Error in Deleting Chaos Environment",
	)
}