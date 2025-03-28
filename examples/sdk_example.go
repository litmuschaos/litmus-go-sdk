package examples

import (
	"os"

	"github.com/litmuschaos/litmus-go-sdk/pkg/logger"
	"github.com/litmuschaos/litmus-go-sdk/pkg/sdk"
)

func CompleteSDKExample() {
	// Initialize the Litmus SDK client
	client, err := sdk.NewClient(sdk.ClientOptions{
		Endpoint: getEnv("LITMUS_ENDPOINT", "http://localhost:8080"),
		Username: getEnv("LITMUS_USERNAME", "admin"),
		Password: getEnv("LITMUS_PASSWORD", "litmus"),
	})
	if err != nil {
		logger.Fatalf("Failed to initialize client: %v", err)
	}

	// ======== Project Operations ========

	// List all projects
	projects, err := client.Projects().List()
	if err != nil {
		logger.Fatalf("Failed to list projects: %v", err)
	}
	logger.InfoWithValues("Projects", map[string]interface{}{
		"projects": projects.Data.Projects,
	})

	// Create a new project
	newProject, err := client.Projects().Create("my-new-sdk-project")
	if err != nil {
		logger.Fatalf("Failed to create project: %v", err)
	}
	logger.InfoWithValues("Created project", map[string]interface{}{
		"name": newProject.Data.Name,
		"id":   newProject.Data.ID,
	})

	// ======== Environment Operations ========

	// List environments
	environments, err := client.Environments().List()
	if err != nil {
		logger.Fatalf("Failed to list environments: %v", err)
	}
	logger.InfoWithValues("Environments", map[string]interface{}{
		"environments": environments,
	})

	// Create environment
	envConfig := map[string]interface{}{
		"namespace": "litmus",
		"type":      "kubernetes",
	}
	newEnv, err := client.Environments().Create("production", envConfig)
	if err != nil {
		logger.Fatalf("Failed to create environment: %v", err)
	}
	logger.InfoWithValues("Created environment", map[string]interface{}{
		"environment": newEnv,
	})

	// ======== Experiment Operations ========

	// List experiments
	experiments, err := client.Experiments().List()
	if err != nil {
		logger.Fatalf("Failed to list experiments: %v", err)
	}
	logger.InfoWithValues("Experiments", map[string]interface{}{
		"experiments": experiments,
	})

	// Create experiment
	expConfig := map[string]interface{}{
		"type":      "pod-delete",
		"target":    "deployment/nginx",
		"namespace": "default",
		"duration":  30,
	}
	newExp, err := client.Experiments().Create("nginx-availability-test", expConfig)
	if err != nil {
		logger.Fatalf("Failed to create experiment: %v", err)
	}
	logger.InfoWithValues("Created experiment", map[string]interface{}{
		"experiment": newExp,
	})

	// Run experiment
	runResult, err := client.Experiments().Run("experiment-id-123")
	if err != nil {
		logger.Fatalf("Failed to run experiment: %v", err)
	}
	logger.InfoWithValues("Experiment run", map[string]interface{}{
		"result": runResult,
	})

	// ======== Infrastructure Operations ========

	// List infrastructure
	infraList, err := client.Infrastructure().List()
	if err != nil {
		logger.Fatalf("Failed to list infrastructure: %v", err)
	}
	logger.InfoWithValues("Infrastructure", map[string]interface{}{
		"infrastructure": infraList,
	})

	// Create infrastructure
	infraConfig := map[string]interface{}{
		"type":     "kubernetes",
		"provider": "gcp",
		"region":   "us-central1",
	}
	newInfra, err := client.Infrastructure().Create("gcp-cluster", infraConfig)
	if err != nil {
		logger.Fatalf("Failed to create infrastructure: %v", err)
	}
	logger.InfoWithValues("Created infrastructure", map[string]interface{}{
		"infrastructure": newInfra,
	})

	// Connect to infrastructure
	connectParams := map[string]string{
		"kubeconfig": "/path/to/kubeconfig",
	}
	err = client.Infrastructure().Connect("infra-id-123", connectParams)
	if err != nil {
		logger.Fatalf("Failed to connect to infrastructure: %v", err)
	}
	logger.Info("Connected to infrastructure")

	// ======== Probe Operations ========

	// List probes
	probes, err := client.Probes().List()
	if err != nil {
		logger.Fatalf("Failed to list probes: %v", err)
	}
	logger.InfoWithValues("Probes", map[string]interface{}{
		"probes": probes,
	})

	// Create probe
	probeConfig := map[string]interface{}{
		"type":    "http",
		"url":     "https://api.example.com/health",
		"method":  "GET",
		"timeout": 5,
	}
	newProbe, err := client.Probes().Create("api-health-check", probeConfig)
	if err != nil {
		logger.Fatalf("Failed to create probe: %v", err)
	}
	logger.InfoWithValues("Created probe", map[string]interface{}{
		"probe": newProbe,
	})

	// Execute probe
	executeParams := map[string]string{
		"headers": "Content-Type: application/json",
	}
	probeResult, err := client.Probes().Execute("probe-id-123", executeParams)
	if err != nil {
		logger.Fatalf("Failed to execute probe: %v", err)
	}
	logger.InfoWithValues("Probe execution", map[string]interface{}{
		"result": probeResult,
	})
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
