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
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/litmuschaos/litmus-go-sdk/pkg/apis"
	"github.com/litmuschaos/litmus-go-sdk/pkg/logger"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	"github.com/litmuschaos/litmus-go-sdk/pkg/utils"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
)

// CreateExperiment sends GraphQL API request for creating a Experiment
func CreateExperiment(pid string, requestData model.SaveChaosExperimentRequest, cred types.Credentials) (RunExperimentResponse, error) {

	// Query to Save the Experiment
	var gqlReq SaveChaosExperimentGraphQLRequest

	gqlReq.Query = SaveExperimentQuery
	gqlReq.Variables.ProjectID = pid
	gqlReq.Variables.SaveChaosExperimentRequest = requestData

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return RunExperimentResponse{}, err
	}

	resp, err := apis.SendRequest(
		apis.SendRequestParams{
			Endpoint: cred.ServerEndpoint + utils.GQLAPIPath,
			Token:    cred.Token,
		},
		query,
		string(types.Post),
	)
	if err != nil {
		return RunExperimentResponse{}, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()
	if err != nil {
		logger.ErrorWithValues("Error in saving Chaos Experiment", map[string]interface{}{
			"error": err.Error(),
		})
		return RunExperimentResponse{}, errors.New("Error in saving Chaos Experiment: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var savedExperiment SaveExperimentData

		err = json.Unmarshal(bodyBytes, &savedExperiment)

		if err != nil {
			logger.ErrorWithValues("Error in saving Chaos Experiment", map[string]interface{}{
				"error": err.Error(),
			})
			return RunExperimentResponse{}, errors.New("Error in saving Chaos Experiment: " + err.Error())
		}

		// Errors present
		if len(savedExperiment.Errors) > 0 {
			logger.ErrorWithValues("Error in saving Chaos Experiment", map[string]interface{}{
				"error": savedExperiment.Errors[0].Message,
			})
			return RunExperimentResponse{}, errors.New(savedExperiment.Errors[0].Message)
		}

	} else {
		logger.ErrorWithValues("Error in saving Chaos Experiment", map[string]interface{}{
			"status_code": resp.StatusCode,
		})
		return RunExperimentResponse{}, errors.New("error in saving Chaos Experiment")
	}

	// Query to Run the Chaos Experiment
	runQuery := `{"query":"mutation{ \n runChaosExperiment(experimentID:  \"` + requestData.ID + `\", projectID:  \"` + pid + `\"){\n notifyID \n}}"}`
	resp, err = apis.SendRequest(apis.SendRequestParams{Endpoint: cred.ServerEndpoint + utils.GQLAPIPath, Token: cred.Token}, []byte(runQuery), string(types.Post))

	if err != nil {
		logger.ErrorWithValues("Error in Running Chaos Experiment", map[string]interface{}{
			"error": err.Error(),
		})
		return RunExperimentResponse{}, errors.New("Error in Running Chaos Experiment: " + err.Error())
	}

	bodyBytes, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logger.ErrorWithValues("Error in Running Chaos Experiment", map[string]interface{}{
			"error": err.Error(),
		})
		return RunExperimentResponse{}, errors.New("Error in Running Chaos Experiment: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var runExperiment RunExperimentResponse
		err = json.Unmarshal(bodyBytes, &runExperiment)
		if err != nil {
			logger.ErrorWithValues("Error in Running Chaos Experiment", map[string]interface{}{
				"error": err.Error(),
			})
			return RunExperimentResponse{}, errors.New("Error in Running Chaos Experiment: " + err.Error())
		}

		if len(runExperiment.Errors) > 0 {
			logger.ErrorWithValues("Error in Running Chaos Experiment", map[string]interface{}{
				"error": runExperiment.Errors[0].Message,
			})
			return RunExperimentResponse{}, errors.New(runExperiment.Errors[0].Message)
		}
		return runExperiment, nil
	} else {
		logger.ErrorWithValues("Error in Running Chaos Experiment", map[string]interface{}{
			"status_code": resp.StatusCode,
		})
		return RunExperimentResponse{}, err
	}
}

func SaveExperiment(pid string, requestData model.SaveChaosExperimentRequest, cred types.Credentials) (SaveExperimentData, error) {

	// Query to Save the Experiment
	var gqlReq SaveChaosExperimentGraphQLRequest

	gqlReq.Query = SaveExperimentQuery
	gqlReq.Variables.ProjectID = pid
	gqlReq.Variables.SaveChaosExperimentRequest = requestData

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return SaveExperimentData{}, err
	}

	resp, err := apis.SendRequest(
		apis.SendRequestParams{
			Endpoint: cred.ServerEndpoint + utils.GQLAPIPath,
			Token:    cred.Token,
		},
		query,
		string(types.Post),
	)
	if err != nil {
		return SaveExperimentData{}, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()
	if err != nil {
		return SaveExperimentData{}, errors.New("Error in saving Chaos Experiment: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var savedExperiment SaveExperimentData

		err = json.Unmarshal(bodyBytes, &savedExperiment)

		if err != nil {
			return SaveExperimentData{}, errors.New("Error in saving Chaos Experiment: " + err.Error())
		}

		// Errors present
		if len(savedExperiment.Errors) > 0 {
			return SaveExperimentData{}, errors.New(savedExperiment.Errors[0].Message)
		}
		return savedExperiment, nil

	} else {
		return SaveExperimentData{}, errors.New("error in saving Chaos Experiment")
	}

}

func RunExperiment(pid string, eid string, cred types.Credentials) (RunExperimentResponse, error) {
	var err error
	runQuery := `{"query":"mutation{ \n runChaosExperiment(experimentID:  \"` + eid + `\", projectID:  \"` + pid + `\"){\n notifyID \n}}"}`

	resp, err := apis.SendRequest(apis.SendRequestParams{Endpoint: cred.ServerEndpoint + utils.GQLAPIPath, Token: cred.Token}, []byte(runQuery), string(types.Post))

	if err != nil {
		return RunExperimentResponse{}, errors.New("Error in Running Chaos Experiment: " + err.Error())
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return RunExperimentResponse{}, errors.New("Error in Running Chaos Experiment: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var runExperiment RunExperimentResponse
		err = json.Unmarshal(bodyBytes, &runExperiment)
		if err != nil {
			return RunExperimentResponse{}, errors.New("Error in Running Chaos Experiment: " + err.Error())
		}

		if len(runExperiment.Errors) > 0 {
			return RunExperimentResponse{}, errors.New(runExperiment.Errors[0].Message)
		}
		return runExperiment, nil
	} else {
		return RunExperimentResponse{}, err
	}
}

// GetExperimentList sends GraphQL API request for fetching a list of experiments.
func GetExperimentList(pid string, in model.ListExperimentRequest, cred types.Credentials) (ExperimentListData, error) {

	var gqlReq GetChaosExperimentsGraphQLRequest
	var err error

	gqlReq.Query = ListExperimentQuery
	gqlReq.Variables.GetChaosExperimentRequest = in
	gqlReq.Variables.ProjectID = pid

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return ExperimentListData{}, err
	}

	resp, err := apis.SendRequest(
		apis.SendRequestParams{
			Endpoint: cred.ServerEndpoint + utils.GQLAPIPath,
			Token:    cred.Token,
		},
		query,
		string(types.Post),
	)
	if err != nil {
		return ExperimentListData{}, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ExperimentListData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var experimentList ExperimentListData
		err = json.Unmarshal(bodyBytes, &experimentList)
		if err != nil {
			return ExperimentListData{}, err
		}

		if len(experimentList.Errors) > 0 {
			return ExperimentListData{}, errors.New(experimentList.Errors[0].Message)
		}

		return experimentList, nil
	} else {
		return ExperimentListData{}, errors.New("Error while fetching the Chaos Experiments")
	}
}

// GetExperimentRunsList sends GraphQL API request for fetching a list of experiment runs.
func GetExperimentRunsList(pid string, in model.ListExperimentRunRequest, cred types.Credentials) (ExperimentRunListData, error) {

	var gqlReq GetChaosExperimentRunGraphQLRequest
	var err error

	gqlReq.Query = ListExperimentRunsQuery
	gqlReq.Variables.ProjectID = pid
	gqlReq.Variables.GetChaosExperimentRunRequest = in

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return ExperimentRunListData{}, err
	}

	resp, err := apis.SendRequest(
		apis.SendRequestParams{
			Endpoint: cred.ServerEndpoint + utils.GQLAPIPath,
			Token:    cred.Token,
		},
		query,
		string(types.Post),
	)
	if err != nil {
		return ExperimentRunListData{}, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ExperimentRunListData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var experimentRunList ExperimentRunListData
		err = json.Unmarshal(bodyBytes, &experimentRunList)
		if err != nil {
			return ExperimentRunListData{}, err
		}

		if len(experimentRunList.Errors) > 0 {
			return ExperimentRunListData{}, errors.New(experimentRunList.Errors[0].Message)
		}

		return experimentRunList, nil
	} else {
		return ExperimentRunListData{}, errors.New("Error while fetching the Chaos Experiment runs")
	}
}

// DeleteChaosExperiment sends GraphQL API request for deleting a given Chaos Experiment.
func DeleteChaosExperiment(projectID string, experimentID *string, cred types.Credentials) (DeleteChaosExperimentData, error) {

	var gqlReq DeleteChaosExperimentGraphQLRequest
	var err error

	gqlReq.Query = DeleteExperimentQuery
	gqlReq.Variables.ProjectID = projectID
	gqlReq.Variables.ExperimentID = experimentID
	//var experiment_run_id string = ""
	//gqlReq.Variables.ExperimentRunID = &experiment_run_id

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return DeleteChaosExperimentData{}, err
	}

	resp, err := apis.SendRequest(
		apis.SendRequestParams{
			Endpoint: cred.ServerEndpoint + utils.GQLAPIPath,
			Token:    cred.Token,
		},
		query,
		string(types.Post),
	)
	if err != nil {
		return DeleteChaosExperimentData{}, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return DeleteChaosExperimentData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var deletedExperiment DeleteChaosExperimentData
		err = json.Unmarshal(bodyBytes, &deletedExperiment)
		if err != nil {
			return DeleteChaosExperimentData{}, err
		}

		if len(deletedExperiment.Errors) > 0 {
			return DeleteChaosExperimentData{}, errors.New(deletedExperiment.Errors[0].Message)
		}

		return deletedExperiment, nil
	} else {
		return DeleteChaosExperimentData{}, errors.New("Error while deleting the Chaos Experiment")
	}
}
