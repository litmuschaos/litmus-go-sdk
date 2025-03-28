# Litmus Go SDK

A Go SDK for interacting with LitmusChaos.

## Testing

The test suite is configured to use environment variables for test settings:

- `LITMUS_TEST_ENDPOINT`: The endpoint URL for the Litmus API (default: "http://127.0.0.1:39651")
- `LITMUS_TEST_USERNAME`: The username for authentication (default: "admin")
- `LITMUS_TEST_PASSWORD`: The password for authentication (default: "LitmusChaos123@")

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