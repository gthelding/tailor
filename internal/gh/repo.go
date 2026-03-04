package gh

import (
	"github.com/cli/go-gh/v2/pkg/repository"
)

// currentRepo wraps repository.Current for testability.
var currentRepo = func() (repository.Repository, error) {
	return repository.Current()
}

// RepoContext detects the GitHub repository for the current directory.
// It returns the owner and name if a GitHub remote is found.
// When no remote is configured, it returns ok=false with a nil error.
func RepoContext() (owner string, name string, ok bool, err error) {
	repo, repoErr := currentRepo()
	if repoErr != nil {
		return "", "", false, nil
	}
	return repo.Owner, repo.Name, true, nil
}
