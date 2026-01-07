package git_commit

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"

	"github.com/hummelgcp/go/shared/logging"
)

func init() {
	ctx := context.Background()

	if err := logging.Init(ctx, "git-commit-function"); err != nil {
		log.Fatalf("failed to init logging: %v", err)
	}

	functions.HTTP("EmptyCommit", EmptyCommit)
}

// EmptyCommit is the HTTP Cloud Function handler that creates a verified empty commit
func EmptyCommit(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	repo := os.Getenv("GITHUB_REPO")
	if repo == "" {
		logging.Error("GITHUB_REPO not set")
		http.Error(w, "GITHUB_REPO not set", http.StatusInternalServerError)
		return
	}

	if err := ValidateRepoFormat(repo); err != nil {
		logging.Error("Invalid GITHUB_REPO format: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logging.Info("Starting empty commit for repo: %s", repo)

	token, err := GetGitHubToken(ctx)
	if err != nil {
		logging.Error("Failed to get GitHub token: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}

	parts := strings.Split(repo, "/")
	owner := parts[0]
	repoName := parts[1]
	repoNameWithOwner := repo

	logging.Info("Getting repository ID and HEAD SHA for repo: %s", repo)
	repoInfo, err := GetRepositoryInfo(client, token, owner, repoName)
	if err != nil {
		logging.Error("Failed to get repository info: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logging.Info("Repository ID: %s, Current HEAD SHA: %s", repoInfo.ID, repoInfo.HeadSHA)

	logging.Info("Creating verified empty commit")
	commitSHA, err := CreateVerifiedCommit(client, token, repoNameWithOwner, repoInfo.HeadSHA)
	if err != nil {
		logging.Error("Failed to create verified commit: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logging.Info("Created verified commit SHA: %s", commitSHA)

	logging.Info("Successfully created verified empty commit: %s", commitSHA)
	fmt.Fprintf(w, "Created verified empty commit: %s\n", commitSHA)
}
