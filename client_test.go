package main

import (
	"testing"

	"github.com/google/uuid"
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
	project, err := client.CreateProject("test-project-" + projectID)
	assert.NoError(t, err, "Failed to create project")
	assert.NotNil(t, project, "Created project should not be nil")
}
