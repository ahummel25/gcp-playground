# Shared Logging Package

This package provides a Cloud Logging client wrapper for GCP Cloud Functions. It offers structured logging with severity levels and automatic resource detection.

## Features

- **Structured Logging**: Logs with proper severity levels (INFO, ERROR, WARNING, DEBUG)
- **Thread-Safe**: Uses mutex locks to ensure safe concurrent access
- **Automatic Resource Detection**: Cloud Logging automatically detects Cloud Function context
- **Fallback Support**: Falls back to standard output if Cloud Logging initialization fails
- **Reusable**: Can be used across multiple Cloud Functions

## Usage

### Basic Example

```go
package myfunction

import (
    "context"
    "net/http"
    
    "github.com/GoogleCloudPlatform/functions-framework-go/functions"
    "github.com/hummelgcp/go/shared/logging"
)

func init() {
    functions.HTTP("MyFunction", MyFunction)
}

func MyFunction(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    
    // Initialize Cloud Logging
    if err := logging.Init(ctx, "my-function"); err != nil {
        // Fallback handled automatically
    }
    defer logging.Close()
    
    // Use logging functions
    logging.Info(ctx, "Function started")
    logging.Error(ctx, "An error occurred: %v", err)
    logging.Warning(ctx, "This is a warning")
    logging.Debug(ctx, "Debug information: %s", data)
}
```

### Adding to Your Function's go.mod

Add the shared logging package to your function's `go.mod`:

```go
module github.com/hummelgcp/go

require (
    github.com/hummelgcp/go/shared/logging v0.0.0
)

replace github.com/hummelgcp/go/shared/logging => ../shared/logging
```

Then run `go mod tidy` to update dependencies.

## API Reference

### `Init(ctx context.Context, logName string) error`

Initializes the Cloud Logging client and logger. Should be called once at the start of a Cloud Function.

- `ctx`: Context for the operation
- `logName`: Name of the log stream in Cloud Logging (e.g., "my-function")

### `Close() error`

Closes the logging client and flushes any pending log entries. Should be called at the end of a Cloud Function execution (typically with `defer`).

### `Info(ctx context.Context, format string, args ...interface{})`

Logs an INFO level message. Use for general informational messages.

### `Error(ctx context.Context, format string, args ...interface{})`

Logs an ERROR level message. Use for error conditions.

### `Warning(ctx context.Context, format string, args ...interface{})`

Logs a WARNING level message. Use for warning conditions.

### `Debug(ctx context.Context, format string, args ...interface{})`

Logs a DEBUG level message. Use for detailed debugging information.

## Environment Variables

- `PROJECT_ID`: Required. The GCP project ID where logs will be written.

## Best Practices

1. **Always initialize at function start**: Call `logging.Init()` at the beginning of your function handler
2. **Always close at function end**: Use `defer logging.Close()` to ensure logs are flushed
3. **Use appropriate severity levels**: Choose the right level for your message
4. **Include context**: Pass the request context to logging functions
5. **Log meaningful information**: Include relevant details like IDs, timestamps, or error messages

