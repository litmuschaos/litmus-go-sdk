package experiment

import "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

type SaveExperimentResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data SaveExperimentData `json:"data"`
}

type SaveExperimentData struct {
	Message string `json:"saveChaosExperiment"`
}


type SaveChaosExperimentGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID                  string                           `json:"projectID"`
		SaveChaosExperimentRequest model.SaveChaosExperimentRequest `json:"request"`
	} `json:"variables"`
}

type RunExperimentResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data RunExperimentData `json:"data"`
}

type RunExperimentData struct {
    RunChaosExperiment struct {
        NotifyID string `json:"notifyID"`
    } `json:"runChaosExperiment"`
}

type ExperimentListData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data ExperimentList `json:"data"`
}

type ExperimentList struct {
	ListExperimentDetails model.ListExperimentResponse `json:"listExperiment"`
}

type GetChaosExperimentsGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		GetChaosExperimentRequest model.ListExperimentRequest `json:"request"`
		ProjectID                 string                      `json:"projectID"`
	} `json:"variables"`
}

type ExperimentRunListData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data ExperimentRunsList `json:"data"`
}

type ExperimentRunsList struct {
	ListExperimentRunDetails model.ListExperimentRunResponse `json:"listExperimentRun"`
}

type GetChaosExperimentRunGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID                    string                         `json:"projectID"`
		GetChaosExperimentRunRequest model.ListExperimentRunRequest `json:"request"`
	} `json:"variables"`
}

type DeleteChaosExperimentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data DeleteChaosExperimentDetails `json:"data"`
}

type DeleteChaosExperimentDetails struct {
	IsDeleted bool `json:"deleteChaosExperiment"`
}

type DeleteChaosExperimentGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID       string  `json:"projectID"`
		ExperimentID    *string `json:"experimentID"`
		ExperimentRunID *string `json:"experimentRunID"`
	} `json:"variables"`
}

type ExperimentRunData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data ExperimentRunDetails `json:"data"`
}

type ExperimentRunDetails struct {
	ExperimentRun ExperimentRunResponse `json:"getExperimentRun"`
}

type ExperimentStatusData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data ExperimentStatusDetails `json:"data"`
}

type ExperimentStatusDetails struct {
	ExperimentDetails model.GetExperimentResponse `json:"getExperiment"`
}

// UserDetails defines the structure for user information
type UserDetails struct {
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

// ExperimentRunResponse represents the response structure from an experiment run
// This is a custom implementation as we don't have direct access to the model.ExperimentRunResponse
type ExperimentRunResponse struct {
	ProjectID        string  `json:"projectID"`
	ExperimentRunID  string  `json:"experimentRunID"`
	ExperimentID     string  `json:"experimentID"`
	ExperimentName   string  `json:"experimentName"`
	Phase            string  `json:"phase"`
	ResiliencyScore  *float64 `json:"resiliencyScore,omitempty"`
	FaultsPassed     *int     `json:"faultsPassed,omitempty"`
	FaultsFailed     *int     `json:"faultsFailed,omitempty"`
	FaultsAwaited    *int     `json:"faultsAwaited,omitempty"`
	FaultsStopped    *int     `json:"faultsStopped,omitempty"`
	FaultsNa         *int     `json:"faultsNa,omitempty"`
	TotalFaults      *int     `json:"totalFaults,omitempty"`
	UpdatedAt        string   `json:"updatedAt,omitempty"`
	UpdatedBy        *UserDetails `json:"updatedBy,omitempty"`
}
