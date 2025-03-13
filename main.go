package main

import (
	"fmt"
	"log"
	"os"

	"github.com/litmuschaos/litmus-go-sdk/pkg/apis"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
)

func main() {
	// Get credentials from environment variables or flags
	endpoint := getEnv("LITMUS_ENDPOINT", "http://localhost:8080")
	username := getEnv("LITMUS_USERNAME", "admin")
	password := getEnv("LITMUS_PASSWORD", "litmus")

	// Initialize client
	client, err := NewLitmusClient(endpoint, username, password)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// Example API call 1: List projects
	projects, err := client.ListProjects()
	if err != nil {
		log.Fatalf("Failed to list projects: %v", err)
	}
	fmt.Printf("Projects: %+v\n", projects)

	// Example API call 2: Create a project
	newProject, err := client.CreateProject("my-new-project")
	if err != nil {
		log.Fatalf("Failed to create project: %v", err)
	}
	fmt.Printf("Created project: %s with ID: %s\n", newProject.Data.Name, newProject.Data.ID)

	// Example API call 3: Get project details
	details, err := client.GetProjectDetails()
	if err != nil {
		log.Fatalf("Failed to get project details: %v", err)
	}
	fmt.Printf("Project details: %+v\n", details)
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
