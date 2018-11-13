package libetcd

import (
	"fmt"
)

// KeyGitHubAuthToken is the key used to store the GitHub Auth token
const KeyGitHubAuthToken string = "/github/auth/token"

// KeyDirTrackedGHRepos is the key used to store tracked GitHub repositories
const KeyDirTrackedGHRepos string = "/github/repositories/tracked"

// KeyTrackedGHRepoName is the key placed inside of a GitHub repo directory
// which stores its name
const KeyTrackedGHRepoName string = "name"

// KeyTrackedGHRepoWebHookID is the key placed inside of a GitHub repo
// directory which stores its web hook ID
const KeyTrackedGHRepoWebHookID string = "web_hook_id"

// GetTrackedGHRepoDirKey returns the directory key for a tracked
// GitHub repository
func GetTrackedGHRepoDirKey(user, repo string) string {
	return fmt.Sprintf("%s/%s/%s", KeyDirTrackedGHRepos, user, repo)
}

// GetTrackedGHRepoNameKey returns a key to a node which holds a repo's name
func GetTrackedGHRepoNameKey(user, repo string) string {
	return fmt.Sprintf("%s/%s/%s/%s", KeyDirTrackedGHRepos, user, repo,
		KeyTrackedGHRepoName)
}

// GetTrackedGHRepoWebHookIDKey returns a key to a node which holds a repo's
// web hook ID
func GetTrackedGHRepoWebHookIDKey(user, repo string) string {
	return fmt.Sprintf("%s/%s/%s/%s", KeyDirTrackedGHRepos, user, repo,
		KeyTrackedGHRepoWebHookID)
}
