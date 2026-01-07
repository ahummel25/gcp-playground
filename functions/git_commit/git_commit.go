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

	mode := os.Getenv("COMMIT_MODE")
	if mode == "" {
		mode = "scheduled"
	}

	var bucket string
	var today string
	commitBatch := 1

	if mode == "random" {
		bucket = os.Getenv("GIT_COMMIT_STATE_BUCKET")
		today = time.Now().UTC().Format(time.DateOnly)
		doCommit, err := DecideAndRecord(ctx, bucket, today)
		if err != nil {
			logging.Error("DecideAndRecord: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if s, err := LoadState(ctx, bucket); err == nil {
			p, hardRule := CommitDecisionProbability(s, today)
			if hardRule != "" {
				logging.Info("Commit decision for %s: %s (probability 0%%)", today, hardRule)
			} else {
				logging.Info("Commit decision for %s: %t (probability %.1f%%)", today, doCommit, p)
			}
		}
		if !doCommit {
			logging.Info("Skipping commit (random mode): state-based decision for %s", today)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Skipped commit (random mode)\n")
			return
		}
		// Determine how many commits to make this run (heavy days can do 2–5)
		if s, err := LoadState(ctx, bucket); err == nil {
			commitBatch = CommitBatchSize(s, today)
		}
		if commitBatch <= 0 {
			logging.Info("Skipping commit (random mode): daily cap reached for %s", today)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Skipped commit (random mode)\n")
			return
		}
		// doCommit: fall through to perform GitHub commit; we will call RecordCommits after success
	}

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

	logging.Info("Creating verified empty commit(s)")
	commitSHAs := make([]string, 0, commitBatch)
	expectedHead := repoInfo.HeadSHA
	successCount := 0
	for i := 0; i < commitBatch; i++ {
		commitSHA, err := CreateVerifiedCommit(client, token, repoNameWithOwner, expectedHead)
		if err != nil {
			logging.Error("Failed to create verified commit: %v", err)
			if mode == "random" && successCount > 0 {
				if err := RecordCommits(ctx, bucket, today, successCount); err != nil {
					logging.Error("RecordCommits after partial success: %v", err)
				}
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		commitSHAs = append(commitSHAs, commitSHA)
		expectedHead = commitSHA
		successCount++
	}

	logging.Info("Created verified commit SHAs: %s", strings.Join(commitSHAs, ", "))

	if mode == "random" {
		if err := RecordCommits(ctx, bucket, today, successCount); err != nil {
			logging.Error("RecordCommits after successful commits: %v", err)
		}
	}

	fmt.Fprintf(w, "Created verified empty commit(s): %s\n", strings.Join(commitSHAs, ", "))
}
