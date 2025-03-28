package probe

import (
	"fmt"

	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	"github.com/litmuschaos/litmus-go-sdk/pkg/utils"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
)

func GetProbeRequest(pid string, probeID string, cred types.Credentials) (GetProbeResponse, error) {
	return utils.SendGraphQLRequest[GetProbeResponse](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		GetProbeQuery,
		struct {
			ProjectID string `json:"projectID"`
			ProbeName string `json:"probeName"`
		}{
			ProjectID: pid,
			ProbeName: probeID,
		},
		"Error in getting requested probe",
	)
}

func ListProbeRequest(pid string, probetypes []*models.ProbeType, cred types.Credentials) (ListProbeResponse, error) {
	if pid == "" {
		return ListProbeResponse{}, fmt.Errorf("project ID cannot be empty")
	}
	
	return utils.SendGraphQLRequest[ListProbeResponse](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		ListProbeQuery,
		struct {
			ProjectID string                  `json:"projectID"`
			Filter    models.ProbeFilterInput `json:"filter"`
		}{
			ProjectID: pid,
			Filter: models.ProbeFilterInput{
				Type: probetypes,
			},
		},
		"Error in listing probes",
	)
}

func DeleteProbeRequest(pid string, probeid string, cred types.Credentials) (DeleteProbeResponse, error) {
	return utils.SendGraphQLRequest[DeleteProbeResponse](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		DeleteProbeQuery,
		struct {
			ProbeName string `json:"probeName"`
			ProjectID string `json:"projectID"`
		}{
			ProbeName: probeid,
			ProjectID: pid,
		},
		"Error in deleting probe",
	)
}

func GetProbeYAMLRequest(pid string, request models.GetProbeYAMLRequest, cred types.Credentials) (GetProbeYAMLResponse, error) {
	return utils.SendGraphQLRequest[GetProbeYAMLResponse](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		GetProbeYAMLQuery,
		struct {
			ProjectID string                     `json:"projectID"`
			Request   models.GetProbeYAMLRequest `json:"request"`
		}{
			ProjectID: pid,
			Request:   request,
		},
		"Error in getting probe details",
	)
}
