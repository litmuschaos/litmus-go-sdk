package experiment

import "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

// UserDetails defines the structure for user information
type UserDetails struct {
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

// SaveExperimentData represents the response data for saving an experiment
type SaveExperimentData struct {
	Message string `json:"saveChaosExperiment"`
}

// RunExperimentData represents the response data for running an experiment
type RunExperimentData struct {
	RunChaosExperiment struct {
		NotifyID string `json:"notifyID"`
	} `json:"runChaosExperiment"`
}

// ExperimentList represents the response data for listing experiments
type ExperimentList struct {
	ListExperimentDetails model.ListExperimentResponse `json:"listExperiment"`
}


// ExperimentRunsList represents the response data for listing experiment runs
type ExperimentRunsList struct {
	ListExperimentRunDetails model.ListExperimentRunResponse `json:"listExperimentRun"`
}

// ExperimentRunDetails represents the response data for a specific experiment run
type ExperimentRunDetails struct {
	ExperimentRun model.ExperimentRun `json:"getExperimentRun"`
}

// ExperimentStatusDetails represents the response data for experiment status
type ExperimentStatusDetails struct {
	ExperimentDetails model.GetExperimentResponse `json:"getExperiment"`
}

// DeleteChaosExperimentDetails represents the response data for deleting an experiment
type DeleteChaosExperimentDetails struct {
	IsDeleted bool `json:"deleteChaosExperiment"`
}

// SaveChaosExperimentGraphQLRequest represents the GraphQL request for saving an experiment
type SaveChaosExperimentGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID                  string                           `json:"projectID"`
		SaveChaosExperimentRequest model.SaveChaosExperimentRequest `json:"request"`
	} `json:"variables"`
}

// GetChaosExperimentsGraphQLRequest represents the GraphQL request for getting experiments
type GetChaosExperimentsGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		GetChaosExperimentRequest model.ListExperimentRequest `json:"request"`
		ProjectID                 string                      `json:"projectID"`
	} `json:"variables"`
}

// GetChaosExperimentRunGraphQLRequest represents the GraphQL request for getting experiment runs
type GetChaosExperimentRunGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID                    string                         `json:"projectID"`
		GetChaosExperimentRunRequest model.ListExperimentRunRequest `json:"request"`
	} `json:"variables"`
}

// DeleteChaosExperimentGraphQLRequest represents the GraphQL request for deleting an experiment
type DeleteChaosExperimentGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID       string  `json:"projectID"`
		ExperimentID    *string `json:"experimentID"`
		ExperimentRunID *string `json:"experimentRunID"`
	} `json:"variables"`
}
