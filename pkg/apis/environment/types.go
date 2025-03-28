package environment

import model "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

type CreateEnvironmentResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data CreateEnvironmentData `json:"data"`
}

type CreateEnvironmentData struct {
	EnvironmentDetails model.Environment `json:"createEnvironment"`
}

type GetEnvironmentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data GetEnvironment `json:"data"`
}

type GetEnvironment struct {
	EnvironmentDetails model.Environment `json:"getEnvironment"`
}

type ListEnvironmentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data EnvironmentsList `json:"data"`
}

type EnvironmentsList struct {
	ListEnvironmentDetails model.ListEnvironmentResponse `json:"listEnvironments"`
}

type DeleteChaosEnvironmentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data DeleteChaosEnvironmentDetails `json:"data"`
}

type DeleteChaosEnvironmentDetails struct {
	DeleteChaosEnvironment string `json:"deleteChaosExperiment"`
}
