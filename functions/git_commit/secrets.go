package git_commit

import (
	"context"
	"fmt"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// GetGitHubToken retrieves the GitHub token from Secret Manager
func GetGitHubToken(ctx context.Context) (string, error) {
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		return "", fmt.Errorf("PROJECT_ID not set")
	}

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secret manager client: %w", err)
	}
	defer client.Close()

	secretName := fmt.Sprintf("projects/%s/secrets/github_token/versions/latest", projectID)

	resp, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	})
	if err != nil {
		return "", fmt.Errorf("failed to access secret: %w", err)
	}

	return string(resp.Payload.Data), nil
}
