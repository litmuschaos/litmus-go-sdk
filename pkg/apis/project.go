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
package apis

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"

	"github.com/litmuschaos/litmus-go-sdk/pkg/utils"

	"github.com/litmuschaos/litmus-go-sdk/pkg/logger"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
)

// Common response error structure
type ErrorResponse struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type CreateProjectResponse struct {
	Data struct {
		Name string `json:"name"`
		ID   string `json:"projectID"`
	} `json:"data"`
	Errors []ErrorResponse `json:"errors"`
}

type createProjectPayload struct {
	ProjectName string `json:"projectName"`
}

func CreateProjectRequest(projectName string, cred types.Credentials) (CreateProjectResponse, error) {
	endpoint := fmt.Sprintf("%s%s/create_project", cred.Endpoint, utils.AuthAPIPath)

	payloadBytes, err := json.Marshal(createProjectPayload{
		ProjectName: projectName,
	})
	if err != nil {
		return CreateProjectResponse{}, err
	}

	bodyBytes, err := utils.SendHTTPRequest(endpoint, cred.Token, payloadBytes, string(types.Post))
	if err != nil {
		return CreateProjectResponse{}, err
	}

	var project CreateProjectResponse
	err = json.Unmarshal(bodyBytes, &project)
	if err != nil {
		return CreateProjectResponse{}, err
	}

	if len(project.Errors) > 0 {
		return CreateProjectResponse{}, errors.New(project.Errors[0].Message)
	}

	logger.InfoWithValues("Project created", map[string]interface{}{
		"project": project.Data.Name,
	})
	return project, nil
}

type ListProjectResponse struct {
	Message string `json:"message"`
	Data    struct {
		Projects []struct {
			ID        string `json:"projectID"`
			Name      string `json:"name"`
			CreatedAt int64  `json:"createdAt"`
		} `json:"projects"`
		TotalNumberOfProjects int `json:"totalNumberOfProjects"`
	} `json:"data"`
	Errors []ErrorResponse `json:"errors"`
}

func ListProject(cred types.Credentials) (ListProjectResponse, error) {
	endpoint := fmt.Sprintf("%s%s/list_projects", cred.Endpoint, utils.AuthAPIPath)

	bodyBytes, err := utils.SendHTTPRequest(endpoint, cred.Token, []byte{}, string(types.Get))
	if err != nil {
		return ListProjectResponse{}, err
	}

	var data ListProjectResponse
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return ListProjectResponse{}, err
	}

	if len(data.Errors) > 0 {
		return ListProjectResponse{}, errors.New(data.Errors[0].Message)
	}

	return data, nil
}

type ProjectDetails struct {
	Data   Data            `json:"data"`
	Errors []ErrorResponse `json:"errors"`
}

type Data struct {
	ID       string    `json:"ID"`
	Projects []Project `json:"Projects"`
}

type Member struct {
	Role     string `json:"Role"`
	UserID   string `json:"userID"`
	UserName string `json:"username"`
}

type Project struct {
	ID        string   `json:"ProjectID"`
	Name      string   `json:"Name"`
	CreatedAt int64    `json:"CreatedAt"`
	Members   []Member `json:"Members"`
}

// GetProjectDetails fetches details of the input user
func GetProjectDetails(c types.Credentials) (ProjectDetails, error) {
	token, _ := jwt.Parse(c.Token, nil)
	if token == nil {
		return ProjectDetails{}, nil
	}

	username, _ := token.Claims.(jwt.MapClaims)["username"].(string)
	endpoint := fmt.Sprintf("%s%s/get_user_with_project/%s", c.Endpoint, utils.AuthAPIPath, username)

	bodyBytes, err := utils.SendHTTPRequest(endpoint, c.Token, []byte{}, string(types.Get))
	if err != nil {
		return ProjectDetails{}, err
	}

	var project ProjectDetails
	err = json.Unmarshal(bodyBytes, &project)
	if err != nil {
		return ProjectDetails{}, err
	}

	if len(project.Errors) > 0 {
		return ProjectDetails{}, errors.New(project.Errors[0].Message)
	}

	return project, nil
}
