package gh

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/wimpysworld/tailor/internal/config"
)

// repoResponse holds the subset of GitHub repository fields we read.
type repoResponse struct {
	Description              string `json:"description"`
	Homepage                 string `json:"homepage"`
	HasWiki                  bool   `json:"has_wiki"`
	HasDiscussions           bool   `json:"has_discussions"`
	HasProjects              bool   `json:"has_projects"`
	HasIssues                bool   `json:"has_issues"`
	AllowMergeCommit         bool   `json:"allow_merge_commit"`
	AllowSquashMerge         bool   `json:"allow_squash_merge"`
	AllowRebaseMerge         bool   `json:"allow_rebase_merge"`
	SquashMergeCommitTitle   string `json:"squash_merge_commit_title"`
	SquashMergeCommitMessage string `json:"squash_merge_commit_message"`
	MergeCommitTitle         string `json:"merge_commit_title"`
	MergeCommitMessage       string `json:"merge_commit_message"`
	DeleteBranchOnMerge      bool   `json:"delete_branch_on_merge"`
	AllowUpdateBranch        bool   `json:"allow_update_branch"`
	AllowAutoMerge           bool   `json:"allow_auto_merge"`
	WebCommitSignoffRequired bool   `json:"web_commit_signoff_required"`
}

// pvrResponse holds the private vulnerability reporting status.
type pvrResponse struct {
	Enabled bool `json:"enabled"`
}

// ReadRepoSettings fetches repository settings from the GitHub API and returns
// them as a config.RepositorySettings. It makes two API calls: one for the
// standard repository fields and one for private vulnerability reporting.
func ReadRepoSettings(client *api.RESTClient, owner, name string) (*config.RepositorySettings, error) {
	var repo repoResponse
	if err := client.Get(fmt.Sprintf("repos/%s/%s", owner, name), &repo); err != nil {
		return nil, fmt.Errorf("fetching repo settings: %w", err)
	}

	var pvr pvrResponse
	if err := client.Get(fmt.Sprintf("repos/%s/%s/private-vulnerability-reporting", owner, name), &pvr); err != nil {
		return nil, fmt.Errorf("fetching private vulnerability reporting: %w", err)
	}

	s := &config.RepositorySettings{
		HasWiki:                          boolPtr(repo.HasWiki),
		HasDiscussions:                   boolPtr(repo.HasDiscussions),
		HasProjects:                      boolPtr(repo.HasProjects),
		HasIssues:                        boolPtr(repo.HasIssues),
		AllowMergeCommit:                 boolPtr(repo.AllowMergeCommit),
		AllowSquashMerge:                 boolPtr(repo.AllowSquashMerge),
		AllowRebaseMerge:                 boolPtr(repo.AllowRebaseMerge),
		SquashMergeCommitTitle:           stringPtr(repo.SquashMergeCommitTitle),
		SquashMergeCommitMessage:         stringPtr(repo.SquashMergeCommitMessage),
		MergeCommitTitle:                 stringPtr(repo.MergeCommitTitle),
		MergeCommitMessage:               stringPtr(repo.MergeCommitMessage),
		DeleteBranchOnMerge:              boolPtr(repo.DeleteBranchOnMerge),
		AllowUpdateBranch:                boolPtr(repo.AllowUpdateBranch),
		AllowAutoMerge:                   boolPtr(repo.AllowAutoMerge),
		WebCommitSignoffRequired:         boolPtr(repo.WebCommitSignoffRequired),
		PrivateVulnerabilityReportEnabled: boolPtr(pvr.Enabled),
	}

	// Set description and homepage to nil when the API returns empty strings
	// so that omitempty drops them from YAML output.
	if repo.Description != "" {
		s.Description = stringPtr(repo.Description)
	}
	if repo.Homepage != "" {
		s.Homepage = stringPtr(repo.Homepage)
	}

	return s, nil
}

func boolPtr(v bool) *bool       { return &v }
func stringPtr(v string) *string { return &v }
