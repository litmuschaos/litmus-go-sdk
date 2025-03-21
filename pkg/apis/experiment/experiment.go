/*
Copyright © 2025 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package experiment

import (
	"fmt"

	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	"github.com/litmuschaos/litmus-go-sdk/pkg/utils"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
)

// CreateExperiment sends GraphQL API request for creating a Experiment
func CreateExperiment(pid string, requestData model.SaveChaosExperimentRequest, cred types.Credentials) (RunExperimentResponse, error) {
	// Save the experiment
	_, err := utils.SendGraphQLRequest[SaveExperimentData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		SaveExperimentQuery,
		struct {
			ProjectID                  string
			SaveChaosExperimentRequest model.SaveChaosExperimentRequest
		}{
			ProjectID:                  pid,
			SaveChaosExperimentRequest: requestData,
		},
		"Error in saving Chaos Experiment",
	)
	if err != nil {
		utils.LogError("Error in saving Chaos Experiment", err)
		return RunExperimentResponse{}, err
	}

	// Run the experiment
	runExperiment, err := utils.SendGraphQLRequest[RunExperimentResponse](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		RunExperimentQuery,
		struct {
			ExperimentID string
			ProjectID    string
		}{
			ExperimentID: requestData.ID,
			ProjectID:    pid,
		},
		"Error in running Chaos Experiment",
	)
	if err != nil {
		utils.LogError("Error in running Chaos Experiment", err)
		return RunExperimentResponse{}, err
	}

	return runExperiment, nil
}

func SaveExperiment(pid string, requestData model.SaveChaosExperimentRequest, cred types.Credentials) (SaveExperimentData, error) {
	return utils.SendGraphQLRequest[SaveExperimentData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		SaveExperimentQuery,
		struct {
			ProjectID                  string
			SaveChaosExperimentRequest model.SaveChaosExperimentRequest
		}{
			ProjectID:                  pid,
			SaveChaosExperimentRequest: requestData,
		},
		"Error in saving Chaos Experiment",
	)
}

func RunExperiment(pid string, eid string, cred types.Credentials) (RunExperimentResponse, error) {
	return utils.SendGraphQLRequest[RunExperimentResponse](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		RunExperimentQuery,
		struct {
			ExperimentID string
			ProjectID    string
		}{
			ExperimentID: eid,
			ProjectID:    pid,
		},
		"Error in running Chaos Experiment",
	)
}

// GetExperimentList sends GraphQL API request for fetching a list of experiments.
func GetExperimentList(pid string, in model.ListExperimentRequest, cred types.Credentials) (ExperimentListData, error) {
	return utils.SendGraphQLRequest[ExperimentListData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		ListExperimentQuery,
		struct {
			ProjectID                 string
			GetChaosExperimentRequest model.ListExperimentRequest
		}{
			ProjectID:                 pid,
			GetChaosExperimentRequest: in,
		},
		"Error in fetching Chaos Experiments",
	)
}

// GetExperimentRunsList sends GraphQL API request for fetching a list of experiment runs.
func GetExperimentRunsList(pid string, in model.ListExperimentRunRequest, cred types.Credentials) (ExperimentRunListData, error) {
	return utils.SendGraphQLRequest[ExperimentRunListData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		ListExperimentRunsQuery,
		struct {
			ProjectID                    string
			GetChaosExperimentRunRequest model.ListExperimentRunRequest
		}{
			ProjectID:                    pid,
			GetChaosExperimentRunRequest: in,
		},
		"Error in fetching Chaos Experiment runs",
	)
}

// DeleteChaosExperiment sends GraphQL API request for deleting a given Chaos Experiment.
func DeleteChaosExperiment(projectID string, experimentID *string, cred types.Credentials) (DeleteChaosExperimentData, error) {
	return utils.SendGraphQLRequest[DeleteChaosExperimentData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		DeleteExperimentQuery,
		struct {
			ProjectID    string
			ExperimentID *string
		}{
			ProjectID:    projectID,
			ExperimentID: experimentID,
		},
		"Error in deleting Chaos Experiment",
	)
}
