package environment

import model "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

// Type for CreateEnvironment - already correct
type CreateEnvironmentData struct {
    CreateEnvironment model.Environment `json:"createEnvironment"`
}

// Type for ListChaosEnvironments
type ListEnvironmentData struct {
    ListEnvironments model.ListEnvironmentResponse `json:"listEnvironments"`
}

// Type for GetChaosEnvironment
type GetEnvironmentData struct {
    GetEnvironment model.Environment `json:"getEnvironment"`
}

// Type for DeleteEnvironment
type DeleteChaosEnvironmentData struct {
    DeleteEnvironment string `json:"deleteEnvironment"`
}