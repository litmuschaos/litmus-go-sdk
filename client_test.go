package main

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/litmuschaos/litmus-go-sdk/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewLitmusClient(t *testing.T) {
	client, err := NewLitmusClient("http://127.0.0.1:35961", "admin", "LitmusChaos123@")
	assert.NoError(t, err, "Failed to create Litmus client")
	assert.NotNil(t, client, "Client should not be nil")
}

func TestListProjects(t *testing.T) {
	client, err := NewLitmusClient("http://127.0.0.1:35961", "admin", "LitmusChaos123@")
	assert.NoError(t, err)

	projects, err := client.ListProjects()
	assert.NoError(t, err, "Failed to list projects")
	assert.NotNil(t, projects, "Projects list should not be nil")
}

func TestCreateProject(t *testing.T) {
	client, err := NewLitmusClient("http://127.0.0.1:35961", "admin", "LitmusChaos123@")
	assert.NoError(t, err)

	projectID := uuid.New().String()
	projectName := fmt.Sprintf("test-project-%s", projectID)
	project, err := client.CreateProject(projectName)
	assert.NoError(t, err, "Failed to create project")
	assert.NotNil(t, project, "Created project should not be nil")

	logger.InfoWithValues("Project created", map[string]interface{}{
		"name": projectName,
		"id":   project.Data.ID,
	})
}
