package main

import (
	"runtime"

	"github.com/kelseyhightower/envconfig"
	"github.com/litmuschaos/litmus-go-sdk/pkg/logger"
	"github.com/litmuschaos/litmus-go-sdk/pkg/sdk"
	"github.com/litmuschaos/litmus-go-sdk/pkg/utils"
)

func init() {
	logger.Infof("go version: %s", runtime.Version())
	logger.Infof("go os/arch: %s/%s", runtime.GOOS, runtime.GOARCH)

	err := envconfig.Process("", &utils.Config)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func main() {
	// Initialize client using the new SDK interface
	client, err := sdk.NewClient(sdk.ClientOptions{
		Endpoint: utils.Config.Endpoint,
		Username: utils.Config.Username,
		Password: utils.Config.Password,
	})
	if err != nil {
		logger.Fatalf("Failed to initialize client: %v", err)
	}

	// Example API call 1: List projects using the SDK
	projects, err := client.Projects().List()
	if err != nil {
		logger.Fatalf("Failed to list projects: %v", err)
	}
	logger.InfoWithValues("Projects", map[string]interface{}{
		"projects": projects,
	})

	// Example API call 2: Create a project using the SDK
	newProject, err := client.Projects().Create("my-new-project")
	if err != nil {
		logger.Fatalf("Failed to create project: %v", err)
	}
	logger.InfoWithValues("Created project", map[string]interface{}{
		"name": newProject.Data.Name,
		"id":   newProject.Data.ID,
	})

	// Example API call 3: Get project details using the SDK
	details, err := client.Projects().GetDetails()
	if err != nil {
		logger.Fatalf("Failed to get project details: %v", err)
	}
	logger.InfoWithValues("Project details", map[string]interface{}{
		"details": details,
	})
}
