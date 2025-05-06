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
func CreateExperiment(pid string, requestData model.SaveChaosExperimentRequest, cred types.Credentials) (RunExperimentData, error) {
	// Save the experiment
	_, err := utils.SendGraphQLRequest[SaveExperimentData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		SaveExperimentQuery,
		struct {
			ProjectID string                           `json:"projectID"`
			Request   model.SaveChaosExperimentRequest `json:"request"`
		}{
			ProjectID: pid,
			Request:   requestData,
		},
		"Error in saving Chaos Experiment",
	)
	if err != nil {
		utils.LogError("Error in saving Chaos Experiment", err)
		return RunExperimentData{}, err
	}

	// Run the experiment
	runExperiment, err := utils.SendGraphQLRequest[RunExperimentData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		RunExperimentQuery,
		struct {
			ExperimentID string `json:"experimentID"`
			ProjectID    string `json:"projectID"`
		}{
			ExperimentID: requestData.ID,
			ProjectID:    pid,
		},
		"Error in running Chaos Experiment",
	)
	if err != nil {
		utils.LogError("Error in running Chaos Experiment", err)
		return RunExperimentData{}, err
	}

	return runExperiment, nil
}

func SaveExperiment(pid string, requestData model.SaveChaosExperimentRequest, cred types.Credentials) (SaveExperimentData, error) {
	return utils.SendGraphQLRequest[SaveExperimentData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		SaveExperimentQuery,
		struct {
			ProjectID string                           `json:"projectID"`
			Request   model.SaveChaosExperimentRequest `json:"request"`
		}{
			ProjectID: pid,
			Request:   requestData,
		},
		"Error in saving Chaos Experiment",
	)
}

func RunExperiment(pid string, eid string, cred types.Credentials) (RunExperimentData, error) {
	return utils.SendGraphQLRequest[RunExperimentData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		RunExperimentQuery,
		struct {
			ExperimentID string `json:"experimentID"`
			ProjectID    string `json:"projectID"`
		}{
			ExperimentID: eid,
			ProjectID:    pid,
		},
		"Error in running Chaos Experiment",
	)
}

// GetExperimentList sends GraphQL API request for fetching a list of experiments.
func GetExperimentList(pid string, in model.ListExperimentRequest, cred types.Credentials) (ExperimentList, error) {
	return utils.SendGraphQLRequest[ExperimentList](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		ListExperimentQuery,
		struct {
			ProjectID string                      `json:"projectID"`
			Request   model.ListExperimentRequest `json:"request"`
		}{
			ProjectID: pid,
			Request:   in,
		},
		"Error in fetching Chaos Experiments",
	)
}

// GetExperimentRunsList sends GraphQL API request for fetching a list of experiment runs.
func GetExperimentRunsList(pid string, in model.ListExperimentRunRequest, cred types.Credentials) (ExperimentRunsList, error) {
	return utils.SendGraphQLRequest[ExperimentRunsList](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		ListExperimentRunsQuery,
		struct {
			ProjectID string                         `json:"projectID"`
			Request   model.ListExperimentRunRequest `json:"request"`
		}{
			ProjectID: pid,
			Request:   in,
		},
		"Error in fetching Chaos Experiment runs",
	)
}

// DeleteChaosExperiment sends GraphQL API request for deleting a given Chaos Experiment.
func DeleteChaosExperiment(pid string, eid *string, cred types.Credentials) (DeleteChaosExperimentDetails, error) {
	return utils.SendGraphQLRequest[DeleteChaosExperimentDetails](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		DeleteExperimentQuery,
		struct {
			ProjectID       string  `json:"projectID"`
			ExperimentID    *string `json:"experimentID"`
			ExperimentRunID *string `json:"experimentRunID"`
		}{
			ProjectID:    pid,
			ExperimentID: eid,
		},
		"Error in deleting Chaos Experiment",
	)
}

// GetExperimentRun sends GraphQL API request for getting a specific experiment run.
func GetExperimentRun(pid string, runID string, cred types.Credentials) (ExperimentRunDetails, error) {
	return utils.SendGraphQLRequest[ExperimentRunDetails](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		GetExperimentRunQuery,
		struct {
			ProjectID       string `json:"projectID"`
			ExperimentRunID string `json:"experimentRunID"`
		}{
			ProjectID:       pid,
			ExperimentRunID: runID,
		},
		"Error in fetching Chaos Experiment Run",
	)
}
