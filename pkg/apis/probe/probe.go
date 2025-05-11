package probe

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	"github.com/litmuschaos/litmus-go-sdk/pkg/utils"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
)

func GetProbeRequest(pid string, probeName string, cred types.Credentials) (GetProbeResponse, error) {
	if probeName == "" {
		return GetProbeResponse{}, fmt.Errorf("probe name cannot be empty")
	}
	return utils.SendGraphQLRequest[GetProbeResponse](
		fmt.Sprintf("%s%s", cred.Endpoint, utils.GQLAPIPath),
		cred.Token,
		GetProbeQuery,
		struct {
			ProjectID string `json:"projectID"`
			ProbeName string `json:"probeName"`
		}{
			ProjectID: pid,
			ProbeName: probeName,
		},
		"Error in getting requested probe",
	)
}

func ListProbeRequest(pid string, probeTypes []*models.ProbeType, cred types.Credentials) (ListProbeResponse, error) {
	if pid == "" {
		return ListProbeResponse{}, fmt.Errorf("projectID cannot be empty")
	}
	return utils.SendGraphQLRequest[ListProbeResponse](
		fmt.Sprintf("%s%s", cred.Endpoint, utils.GQLAPIPath),
		cred.Token,
		ListProbeQuery,
		struct {
			ProjectID string                  `json:"projectID"`
			Filter    models.ProbeFilterInput `json:"filter"`
		}{
			ProjectID: pid,
			Filter: models.ProbeFilterInput{ // Assuming ProbeFilterInput is the correct type based on usage
				Type: probeTypes,
			},
		},
		"Error in listing probes",
	)
}

func DeleteProbeRequest(pid string, probeName string, cred types.Credentials) (DeleteProbeResponse, error) {
	if probeName == "" {
		return DeleteProbeResponse{}, fmt.Errorf("probe name cannot be empty")
	}
	return utils.SendGraphQLRequest[DeleteProbeResponse](
		fmt.Sprintf("%s%s", cred.Endpoint, utils.GQLAPIPath),
		cred.Token,
		DeleteProbeQuery,
		struct {
			ProbeName string `json:"probeName"`
			ProjectID string `json:"projectID"`
		}{
			ProbeName: probeName,
			ProjectID: pid,
		},
		"Error in deleting probe",
	)
}

func GetProbeYAMLRequest(pid string, request models.GetProbeYAMLRequest, cred types.Credentials) (GetProbeYAMLResponse, error) {
	return utils.SendGraphQLRequest[GetProbeYAMLResponse](
		fmt.Sprintf("%s%s", cred.Endpoint, utils.GQLAPIPath),
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

// CreateProbe creates a new chaos probe
func CreateProbe(request ProbeRequest, projectID string, cred types.Credentials) (*Probe, error) {
	if err := validateProbeRequest(request); err != nil {
		return nil, err
	}

	query := createProbeMutation

	rawResponse, err := utils.SendGraphQLRequest[map[string]interface{}](
		fmt.Sprintf("%s%s", cred.Endpoint, utils.GQLAPIPath),
		cred.Token,
		query,
		struct {
			ProjectID string       `json:"projectID"`
			Request   ProbeRequest `json:"request"`
		}{
			ProjectID: projectID,
			Request:   request,
		},
		"Error in creating probe",
	)

	if err != nil {
		return nil, err
	}

	// Check if we have the addProbe key in response
	addProbeData, ok := rawResponse["addProbe"]
	if !ok {
		return nil, fmt.Errorf("no addProbe data in response")
	}

	// Convert the interface{} to JSON and then to our Probe struct
	jsonData, err := json.Marshal(addProbeData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal probe data: %v", err)
	}

	var probe Probe
	if err := json.Unmarshal(jsonData, &probe); err != nil {
		return nil, fmt.Errorf("failed to unmarshal probe data: %v", err)
	}

	return &probe, nil
}

func validateProbeRequest(request ProbeRequest) error {
	propertiesSet := 0
	if request.KubernetesHTTPProperties != nil {
		propertiesSet++
	}
	if request.KubernetesCMDProperties != nil {
		propertiesSet++
	}
	if request.K8SProperties != nil {
		propertiesSet++
	}
	if request.PROMProperties != nil {
		propertiesSet++
	}

	if propertiesSet == 0 {
		return errors.New("no probe properties provided")
	}
	if propertiesSet > 1 {
		return errors.New("multiple probe property types provided, only one is allowed")
	}

	switch request.Type {
	case ProbeTypeHTTPProbe:
		if request.KubernetesHTTPProperties == nil {
			return errors.New("httpProbe type requires kubernetesHTTPProperties")
		}
	case ProbeTypeCMDProbe:
		if request.KubernetesCMDProperties == nil {
			return errors.New("cmdProbe type requires kubernetesCMDProperties")
		}
	case ProbeTypeK8SProbe:
		if request.K8SProperties == nil {
			return errors.New("k8sProbe type requires k8sProperties")
		}
	case ProbeTypePROMProbe:
		if request.PROMProperties == nil {
			return errors.New("promProbe type requires promProperties")
		}
	default:
		return errors.New("invalid or unsupported probe type for validation")
	}
	return nil
}
