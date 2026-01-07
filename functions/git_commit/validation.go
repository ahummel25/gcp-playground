package git_commit

import (
	"fmt"
	"strings"
)

// ValidateRepoFormat validates that the repository string is in the format "owner/repo"
func ValidateRepoFormat(repo string) error {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return fmt.Errorf("GITHUB_REPO must be in format 'owner/repo', got: %s", repo)
	}
	if parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("GITHUB_REPO owner and repo cannot be empty, got: %s", repo)
	}
	return nil
}
