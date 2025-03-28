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
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"

	"github.com/fatih/color"
	"github.com/litmuschaos/litmus-go-sdk/pkg/logger"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
)

var (
	Red     = color.New(color.FgRed)
	White_B = color.New(color.FgWhite, color.Bold)
	White   = color.New(color.FgWhite)
)

func PrintError(err error) {
	if err != nil {
		Red.Println(err)
		os.Exit(1)
	}
}

// SendHTTPRequest is a utility function to send HTTP requests and handle common response patterns
func SendHTTPRequest(endpoint, token string, payload []byte, method string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unmatched status code: %s", string(bodyBytes))
	}

	return bodyBytes, nil
}

// ProcessResponse is a utility function to process HTTP responses and unmarshal JSON data
func ProcessResponse[T any](bodyBytes []byte, errorPrefix string) (T, error) {
	var result T
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return result, fmt.Errorf("%s: %v", errorPrefix, err)
	}
	return result, nil
}

// GraphQLRequest represents a GraphQL request structure
type GraphQLRequest struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

// GraphQLResponse represents a GraphQL response structure with errors
type GraphQLResponse[T any] struct {
	Data   T `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// SendGraphQLRequest is a utility function to send GraphQL requests and handle responses
func SendGraphQLRequest[T any](endpoint, token string, query string, variables interface{}, errorPrefix string) (T, error) {
	var result T
	gqlReq := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	payload, err := json.Marshal(gqlReq)
	if err != nil {
		return result, fmt.Errorf("%s: error marshaling request: %v", errorPrefix, err)
	}

	bodyBytes, err := SendHTTPRequest(endpoint, token, payload, string(types.Post))
	if err != nil {
		return result, fmt.Errorf("%s: %v", errorPrefix, err)
	}

	var response GraphQLResponse[T]
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return result, fmt.Errorf("%s: error unmarshaling response: %v", errorPrefix, err)
	}

	if len(response.Errors) > 0 {
		return result, fmt.Errorf("%s: GraphQL error: %s", errorPrefix, response.Errors[0].Message)
	}

	// Initialize the result to avoid nil pointer issues during tests
	// This is useful for test scenarios where we expect structured data
	initializeEmptyStruct(&response.Data)

	return response.Data, nil
}

// initializeEmptyStruct recursively initializes maps and slices within a structure to avoid nil values
func initializeEmptyStruct(v interface{}) {
	if v == nil {
		return
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return
	}

	val = val.Elem()
	if !val.IsValid() {
		return
	}

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.Kind() == reflect.Ptr && field.IsNil() && field.CanSet() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			if field.CanAddr() {
				initializeEmptyStruct(field.Addr().Interface())
			}
		}
	case reflect.Map:
		if val.IsNil() && val.CanSet() {
			val.Set(reflect.MakeMap(val.Type()))
		}
	case reflect.Slice:
		if val.IsNil() && val.CanSet() {
			val.Set(reflect.MakeSlice(val.Type(), 0, 0))
		}
	}
}

// LogError is a utility function to log errors with consistent formatting
func LogError(message string, err error) {
	logger.ErrorWithValues(message, map[string]interface{}{
		"error": err.Error(),
	})
}
