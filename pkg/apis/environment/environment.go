package environment

import (
	"encoding/json"
	"fmt"

	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	"github.com/litmuschaos/litmus-go-sdk/pkg/utils"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
)

// CreateEnvironment connects the Infra with the given details
func CreateEnvironment(pid string, request models.CreateEnvironmentRequest, cred types.Credentials) (CreateEnvironmentResponse, error) {
    var response CreateEnvironmentResponse
    
    gqlReq := struct {
        Query     string      `json:"query"`
        Variables interface{} `json:"variables"`
    }{
        Query: CreateEnvironmentQuery,
        Variables: struct {
            ProjectID string                          `json:"projectID"`
            Request   models.CreateEnvironmentRequest `json:"request"`
        }{
            ProjectID: pid,
            Request:   request,
        },
    }
    
    payload, err := json.Marshal(gqlReq)
    if err != nil {
        return response, fmt.Errorf("Error marshaling request: %v", err)
    }
    
    bodyBytes, err := utils.SendHTTPRequest(
        fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
        cred.Token,
        payload, 
        string(types.Post),
    )
    if err != nil {
        return response, fmt.Errorf("Error sending request: %v", err)
    }
    
    if err := json.Unmarshal(bodyBytes, &response); err != nil {
        return response, fmt.Errorf("Error unmarshaling response: %v", err)
    }
    
    if len(response.Errors) > 0 {
        return response, fmt.Errorf("GraphQL error: %s", response.Errors[0].Message)
    }
    
    return response, nil
}

func ListChaosEnvironments(pid string, cred types.Credentials) (ListEnvironmentData, error) {
	if pid == "" {
		return ListEnvironmentData{}, fmt.Errorf("project ID cannot be empty")
	}
	
	return utils.SendGraphQLRequest[ListEnvironmentData](
		fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
		cred.Token,
		ListEnvironmentQuery,
		struct {
			ProjectID string                        `json:"projectID"`
			Request   models.ListEnvironmentRequest `json:"request"`
		}{
			ProjectID: pid,
			Request:   models.ListEnvironmentRequest{},
		},
		"Error in Getting Chaos Environment List",
	)
}

func GetChaosEnvironment(pid string, envid string, cred types.Credentials) (GetEnvironmentData, error) {
    var response GetEnvironmentData
    
    gqlReq := struct {
        Query     string      `json:"query"`
        Variables interface{} `json:"variables"`
    }{
        Query: GetEnvironmentQuery,
        Variables: struct {
            ProjectID     string `json:"projectID"`
            EnvironmentID string `json:"environmentID"`
        }{
            ProjectID:     pid,
            EnvironmentID: envid,
        },
    }
    
    payload, err := json.Marshal(gqlReq)
    if err != nil {
        return response, fmt.Errorf("Error marshaling request: %v", err)
    }
    
    bodyBytes, err := utils.SendHTTPRequest(
        fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
        cred.Token,
        payload, 
        string(types.Post),
    )
    if err != nil {
        return response, fmt.Errorf("Error sending request: %v", err)
    }
    
    fmt.Printf("Raw GraphQL response:\n%s\n", string(bodyBytes))
    
    if err := json.Unmarshal(bodyBytes, &response); err != nil {
        return response, fmt.Errorf("Error unmarshaling response: %v", err)
    }
    
    fmt.Printf("Unmarshaled response: %+v\n", response)
    
    if len(response.Errors) > 0 {
        return response, fmt.Errorf("GraphQL error: %s", response.Errors[0].Message)
    }
    
    return response, nil
}

func DeleteEnvironment(pid string, envid string, cred types.Credentials) (DeleteChaosEnvironmentData, error) {
    var response DeleteChaosEnvironmentData
    
    gqlReq := struct {
        Query     string      `json:"query"`
        Variables interface{} `json:"variables"`
    }{
        Query: DeleteEnvironmentQuery,
        Variables: struct {
            ProjectID     string `json:"projectID"`
            EnvironmentID string `json:"environmentID"`
        }{
            ProjectID:     pid,
            EnvironmentID: envid,
        },
    }
    
    payload, err := json.Marshal(gqlReq)
    if err != nil {
        return response, fmt.Errorf("Error marshaling request: %v", err)
    }
    
    bodyBytes, err := utils.SendHTTPRequest(
        fmt.Sprintf("%s%s", cred.ServerEndpoint, utils.GQLAPIPath),
        cred.Token,
        payload, 
        string(types.Post),
    )
    if err != nil {
        return response, fmt.Errorf("Error sending request: %v", err)
    }
    
    fmt.Printf("Raw GraphQL response:\n%s\n", string(bodyBytes))
    
    if err := json.Unmarshal(bodyBytes, &response); err != nil {
        return response, fmt.Errorf("Error unmarshaling response: %v", err)
    }
    
    fmt.Printf("Unmarshaled response: %+v\n", response)
    
    if len(response.Errors) > 0 {
        return response, fmt.Errorf("GraphQL error: %s", response.Errors[0].Message)
    }
    
    return response, nil
}
