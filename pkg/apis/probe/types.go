package probe

import model "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

type GetProbeResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data GetProbeResponseData `json:"data"`
}

type GetProbeResponseData struct {
	GetProbe model.Probe `json:"getProbe"`
}

type ListProbeResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data ListProbeResponseData `json:"data"`
}

type ListProbeResponseData struct {
	Probes []model.Probe `json:"listProbes"`
}

type DeleteProbeResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data DeleteProbeResponseData `json:"data"`
}

type DeleteProbeResponseData struct {
	DeleteProbe bool `json:"deleteProbe"`
}

type GetProbeYAMLResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data GetProbeYAMLResponseData `json:"data"`
}

type GetProbeYAMLResponseData struct {
	GetProbeYAML string `json:"getProbeYAML"`
}
