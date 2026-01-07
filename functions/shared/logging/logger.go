// Package logging provides a Cloud Logging client wrapper for GCP Cloud Functions.
package logging

import (
	"context"
	"fmt"
	"os"
	"sync"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
)

var (
	initOnce sync.Once
	logger   *logging.Logger
	client   *logging.Client
	initErr  error
)

// Init initializes Cloud Logging exactly once per process.
func Init(ctx context.Context, logName string) error {
	initOnce.Do(func() {
		projectID := os.Getenv("PROJECT_ID")
		if projectID == "" {
			pid, err := metadata.ProjectIDWithContext(ctx)
			if err != nil {
				initErr = fmt.Errorf("failed to get project ID: %w", err)
				return
			}
			projectID = pid
		}

		client, initErr = logging.NewClient(ctx, projectID)
		if initErr != nil {
			return
		}

		logger = client.Logger(logName)
	})

	return initErr
}

// Info logs an INFO level message.
func Info(format string, args ...any) {
	if logger == nil {
		// Should never happen in production
		fmt.Printf("[INFO] "+format+"\n", args...)
		return
	}

	logger.Log(logging.Entry{
		Severity: logging.Info,
		Payload:  fmt.Sprintf(format, args...),
	})
}

// Warning logs a WARNING level message.
func Warning(format string, args ...any) {
	if logger == nil {
		fmt.Printf("[WARNING] "+format+"\n", args...)
		return
	}

	logger.Log(logging.Entry{
		Severity: logging.Warning,
		Payload:  fmt.Sprintf(format, args...),
	})
}

// Error logs an ERROR level message.
func Error(format string, args ...any) {
	if logger == nil {
		fmt.Printf("[ERROR] "+format+"\n", args...)
		return
	}

	logger.Log(logging.Entry{
		Severity: logging.Error,
		Payload:  fmt.Sprintf(format, args...),
	})
}
