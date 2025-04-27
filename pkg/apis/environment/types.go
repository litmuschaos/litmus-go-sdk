package environment

import model "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

type CreateEnvironmentResponse struct {
    Errors []struct {
        Message string   `json:"message"`
        Path    []string `json:"path"`
    } `json:"errors"`
    Data struct {
        CreateEnvironment struct {
            EnvironmentID string `json:"environmentID"`
            Name          string `json:"name"`
        } `json:"createEnvironment"`
    } `json:"data"`
}

type CreateEnvironmentData struct {
	EnvironmentDetails model.Environment `json:"createEnvironment"`
}

// For GetEnvironmentData
type GetEnvironmentData struct {
    Errors []struct {
        Message string   `json:"message"`
        Path    []string `json:"path"`
    } `json:"errors"`
    Data struct {
        GetEnvironment struct {
            EnvironmentID string     `json:"environmentID"`
            Name          string     `json:"name"`
            CreatedAt     string     `json:"createdAt"`
            UpdatedAt     string     `json:"updatedAt"`
            CreatedBy     struct {
                Username string `json:"username"`
            } `json:"createdBy"`
            UpdatedBy struct {
                Username string `json:"username"`
            } `json:"updatedBy"`
            InfraIDs []string `json:"infraIDs"`
            Type     string   `json:"type"`
            Tags     []string `json:"tags"`
        } `json:"getEnvironment"`
    } `json:"data"`
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

// For DeleteChaosEnvironmentData
type DeleteChaosEnvironmentData struct {
    Errors []struct {
        Message string   `json:"message"`
        Path    []string `json:"path"`
    } `json:"errors"`
    Data struct {
        DeleteEnvironment string `json:"deleteEnvironment"`
    } `json:"data"`
}

type DeleteChaosEnvironmentDetails struct {
	DeleteChaosEnvironment string `json:"deleteChaosExperiment"`
}
