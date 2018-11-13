package libetcd

import (
	"fmt"
)

// KeyGitHubAuthToken is the key used to store the GitHub Auth token
const KeyGitHubAuthToken string = "/github/auth/token"

// KeyTrackedGitHubRepos is the key used to store tracked GitHub repositories
const KeyTrackedGitHubRepos string = "/github/repositories/tracked"

// KeyTrackedGitHubReposWebHookIDs is the key of the directory used to store
// IDs of GitHub web hooks
const KeyTrackedGitHubReposWebHookIDs string = "/github/repositories/web_hooks/tracked"

// GetTrackedGitHubRepoKey formats a key for a tracked GitHub repo
func GetTrackedGitHubRepoKey(user, repo string) string {
	return fmt.Sprintf("%s/%s/%s", KeyTrackedGitHubRepos, user, repo)
}

// GetTrackedGitRepoWebHookIDKey formats a key for a tracked GitHub repo's
// web hook ID
func GetTrackedGitRepoWebHookIDKey(user, repo string) string {
	return fmt.Sprintf("%s/%s/%s", KeyTrackedGitHubReposWebHookIDs, user,
		repo)
}
