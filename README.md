# Litmus Go SDK

The Litmus Go SDK is a client library for interacting with the Litmus Chaos platform. This SDK provides a simple and intuitive interface for managing chaos experiments, infrastructure, and more.

## Installation

```bash
go get github.com/litmuschaos/litmus-go-sdk
```

## Usage

### Initializing the Client

```go
import (
    "github.com/litmuschaos/litmus-go-sdk/pkg/sdk"
)

// Create a new client
client, err := sdk.NewClient(sdk.ClientOptions{
    Endpoint: "https://litmus.example.com",
    Username: "admin",
    Password: "password",
})
if err != nil {
    // Handle error
}
```

## SDK Design

The Litmus Go SDK is designed with the following principles:

1. **Interface-based design**: Every feature group has a dedicated interface.
2. **Method chaining**: Operations can be chained together for better readability.
3. **Fluent API**: The API is designed to be intuitive and self-documenting.
4. **Error handling**: All operations return an error along with the result.

The SDK is organized into the following components:

- **Client**: The main entry point for the SDK, providing access to all features.
- **Projects**: Operations related to project management.
- **Auth**: Authentication and authorization operations.
- **Environments**: Operations for managing environments.
- **Experiments**: Operations for creating and running chaos experiments.
- **Infrastructure**: Operations for managing infrastructure resources.
- **Probes**: Operations for creating and executing probes.

### Working with Projects

```go
// List all projects
projects, err := client.Projects().List()
if err != nil {
    // Handle error
}

// Create a new project
newProject, err := client.Projects().Create("my-new-project")
if err != nil {
    // Handle error
}

// Get project details
details, err := client.Projects().GetDetails()
if err != nil {
    // Handle error
}
```

### Working with Environments

```go
// List all environments
environments, err := client.Environments().List()
if err != nil {
    // Handle error
}

// Create a new environment
envConfig := map[string]interface{}{
    "namespace": "litmus",
    "type":      "kubernetes",
}
newEnv, err := client.Environments().Create("production", envConfig)
if err != nil {
    // Handle error
}

// Get environment details
env, err := client.Environments().Get("env-id")
if err != nil {
    // Handle error
}
```

### Working with Experiments

```go
// List all experiments
experiments, err := client.Experiments().List()
if err != nil {
    // Handle error
}

// Create a new experiment
expConfig := map[string]interface{}{
    "type":      "pod-delete",
    "target":    "deployment/nginx",
    "namespace": "default",
    "duration":  30,
}
newExp, err := client.Experiments().Create("nginx-test", expConfig)
if err != nil {
    // Handle error
}

// Run an experiment
result, err := client.Experiments().Run("experiment-id")
if err != nil {
    // Handle error
}

// Stop a running experiment
err = client.Experiments().Stop("experiment-id")
if err != nil {
    // Handle error
}
```

### Working with Infrastructure

```go
// List all infrastructure resources
infra, err := client.Infrastructure().List()
if err != nil {
    // Handle error
}

// Create a new infrastructure resource
infraConfig := map[string]interface{}{
    "type":     "kubernetes",
    "provider": "gcp",
    "region":   "us-central1",
}
newInfra, err := client.Infrastructure().Create("gcp-cluster", infraConfig)
if err != nil {
    // Handle error
}

// Connect to infrastructure
connectParams := map[string]string{
    "kubeconfig": "/path/to/kubeconfig",
}
err = client.Infrastructure().Connect("infra-id", connectParams)
if err != nil {
    // Handle error
}
```

### Working with Probes

```go
// List all probes
probes, err := client.Probes().List()
if err != nil {
    // Handle error
}

// Create a new probe
probeConfig := map[string]interface{}{
    "type":    "http",
    "url":     "https://api.example.com/health",
    "method":  "GET",
    "timeout": 5,
}
newProbe, err := client.Probes().Create("health-check", probeConfig)
if err != nil {
    // Handle error
}

// Execute a probe
executeParams := map[string]string{
    "headers": "Content-Type: application/json",
}
result, err := client.Probes().Execute("probe-id", executeParams)
if err != nil {
    // Handle error
}
```

## Examples

For more examples, see the [examples](./examples) directory.

## License

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## Testing

The test suite is configured to use environment variables for test settings:

- `LITMUS_TEST_ENDPOINT`: The endpoint URL for the Litmus API (default: "http://127.0.0.1:39651")
- `LITMUS_TEST_USERNAME`: The username for authentication (default: "admin")
- `LITMUS_TEST_PASSWORD`: The password for authentication (default: "litmus")

### Running Tests

To run tests with default settings:

```bash
go test ./...
```

To run tests with custom settings:

```bash
LITMUS_TEST_ENDPOINT="https://your-litmus-instance.com" \
LITMUS_TEST_USERNAME="your-username" \
LITMUS_TEST_PASSWORD="your-password" \
go test ./...
```

The tests use a table-driven format which makes it easy to add new test cases.