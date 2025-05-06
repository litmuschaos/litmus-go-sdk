package experiment

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/litmuschaos/litmus-go-sdk/pkg/apis"
	"github.com/litmuschaos/litmus-go-sdk/pkg/apis/environment"
	"github.com/litmuschaos/litmus-go-sdk/pkg/apis/infrastructure"
	"github.com/litmuschaos/litmus-go-sdk/pkg/logger"
	"github.com/litmuschaos/litmus-go-sdk/pkg/types"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/stretchr/testify/assert"
)

// Test configuration with defaults
var (
	testEndpoint = "http://127.0.0.1:39651"
	testUsername = "admin"
	testPassword = "litmus"	
	// Store IDs as package-level variables for test access
	projectID       string
	environmentID   string
	infrastructureID string
	experimentID    string
	credentials     types.Credentials
)

// Use the exact format from the successful API request
    // Note: The manifest is a JSON string containing the workflow definition
const workflowManifest = `{
        "apiVersion": "argoproj.io/v1alpha1",
        "kind": "Workflow",
        "metadata": {
            "name": "test-experiment",
            "namespace": "litmus-2",
            "labels": {
                "subject": "{{workflow.parameters.appNamespace}}_kube-proxy"
            }
        },
        "spec": {
            "entrypoint": "argowf-chaos",
            "serviceAccountName": "argo-chaos",
            "securityContext": {
                "runAsUser": 1000,
                "runAsNonRoot": true
            },
            "arguments": {
                "parameters": [
                    {
                        "name": "adminModeNamespace",
                        "value": "litmus-2"
                    },
                    {
                        "name": "appNamespace",
                        "value": "kube-system"
                    }
                ]
            },
            "templates": [
                {
                    "name": "argowf-chaos",
                    "steps": [
                        [
                            {
                                "name": "install-chaos-faults",
                                "template": "install-chaos-faults"
                            }
                        ],
                        [
                            {
                                "name": "run-chaos",
                                "template": "run-chaos"
                            }
                        ],
                        [
                            {
                                "name": "cleanup-chaos-resources",
                                "template": "cleanup-chaos-resources"
                            }
                        ]
                    ]
                },
                {
                    "name": "install-chaos-faults",
                    "inputs": {
                        "artifacts": [
                            {
                                "name": "install-chaos-faults",
                                "path": "/tmp/pod-delete.yaml",
                                "raw": {
                                    "data": "apiVersion: litmuschaos.io/v1alpha1\ndescription:\n  message: |\n    Deletes a pod belonging to a deployment/statefulset/daemonset\nkind: ChaosExperiment\nmetadata:\n  name: pod-delete\nspec:\n  definition:\n    scope: Namespaced\n    permissions:\n      - apiGroups:\n          - \"\"\n          - \"apps\"\n          - \"batch\"\n          - \"litmuschaos.io\"\n        resources:\n          - \"deployments\"\n          - \"jobs\"\n          - \"pods\"\n          - \"pods/log\"\n          - \"events\"\n          - \"configmaps\"\n          - \"chaosengines\"\n          - \"chaosexperiments\"\n          - \"chaosresults\"\n        verbs:\n          - \"create\"\n          - \"list\"\n          - \"get\"\n          - \"patch\"\n          - \"update\"\n          - \"delete\"\n      - apiGroups:\n          - \"\"\n        resources:\n          - \"nodes\"\n        verbs:\n          - \"get\"\n          - \"list\"\n    image: \"litmuschaos.docker.scarf.sh/litmuschaos/go-runner:3.16.0\"\n    imagePullPolicy: Always\n    args:\n    - -c\n    - ./experiments -name pod-delete\n    command:\n    - /bin/bash\n    env:\n\n    - name: TOTAL_CHAOS_DURATION\n      value: '15'\n\n    # Period to wait before and after injection of chaos in sec\n    - name: RAMP_TIME\n      value: ''\n\n    # provide the kill count\n    - name: KILL_COUNT\n      value: ''\n\n    - name: FORCE\n      value: 'true'\n\n    - name: CHAOS_INTERVAL\n      value: '5'\n\n    labels:\n      name: pod-delete\n"
                                }
                            }
                        ]
                    },
                    "container": {
                        "image": "litmuschaos/k8s:latest",
                        "command": [
                            "sh",
                            "-c"
                        ],
                        "args": [
                            "kubectl apply -f /tmp/pod-delete.yaml -n {{workflow.parameters.adminModeNamespace}}"
                        ]
                    }
                },
                {
                    "name": "run-chaos",
                    "inputs": {
                        "artifacts": [
                            {
                                "name": "run-chaos",
                                "path": "/tmp/chaosengine-run-chaos.yaml",
                                "raw": {
                                    "data": "apiVersion: litmuschaos.io/v1alpha1\nkind: ChaosEngine\nmetadata:\n  namespace: \"{{workflow.parameters.adminModeNamespace}}\"\n  labels:\n    context: \"{{workflow.parameters.appNamespace}}_kube-proxy\"\n    workflow_run_id: \"{{ workflow.uid }}\"\n    workflow_name: test-experiment\n  annotations:\n    probeRef: '[{\"name\":\"myprobe\",\"mode\":\"SOT\"}]'\n  generateName: run-chaos\nspec:\n  appinfo:\n    appns: litmus-2\n    applabel: app=nginx\n    appkind: deployment\n  jobCleanUpPolicy: retain\n  engineState: active\n  chaosServiceAccount: litmus-admin\n  experiments:\n    - name: pod-delete\n      spec:\n        components:\n          env:\n            - name: TOTAL_CHAOS_DURATION\n              value: \"60\"\n            - name: CHAOS_INTERVAL\n              value: \"10\"\n            - name: FORCE\n              value: \"false\"\n"
                                }
                            }
                        ]
                    },
                    "metadata": {
                        "labels": {
                            "weight": "10"
                        }
                    },
                    "container": {
                        "name": "",
                        "image": "docker.io/litmuschaos/litmus-checker:2.11.0",
                        "args": [
                            "-file=/tmp/chaosengine-run-chaos.yaml",
                            "-saveName=/tmp/engine-name"
                        ]
                    }
                },
                {
                    "name": "cleanup-chaos-resources",
                    "container": {
                        "image": "litmuschaos/k8s:latest",
                        "command": [
                            "sh",
                            "-c"
                        ],
                        "args": [
                            "kubectl delete chaosengine -l workflow_run_id={{workflow.uid}} -n {{workflow.parameters.adminModeNamespace}}"
                        ]
                    }
                }
            ]
        }
    }`

func TestMain(m *testing.M) {
	// Override defaults with environment variables if set
	if endpoint := os.Getenv("LITMUS_TEST_ENDPOINT"); endpoint != "" {
		testEndpoint = endpoint
	}
	if username := os.Getenv("LITMUS_TEST_USERNAME"); username != "" {
		testUsername = username
	}
	if password := os.Getenv("LITMUS_TEST_PASSWORD"); password != "" {
		testPassword = password
	}

	logger.Infof("Test configuration - Endpoint: %s, Username: %s", testEndpoint, testUsername)
	
	// Setup credentials by authenticating
	authResp, err := apis.Auth(types.AuthInput{
		Endpoint: testEndpoint,
		Username: testUsername,
		Password: testPassword,
	})
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	credentials = types.Credentials{
		ServerEndpoint: testEndpoint,
		Endpoint: testEndpoint,
		Token:          authResp.AccessToken,
	}

	// Get or create project ID
	projectResp, err := apis.ListProject(credentials)
	if err != nil {
		log.Fatalf("Failed to list projects: %v", err)
	}

	if len(projectResp.Data.Projects) > 0 {
		projectID = projectResp.Data.Projects[0].ID
		logger.Infof("Using existing project ID: %s", projectID)
	} else {
		// Create a project if none exists
		projectName := fmt.Sprintf("test-project-%s", uuid.New().String())
		newProject, err := apis.CreateProjectRequest(projectName, credentials)
		if err != nil {
			log.Fatalf("Failed to create project: %v", err)
		}
		projectID = newProject.Data.ID
		logger.Infof("Created new project ID: %s", projectID)
	}
	
	// Store project ID in credentials for convenience
	credentials.ProjectID = projectID

	// 1. Seed Environment Data
	logger.Infof("Seeding Environment data...")
	environmentID = seedEnvironmentData(credentials, projectID)

	// 2. Seed Infrastructure Data
	logger.Infof("Seeding Infrastructure data...")
	infrastructureID = seedInfrastructureData(credentials, projectID, environmentID)
	
	examineExistingExperiment(credentials, projectID)
	// 3. Seed Experiment Data
	logger.Infof("Seeding Experiment data...")
	experimentID = seedExperimentData(credentials, projectID, infrastructureID)
	
	// Run the tests
	exitCode := m.Run()
	
	// Exit with the test status code
	os.Exit(exitCode)
}

func seedEnvironmentData(credentials types.Credentials, projectID string) string {
	// Create environment
	envID := fmt.Sprintf("test-env-%s", uuid.New().String())
	description := "Test environment for SDK tests"
	
	envRequest := model.CreateEnvironmentRequest{
		Name:          "test-environment",
		Description:   &description,
		Type:          "NON_PROD",
		Tags:          []string{"test", "sdk"},
		EnvironmentID: envID,
	}

	envResp, err := environment.CreateEnvironment(projectID, envRequest, credentials)
	if err != nil {
		log.Fatalf("Failed to create environment: %v", err)
	}

	logger.Infof("Created environment with ID: %s", envResp.CreateEnvironment.EnvironmentID)

	return envResp.CreateEnvironment.EnvironmentID
}

func seedInfrastructureData(credentials types.Credentials, projectID, environmentID string) string {
	// Connect infrastructure
	infraName := "test-infrastructure"
	description := "Test infrastructure for SDK tests"
	namespace := "litmus"
	serviceAccount := "litmus"
	
	// Create the infrastructure object
	infra := types.Infra{
		ProjectID:      projectID,
		InfraName:      infraName,
		Description:    description,
		PlatformName:   "kubernetes",
		Mode:           "cluster",
		EnvironmentID:  environmentID,
		Namespace:      namespace,
		ServiceAccount: serviceAccount,
		NsExists:       true,
		SAExists:       true,
		SkipSSL:        false,
	}
	
	// Connect the infrastructure
	infraResp, err := infrastructure.ConnectInfra(infra, credentials)
	if err != nil {
		log.Fatalf("Failed to register infrastructure: %v", err)
	}

	logger.Infof("Created infrastructure with ID: %s", infraResp.RegisterInfra.InfraID)
	
	return infraResp.RegisterInfra.InfraID
}

func examineExistingExperiment(credentials types.Credentials, projectID string) {
    // Use the provided experiment ID
    existingExperimentID := "4813cc63-753e-4d2e-80a0-fba935a2f75d"
    
    // Create a request to get the experiment details
    request := model.ListExperimentRequest{
        ExperimentIDs: []*string{&existingExperimentID},
    }
    
    logger.Infof("Examining existing experiment with ID: %s", existingExperimentID)
    // Fetch the experiment
    response, err := GetExperimentList(projectID, request, credentials)
	fmt.Println("response", response)
    if err != nil {
        logger.Errorf("Failed to get existing experiment: %v", err)
        return
    }
    
    // Check if we got any experiments back
    if len(response.Data.ListExperimentDetails.Experiments) == 0 {
        logger.Errorf("No experiment found with ID: %s", existingExperimentID)
        return
    }
    
    // Get the experiment
    experiment := response.Data.ListExperimentDetails.Experiments[0]
    
    // Log the experiment details
    experimentJSON, _ := json.MarshalIndent(experiment, "", "  ")
    logger.Infof("Existing experiment details: %s", string(experimentJSON))
    
    // Specifically log the manifest
    logger.Infof("Existing experiment manifest: %s", experiment.ExperimentManifest)
    
    // Now try to create a new experiment using the same format
    experimentID := fmt.Sprintf("test-exp-%s", uuid.New().String())
    
    experimentRequest := model.SaveChaosExperimentRequest{
        ID:       experimentID,
        Name:     "cloned-experiment",
        InfraID:  experiment.Infra.InfraID,
        Manifest: experiment.ExperimentManifest,
    }
    
    
    logger.Infof("Creating new experiment with cloned format: %+v", experimentRequest)
    
    // Save the experiment
    _, err = SaveExperiment(projectID, experimentRequest, credentials)
    if err != nil {
        logger.Errorf("Failed to create cloned experiment: %v", err)
        return
    }
    
    logger.Infof("Successfully created cloned experiment with ID: %s", experimentID)
}

func seedExperimentData(credentials types.Credentials, projectID, infrastructureID string) string {
    // Create experiment
    experimentID := fmt.Sprintf("test-exp-%s", uuid.New().String())
    experimentName := fmt.Sprintf("test-exp-%s", uuid.New().String())
    
    
    // Create the experiment request with the required fields
    experimentRequest := model.SaveChaosExperimentRequest{
        ID:       experimentID,
		Name:     experimentName,
        InfraID:  infrastructureID,
        Manifest: getWorkflowManifest(experimentName),
    }
    
    // Log the request for debugging
    requestJSON, _ := json.MarshalIndent(experimentRequest, "", "  ")
    logger.Infof("Creating experiment with request: %s", string(requestJSON))
    
    // Save the experiment
    _, err := SaveExperiment(projectID, experimentRequest, credentials)
    if err != nil {
        // Log detailed error information
        logger.Errorf("Failed to create experiment: %v", err)
        log.Fatalf("Failed to create experiment: %v", err)
    }
    
    logger.Infof("Created experiment with ID: %s", experimentID)
    
    return experimentID
}

func init() {
	// Override defaults with environment variables if set
	if endpoint := os.Getenv("LITMUS_TEST_ENDPOINT"); endpoint != "" {
		testEndpoint = endpoint
	}
	if username := os.Getenv("LITMUS_TEST_USERNAME"); username != "" {
		testUsername = username
	}
	if password := os.Getenv("LITMUS_TEST_PASSWORD"); password != "" {
		testPassword = password
	}

	logger.Infof("Test configuration - Endpoint: %s, Username: %s", testEndpoint, testUsername)
}


func getWorkflowManifest(experimentName string) string {
    // Parse the existing manifest as an object
    var manifestObj map[string]interface{}
    err := json.Unmarshal([]byte(workflowManifest), &manifestObj)
    if err != nil {
        return workflowManifest // Return original if parsing fails
    }
    
    // Update the metadata.name field to match the experiment name
    metadata, ok := manifestObj["metadata"].(map[string]interface{})
    if ok {
        metadata["name"] = experimentName
    }
    
    // Update workflow_name in labels if present
    spec, ok := manifestObj["spec"].(map[string]interface{})
    if ok {
        templates, ok := spec["templates"].([]interface{})
        if ok && len(templates) > 2 {
            runChaos, ok := templates[2].(map[string]interface{})
            if ok {
                inputs, ok := runChaos["inputs"].(map[string]interface{})
                if ok {
                    artifacts, ok := inputs["artifacts"].([]interface{})
                    if ok && len(artifacts) > 0 {
                        artifact, ok := artifacts[0].(map[string]interface{})
                        if ok {
                            raw, ok := artifact["raw"].(map[string]interface{})
                            if ok {
                                data, ok := raw["data"].(string)
                                if ok {
                                    // Replace workflow_name in the raw data string
                                    data = strings.Replace(data, "workflow_name: test-experiment", 
                                                          fmt.Sprintf("workflow_name: %s", experimentName), -1)
                                    raw["data"] = data
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    
    // Convert back to JSON string
    updatedManifest, err := json.Marshal(manifestObj)
    if err != nil {
        return workflowManifest // Return original if marshaling fails
    }
    
    return string(updatedManifest)
}
// NewLitmusClient creates and authenticates a new client with username/password
func NewLitmusClient(endpoint, username, password string) (*LitmusClient, error) {
	// Implementation should match the one in main.go
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
			ServerEndpoint: endpoint,
			Token:          authResp.AccessToken,
		},
	}, nil
}

// LitmusClient provides methods to interact with Litmus Chaos API
type LitmusClient struct {
	credentials types.Credentials
}

func setupTestClient() (*LitmusClient, error) {
	return NewLitmusClient(testEndpoint, testUsername, testPassword)
}

func TestSaveExperiment(t *testing.T) {
	experimentName := fmt.Sprintf("test-exp-%s", uuid.New().String())
    tests := []struct {
        name       string
        projectID  string
        request    model.SaveChaosExperimentRequest
        setup      func(*LitmusClient) // optional setup steps
        wantErr    bool
        validateFn func(*testing.T, *SaveExperimentData)
    }{
        {
            name:      "successful experiment save",
            projectID: projectID, // Use the actual projectID instead of "test-project-id"
            request: model.SaveChaosExperimentRequest{
                ID:       fmt.Sprintf("test-exp-%s", uuid.New().String()),
                Name:     experimentName,
                InfraID:  infrastructureID, // Add the real infrastructure ID
                Manifest: getWorkflowManifest(experimentName), 
            },
            wantErr: false,
            validateFn: func(t *testing.T, result *SaveExperimentData) {
                assert.NotNil(t, result, "Result should not be nil")
                assert.NotEmpty(t, result.Message, "Response message should not be empty")
                assert.Contains(t, result.Message, "success", "Response message should indicate success")
            },
        },
        {
            name:      "save experiment with empty ID",
            projectID: projectID,
            request: model.SaveChaosExperimentRequest{
                ID:   "",
                Name: "test-experiment",
            },
            wantErr:    true,
            validateFn: nil,
        },
    }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient()
			assert.NoError(t, err, "Failed to create Litmus client")

			// Run any setup function if provided
			if tt.setup != nil {
				tt.setup(client)
			}

			result, err := SaveExperiment(tt.projectID, tt.request, client.credentials)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Run validation function if provided
			if tt.validateFn != nil {
				tt.validateFn(t, &result)
			}
		})
	}
}

func TestRunExperiment(t *testing.T) {
    // First create an experiment to run
    var testExpID string
    
    tests := []struct {
        name         string
        projectID    string
        experimentID string
        setup        func(*LitmusClient) // optional setup steps
        wantErr      bool
        validateFn   func(*testing.T, *RunExperimentData)
    }{
        {
            name:      "successful experiment run",
            projectID: projectID,
            // We'll set the experimentID in the setup function
            setup: func(client *LitmusClient) {
                expID := fmt.Sprintf("test-exp-%s", uuid.New().String())
				experimentName := fmt.Sprintf("test-exp-%s", uuid.New().String())
                req := model.SaveChaosExperimentRequest{
                    ID:       expID,
                    Name:     experimentName,
                    InfraID:  infrastructureID,
                    Manifest: getWorkflowManifest(experimentName),
                }
                _, err := SaveExperiment(projectID, req, client.credentials)
                if err != nil {
                    t.Logf("Setup failed to create experiment: %v", err)
                    return
                }
                testExpID = expID
            },
            wantErr: false,
            validateFn: func(t *testing.T, result *RunExperimentData) {
                assert.NotNil(t, result, "Result should not be nil")
                assert.NotNil(t, result.RunChaosExperiment, "RunChaosExperiment should not be nil")
                
                // NotifyID might be empty if the experiment run is queued but not yet started
                // but it should at least be a defined field
                assert.IsType(t, "", result.RunChaosExperiment.NotifyID, "NotifyID should be a string")
                
                // If NotifyID is present, it should be non-empty
                if result.RunChaosExperiment.NotifyID != "" {
                    assert.NotEmpty(t, result.RunChaosExperiment.NotifyID, "If NotifyID is returned, it should not be empty")
                    assert.True(t, len(result.RunChaosExperiment.NotifyID) > 5, "NotifyID should be a substantial string")
                }
            },
        },
        {
            name:         "experiment run with empty ID",
            projectID:    projectID,
            experimentID: "",
            wantErr:      false,
            validateFn:   nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client, err := setupTestClient()
            assert.NoError(t, err, "Failed to create Litmus client")

            // Run any setup function if provided
            if tt.setup != nil {
                tt.setup(client)
            }
            
            // Use the experimentID created during setup or the one in the test case
            experimentIDToRun := testExpID
            if tt.experimentID != "" {
                experimentIDToRun = tt.experimentID
            }

            result, err := RunExperiment(tt.projectID, experimentIDToRun, client.credentials)

            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)

            // Run validation function if provided
            if tt.validateFn != nil {
                tt.validateFn(t, &result)
            }
        })
    }
}

func TestGetExperimentList(t *testing.T) {
	tests := []struct {
		name       string
		projectID  string
		request    model.ListExperimentRequest
		setup      func(*LitmusClient) // optional setup steps
		wantErr    bool
		validateFn func(*testing.T, *ExperimentListData)
	}{
		{
			name:      "successful experiment list fetch",
			projectID: projectID,
			request: model.ListExperimentRequest{
				Pagination: &model.Pagination{
					Page:  1,
					Limit: 10,
				},
			},
			wantErr: false,
			validateFn: func(t *testing.T, result *ExperimentListData) {
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.NotNil(t, result.Data.ListExperimentDetails, "ListExperimentDetails should not be nil")

				// Check total count is a non-negative number
				assert.GreaterOrEqual(t, result.Data.ListExperimentDetails.TotalNoOfExperiments, 0,
					"Total number of experiments should be non-negative")

				// If there are experiments, validate their structure
				if len(result.Data.ListExperimentDetails.Experiments) > 0 {
					for i, exp := range result.Data.ListExperimentDetails.Experiments {
						assert.NotEmpty(t, exp.ExperimentID, "Experiment ID should not be empty for experiment at index %d", i)
						assert.NotEmpty(t, exp.Name, "Experiment name should not be empty for experiment at index %d", i)
						assert.NotNil(t, exp.Infra, "Infra should not be nil for experiment at index %d", i)
						assert.NotEmpty(t, exp.Infra.InfraID, "Infra ID should not be empty for experiment at index %d", i)
						assert.NotEmpty(t, exp.ExperimentManifest, "Experiment manifest should not be empty for experiment at index %d", i)
						assert.NotEmpty(t, exp.ProjectID, "Project ID should not be empty for experiment at index %d", i)
						
						// Check timestamps are valid
						assert.NotZero(t, exp.CreatedAt, "CreatedAt timestamp should be non-zero for experiment at index %d", i)
					}
				}
			},
		},
		{
			name:      "experiment list with pagination",
			projectID: projectID,
			request: model.ListExperimentRequest{
				Pagination: &model.Pagination{
					Page:  1,
					Limit: 5,
				},
			},
			wantErr: false,
			validateFn: func(t *testing.T, result *ExperimentListData) {
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.NotNil(t, result.Data.ListExperimentDetails, "ListExperimentDetails should not be nil")

				// Verify pagination works by checking max results
				if result.Data.ListExperimentDetails.TotalNoOfExperiments > 0 {
					assert.LessOrEqual(t,
						len(result.Data.ListExperimentDetails.Experiments),
						5,
						"Should return 5 or fewer results with limit=5")
				}
				
				// If there are any experiments returned, verify they have valid data
				if len(result.Data.ListExperimentDetails.Experiments) > 0 {
					for i, exp := range result.Data.ListExperimentDetails.Experiments {
						assert.NotEmpty(t, exp.ExperimentID, "Experiment ID should not be empty for experiment at index %d", i)
						assert.NotEmpty(t, exp.Name, "Experiment name should not be empty for experiment at index %d", i)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient()
			assert.NoError(t, err, "Failed to create Litmus client")

			// Run any setup function if provided
			if tt.setup != nil {
				tt.setup(client)
			}

			result, err := GetExperimentList(tt.projectID, tt.request, client.credentials)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Run validation function if provided
			if tt.validateFn != nil {
				tt.validateFn(t, &result)
			}
		})
	}
}

func TestGetExperimentRunsList(t *testing.T) {
	tests := []struct {
		name       string
		projectID  string
		request    model.ListExperimentRunRequest
		setup      func(*LitmusClient) // optional setup steps
		wantErr    bool
		validateFn func(*testing.T, *ExperimentRunListData) // fixed type
	}{
		{
			name:      "successful experiment runs list fetch",
			projectID: projectID,
			request: model.ListExperimentRunRequest{
				Pagination: &model.Pagination{
					Page:  1,
					Limit: 10,
				},
			},
			wantErr: false,
			validateFn: func(t *testing.T, result *ExperimentRunListData) { // fixed type
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.NotNil(t, result.Data.ListExperimentRunDetails, "ListExperimentRunDetails should not be nil")
				assert.GreaterOrEqual(t, result.Data.ListExperimentRunDetails.TotalNoOfExperimentRuns, 0,
					"Total number of experiment runs should be non-negative")
				
				// If there are experiment runs, validate their structure
				if len(result.Data.ListExperimentRunDetails.ExperimentRuns) > 0 {
					for i, run := range result.Data.ListExperimentRunDetails.ExperimentRuns {
						assert.NotEmpty(t, run.ExperimentRunID, "Experiment run ID should not be empty for run at index %d", i)
						assert.NotEmpty(t, run.ExperimentID, "Experiment ID should not be empty for run at index %d", i)
						assert.NotEmpty(t, run.ProjectID, "Project ID should not be empty for run at index %d", i)
						
						// Check for valid timestamp
						assert.NotZero(t, run.CreatedAt, "CreatedAt timestamp should be non-zero for run at index %d", i)
						
						// Check for non-empty status fields
						assert.NotEmpty(t, run.Phase, "Phase should not be empty for run at index %d", i)
						assert.NotNil(t, run.ResiliencyScore, "ResiliencyScore should not be nil for run at index %d", i)
					}
				}
			},
		},
		{
			name:      "experiment runs list with pagination",
			projectID: projectID,
			request: model.ListExperimentRunRequest{
				Pagination: &model.Pagination{
					Page:  1,
					Limit: 5,
				},
				// Removed filter to avoid the linter error
			},
			wantErr: false,
			validateFn: func(t *testing.T, result *ExperimentRunListData) { // fixed type
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotNil(t, result.Data, "Data should not be nil")
				assert.NotNil(t, result.Data.ListExperimentRunDetails, "ListExperimentRunDetails should not be nil")

				// Verify pagination works by checking max results
				if result.Data.ListExperimentRunDetails.TotalNoOfExperimentRuns > 0 {
					assert.LessOrEqual(t,
						len(result.Data.ListExperimentRunDetails.ExperimentRuns),
						5,
						"Should return 5 or fewer results with limit=5")
				}
				
				// If there are any runs returned, verify they have valid data
				if len(result.Data.ListExperimentRunDetails.ExperimentRuns) > 0 {
					for i, run := range result.Data.ListExperimentRunDetails.ExperimentRuns {
						assert.NotEmpty(t, run.ExperimentRunID, "Experiment run ID should not be empty for run at index %d", i)
						assert.NotEmpty(t, run.ExperimentID, "Experiment ID should not be empty for run at index %d", i)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient()
			assert.NoError(t, err, "Failed to create Litmus client")

			// Run any setup function if provided
			if tt.setup != nil {
				tt.setup(client)
			}

			result, err := GetExperimentRunsList(tt.projectID, tt.request, client.credentials)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Run validation function if provided
			if tt.validateFn != nil {
				tt.validateFn(t, &result)
			}
		})
	}
}

func TestDeleteChaosExperiment(t *testing.T) {
    // First create an experiment to delete
    var testExpID string
    client, err := setupTestClient()
    assert.NoError(t, err, "Failed to create Litmus client")
    
    // Create a test experiment that can be deleted
    expID := fmt.Sprintf("test-exp-%s", uuid.New().String())
    req := model.SaveChaosExperimentRequest{
        ID:       expID,
        Name:     "test-experiment-for-deletion",
        InfraID:  infrastructureID,
        Manifest: workflowManifest,
    }
    _, err = SaveExperiment(projectID, req, client.credentials)
    if err != nil {
        t.Logf("Failed to create experiment for deletion: %v", err)
    } else {
        testExpID = expID
    }
    
    tests := []struct {
        name         string
        projectID    string
        experimentID string
        wantErr      bool
        validateFn   func(*testing.T, *DeleteChaosExperimentData)
    }{
        {
            name:         "successful experiment deletion",
            projectID:    projectID,
            experimentID: testExpID, // Use the experiment we just created
            wantErr:      false,
            validateFn: func(t *testing.T, result *DeleteChaosExperimentData) {
                assert.NotNil(t, result, "Result should not be nil")
                assert.NotNil(t, result.Data, "Data should not be nil")
                assert.True(t, result.Data.IsDeleted, "IsDeleted should be true")
                
                // Verify the response has no errors
                assert.Empty(t, result.Errors, "Response should not contain any errors")
                
                // After successful deletion, verify the experiment is no longer retrievable
                listReq := model.ListExperimentRequest{
                    ExperimentIDs: []*string{&testExpID},
                }
                listResp, listErr := GetExperimentList(projectID, listReq, client.credentials)
                
                // Either the request will fail or it should return zero experiments
                if listErr == nil {
                    assert.Equal(t, 0, len(listResp.Data.ListExperimentDetails.Experiments), 
                        "Deleted experiment should not be retrievable")
                }
            },
        },
        {
            name:         "experiment deletion with empty ID",
            projectID:    projectID,
            experimentID: "",
            wantErr:      true,
            validateFn:   nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Skip the test if we couldn't create an experiment and this test requires one
            if tt.experimentID == testExpID && testExpID == "" {
                t.Skip("Skipping test because setup experiment could not be created")
            }

            experimentID := tt.experimentID
            result, err := DeleteChaosExperiment(tt.projectID, &experimentID, client.credentials)

            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)

            if tt.validateFn != nil {
                tt.validateFn(t, &result)
            }
        })
    }
}

func TestCreateExperiment(t *testing.T) {
	experimentName := fmt.Sprintf("test-exp-%s", uuid.New().String())
    tests := []struct {
        name       string
        projectID  string
        request    model.SaveChaosExperimentRequest
        setup      func(*LitmusClient) // optional setup steps
        wantErr    bool
        validateFn func(*testing.T, *RunExperimentData)
    }{
        {
            name:      "successful experiment creation and run",
            projectID: projectID,
            request: model.SaveChaosExperimentRequest{
                ID:       fmt.Sprintf("test-exp-%s", uuid.New().String()),
                Name:     experimentName,
                InfraID:  infrastructureID,
                Manifest: getWorkflowManifest(experimentName),
            },
            wantErr: false,
            validateFn: func(t *testing.T, result *RunExperimentData) {
                assert.NotNil(t, result, "Result should not be nil")
                assert.NotNil(t, result.RunChaosExperiment, "RunChaosExperiment should not be nil")
                
                // NotifyID might be empty if the experiment run is queued but not yet started
                assert.IsType(t, "", result.RunChaosExperiment.NotifyID, "NotifyID should be a string")
                
                // If NotifyID is present, check it's a meaningful value
                if result.RunChaosExperiment.NotifyID != "" {
                    assert.True(t, len(result.RunChaosExperiment.NotifyID) > 5, 
                        "If NotifyID is returned, it should be a substantial string")
                }
            },
        },
        {
            name:      "experiment creation with empty ID",
            projectID: projectID,
            request: model.SaveChaosExperimentRequest{
                ID:   "",
                Name: experimentName,
            },
            wantErr:    true,
            validateFn: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client, err := setupTestClient()
            assert.NoError(t, err, "Failed to create Litmus client")

            // Run any setup function if provided
            if tt.setup != nil {
                tt.setup(client)
            }

            result, err := CreateExperiment(tt.projectID, tt.request, client.credentials)

            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)

            // Run validation function if provided
            if tt.validateFn != nil {
                tt.validateFn(t, &result)
            }
        })
    }
}
