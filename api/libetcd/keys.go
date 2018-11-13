package libetcd

import (
	"fmt"
)

// KeyGitHubAuthToken is the key used to store the GitHub Auth token
const KeyGitHubAuthToken string = "/github/auth/token"

// KeyTrackedGitHubRepos is the key used to store tracked GitHub repositories
const KeyTrackedGitHubRepos string = "/github/repositories/tracked"

// GetTrackedGitHubRepoKey formats a key for a tracked GitHub repo
func GetTrackedGitHubRepoKey(user, repo string) string {
	return fmt.Sprintf("%s/%s/%s", KeyTrackedGitHubRepos, user, repo)
}
