package git_commit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// RepositoryInfo contains repository ID and current HEAD SHA
type RepositoryInfo struct {
	ID      string
	HeadSHA string
}

// GetRepositoryInfo retrieves repository ID and current HEAD SHA using GraphQL
func GetRepositoryInfo(client *http.Client, token, owner, repoName string) (RepositoryInfo, error) {
	query := `
		query($owner: String!, $repo: String!) {
			repository(owner: $owner, name: $repo) {
				id
				defaultBranchRef {
					target {
						... on Commit {
							oid
						}
					}
				}
			}
		}
	`

	reqBody := GraphQLRequest{
		Query: query,
		Variables: map[string]interface{}{
			"owner": owner,
			"repo":  repoName,
		},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return RepositoryInfo{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(b))
	if err != nil {
		return RepositoryInfo{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return RepositoryInfo{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RepositoryInfo{}, fmt.Errorf("GitHub GraphQL API returned status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Repository RepositoryData `json:"repository"`
		} `json:"data"`
		Errors []GraphQLError `json:"errors"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return RepositoryInfo{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Errors) > 0 {
		return RepositoryInfo{}, fmt.Errorf("GraphQL errors: %v", result.Errors)
	}

	repoID := result.Data.Repository.ID
	headSHA := result.Data.Repository.DefaultBranchRef.Target.OID

	if repoID == "" {
		return RepositoryInfo{}, fmt.Errorf("empty repository ID in response")
	}

	if headSHA == "" {
		return RepositoryInfo{}, fmt.Errorf("empty HEAD SHA in response")
	}

	return RepositoryInfo{
		ID:      repoID,
		HeadSHA: headSHA,
	}, nil
}

// CreateVerifiedCommit creates a verified empty commit using GraphQL createCommitOnBranch mutation
func CreateVerifiedCommit(client *http.Client, token, repoNameWithOwner, expectedHeadOID string) (string, error) {
	message := fmt.Sprintf("chore: scheduled empty commit (%s)", time.Now().UTC().Format(time.RFC3339))

	query := `
		mutation($input: CreateCommitOnBranchInput!) {
			createCommitOnBranch(input: $input) {
				commit {
					oid
				}
			}
		}
	`

	branchInput := map[string]interface{}{
		"repositoryNameWithOwner": repoNameWithOwner,
		"branchName":              "main",
	}

	input := map[string]interface{}{
		"branch":          branchInput,
		"expectedHeadOid": expectedHeadOID,
		"message": map[string]interface{}{
			"headline": message,
		},
		"fileChanges": map[string]interface{}{},
	}

	reqBody := GraphQLRequest{
		Query: query,
		Variables: map[string]interface{}{
			"input": input,
		},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes := make([]byte, 1024)
		n, _ := resp.Body.Read(bodyBytes)
		return "", fmt.Errorf("GitHub GraphQL API returned status %d: %s", resp.StatusCode, string(bodyBytes[:n]))
	}

	var result struct {
		Data struct {
			CreateCommitOnBranch CommitData `json:"createCommitOnBranch"`
		} `json:"data"`
		Errors []GraphQLError `json:"errors"`
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Errors) > 0 {
		return "", fmt.Errorf("GraphQL errors: %v", result.Errors)
	}

	commitOID := result.Data.CreateCommitOnBranch.Commit.OID
	if commitOID == "" {
		return "", fmt.Errorf("empty commit SHA in response")
	}

	return commitOID, nil
}
