package ci

import "os"

// IsRunningOnGithubActions returns true if the current process is running on Github Actions.
func IsRunningOnGithubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}
