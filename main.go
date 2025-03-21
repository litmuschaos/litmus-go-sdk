package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/kelseyhightower/envconfig"
	"github.com/litmuschaos/litmus-go-sdk/pkg/apis"
	"github.com/litmuschaos/litmus-go-sdk/pkg/logger"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
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
	// Initialize client
	client, err := NewLitmusClient(utils.Config.Endpoint, utils.Config.Username, utils.Config.Password)
	if err != nil {
		logger.Fatalf("Failed to initialize client: %v", err)
	}

	// Example API call 1: List projects
	projects, err := client.ListProjects()
	if err != nil {
		logger.Fatalf("Failed to list projects: %v", err)
	}
	logger.InfoWithValues("Projects", map[string]interface{}{
		"projects": projects,
	})

	// Example API call 2: Create a project
	newProject, err := client.CreateProject("my-new-project")
	if err != nil {
		logger.Fatalf("Failed to create project: %v", err)
	}
	logger.InfoWithValues("Created project", map[string]interface{}{
		"name": newProject.Data.Name,
		"id":   newProject.Data.ID,
	})

	// Example API call 3: Get project details
	details, err := client.GetProjectDetails()
	if err != nil {
		logger.Fatalf("Failed to get project details: %v", err)
	}
	logger.InfoWithValues("Project details", map[string]interface{}{
		"details": details,
	})
}

// LitmusClient provides methods to interact with Litmus Chaos API
type LitmusClient struct {
	credentials types.Credentials
}

// NewLitmusClient creates and authenticates a new client
func NewLitmusClient(endpoint, username, password string) (*LitmusClient, error) {
	authResp, err := apis.Auth(types.AuthInput{
		Endpoint: endpoint,
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return &LitmusClient{
		credentials: types.Credentials{
			Endpoint: endpoint,
			Token:    authResp.AccessToken,
		},
	}, nil
}

// ListProjects retrieves all projects
func (c *LitmusClient) ListProjects() (apis.ListProjectResponse, error) {
	return apis.ListProject(c.credentials)
}

// CreateProject creates a new project
func (c *LitmusClient) CreateProject(name string) (apis.CreateProjectResponse, error) {
	return apis.CreateProjectRequest(name, c.credentials)
}

// GetProjectDetails retrieves detailed information about projects
func (c *LitmusClient) GetProjectDetails() (apis.ProjectDetails, error) {
	return apis.GetProjectDetails(c.credentials)
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
