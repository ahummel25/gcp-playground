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
)

func init() {
	functions.HTTP("EmptyCommit", EmptyCommit)
}

// EmptyCommit is the HTTP Cloud Function handler that creates a verified empty commit
func EmptyCommit(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	repo := os.Getenv("GITHUB_REPO")
	if repo == "" {
		log.Printf("ERROR: GITHUB_REPO not set")
		http.Error(w, "GITHUB_REPO not set", http.StatusInternalServerError)
		return
	}

	if err := ValidateRepoFormat(repo); err != nil {
		log.Printf("ERROR: Invalid GITHUB_REPO format: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Starting empty commit for repo: %s", repo)

	token, err := GetGitHubToken(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to get GitHub token: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}

	parts := strings.Split(repo, "/")
	owner := parts[0]
	repoName := parts[1]
	repoNameWithOwner := repo

	log.Printf("INFO: Getting repository ID and HEAD SHA for repo: %s", repo)
	repoInfo, err := GetRepositoryInfo(client, token, owner, repoName)
	if err != nil {
		log.Printf("ERROR: Failed to get repository info: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("INFO: Repository ID: %s, Current HEAD SHA: %s", repoInfo.ID, repoInfo.HeadSHA)

	log.Printf("INFO: Creating verified empty commit")
	commitSHA, err := CreateVerifiedCommit(client, token, repoNameWithOwner, repoInfo.HeadSHA)
	if err != nil {
		log.Printf("ERROR: Failed to create verified commit: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("INFO: Created verified commit SHA: %s", commitSHA)

	log.Printf("INFO: Successfully created verified empty commit: %s", commitSHA)
	fmt.Fprintf(w, "Created verified empty commit: %s\n", commitSHA)
}
