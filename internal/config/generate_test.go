package config

import (
	"io/fs"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/wimpysworld/tailor"
	"github.com/wimpysworld/tailor/internal/swatch"
)

func TestDefaultConfigMatchesEmbedded(t *testing.T) {
	// Parse the embedded config directly for comparison.
	data, err := fs.ReadFile(tailor.SwatchFS, embeddedConfigPath)
	if err != nil {
		t.Fatalf("reading embedded config: %v", err)
	}
	var want Config
	if err := yaml.Unmarshal(data, &want); err != nil {
		t.Fatalf("unmarshalling embedded config: %v", err)
	}

	got, err := DefaultConfig("MIT")
	if err != nil {
		t.Fatalf("DefaultConfig() error: %v", err)
	}

	// License should be the value we passed, not the embedded one.
	if got.License != "MIT" {
		t.Errorf("License = %q, want %q", got.License, "MIT")
	}

	// Repository settings should match the embedded config exactly.
	if got.Repository == nil {
		t.Fatal("Repository is nil, want non-nil")
	}
	assertBoolPtr(t, "has_wiki", got.Repository.HasWiki, false)
	assertBoolPtr(t, "has_discussions", got.Repository.HasDiscussions, false)
	assertBoolPtr(t, "has_projects", got.Repository.HasProjects, false)
	assertBoolPtr(t, "has_issues", got.Repository.HasIssues, true)
	assertBoolPtr(t, "allow_merge_commit", got.Repository.AllowMergeCommit, false)
	assertBoolPtr(t, "allow_squash_merge", got.Repository.AllowSquashMerge, true)
	assertBoolPtr(t, "allow_rebase_merge", got.Repository.AllowRebaseMerge, true)
	assertStringPtr(t, "squash_merge_commit_title", got.Repository.SquashMergeCommitTitle, "PR_TITLE")
	assertStringPtr(t, "squash_merge_commit_message", got.Repository.SquashMergeCommitMessage, "PR_BODY")
	assertBoolPtr(t, "delete_branch_on_merge", got.Repository.DeleteBranchOnMerge, true)
	assertBoolPtr(t, "allow_update_branch", got.Repository.AllowUpdateBranch, true)
	assertBoolPtr(t, "allow_auto_merge", got.Repository.AllowAutoMerge, true)
	assertBoolPtr(t, "web_commit_signoff_required", got.Repository.WebCommitSignoffRequired, false)
	assertBoolPtr(t, "private_vulnerability_reporting_enabled", got.Repository.PrivateVulnerabilityReportEnabled, true)

	// Fields absent from the embedded config should remain nil.
	if got.Repository.Description != nil {
		t.Errorf("Description = %q, want nil", *got.Repository.Description)
	}
	if got.Repository.Homepage != nil {
		t.Errorf("Homepage = %q, want nil", *got.Repository.Homepage)
	}
	if got.Repository.MergeCommitTitle != nil {
		t.Errorf("MergeCommitTitle = %q, want nil", *got.Repository.MergeCommitTitle)
	}
	if got.Repository.MergeCommitMessage != nil {
		t.Errorf("MergeCommitMessage = %q, want nil", *got.Repository.MergeCommitMessage)
	}

	// Swatch count and ordering must match exactly.
	if len(got.Swatches) != len(want.Swatches) {
		t.Fatalf("Swatches count = %d, want %d", len(got.Swatches), len(want.Swatches))
	}
	for i, g := range got.Swatches {
		w := want.Swatches[i]
		if g.Source != w.Source || g.Destination != w.Destination || g.Alteration != w.Alteration {
			t.Errorf("swatch[%d] = {%q, %q, %q}, want {%q, %q, %q}",
				i, g.Source, g.Destination, g.Alteration, w.Source, w.Destination, w.Alteration)
		}
	}
}

func TestDefaultConfigSwatchCount(t *testing.T) {
	cfg, err := DefaultConfig("MIT")
	if err != nil {
		t.Fatalf("DefaultConfig() error: %v", err)
	}
	if len(cfg.Swatches) != 16 {
		t.Errorf("Swatches count = %d, want 16", len(cfg.Swatches))
	}
}

func TestDefaultConfigSwatchOrder(t *testing.T) {
	cfg, err := DefaultConfig("MIT")
	if err != nil {
		t.Fatalf("DefaultConfig() error: %v", err)
	}

	first := cfg.Swatches[0]
	if first.Source != ".github/workflows/tailor.yml" {
		t.Errorf("first swatch Source = %q, want %q", first.Source, ".github/workflows/tailor.yml")
	}
	if first.Alteration != swatch.Always {
		t.Errorf("first swatch Alteration = %q, want %q", first.Alteration, swatch.Always)
	}

	last := cfg.Swatches[len(cfg.Swatches)-1]
	if last.Source != ".tailor/config.yml" {
		t.Errorf("last swatch Source = %q, want %q", last.Source, ".tailor/config.yml")
	}
	if last.Alteration != swatch.FirstFit {
		t.Errorf("last swatch Alteration = %q, want %q", last.Alteration, swatch.FirstFit)
	}
}

func stringPtr(v string) *string { return &v }

func TestMergeRepoSettings(t *testing.T) {
	tests := []struct {
		name        string
		live        *RepositorySettings
		description string
		wantDesc    *string // nil means expect nil
		wantHome    *string
	}{
		{
			name: "live settings override defaults entirely",
			live: &RepositorySettings{
				Description: stringPtr("live desc"),
				Homepage:    stringPtr("https://live.example.com"),
				HasWiki:     boolPtr(true),
				HasIssues:   boolPtr(false),
			},
			description: "",
			wantDesc:    stringPtr("live desc"),
			wantHome:    stringPtr("https://live.example.com"),
		},
		{
			name: "description flag overrides live description",
			live: &RepositorySettings{
				Description: stringPtr("live desc"),
				Homepage:    stringPtr("https://live.example.com"),
			},
			description: "flag desc",
			wantDesc:    stringPtr("flag desc"),
			wantHome:    stringPtr("https://live.example.com"),
		},
		{
			name: "empty description from live produces nil",
			live: &RepositorySettings{
				Description: stringPtr(""),
				Homepage:    stringPtr("https://live.example.com"),
			},
			description: "",
			wantDesc:    nil,
			wantHome:    stringPtr("https://live.example.com"),
		},
		{
			name: "empty homepage from live produces nil",
			live: &RepositorySettings{
				Description: stringPtr("live desc"),
				Homepage:    stringPtr(""),
			},
			description: "",
			wantDesc:    stringPtr("live desc"),
			wantHome:    nil,
		},
		{
			name: "non-empty description flag with empty live description sets flag value",
			live: &RepositorySettings{
				Description: stringPtr(""),
				Homepage:    stringPtr("https://live.example.com"),
			},
			description: "flag desc",
			wantDesc:    stringPtr("flag desc"),
			wantHome:    stringPtr("https://live.example.com"),
		},
		{
			name: "empty description flag with non-empty live description preserves live value",
			live: &RepositorySettings{
				Description: stringPtr("live desc"),
				Homepage:    stringPtr("https://live.example.com"),
			},
			description: "",
			wantDesc:    stringPtr("live desc"),
			wantHome:    stringPtr("https://live.example.com"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				License: "MIT",
				Repository: &RepositorySettings{
					HasWiki:   boolPtr(false),
					HasIssues: boolPtr(true),
				},
			}

			MergeRepoSettings(cfg, tt.live, tt.description)

			// Repository must point to the live object.
			if cfg.Repository != tt.live {
				t.Fatal("Repository was not replaced with live settings")
			}

			// Check Description.
			if tt.wantDesc == nil {
				if cfg.Repository.Description != nil {
					t.Errorf("Description = %q, want nil", *cfg.Repository.Description)
				}
			} else {
				assertStringPtr(t, "description", cfg.Repository.Description, *tt.wantDesc)
			}

			// Check Homepage.
			if tt.wantHome == nil {
				if cfg.Repository.Homepage != nil {
					t.Errorf("Homepage = %q, want nil", *cfg.Repository.Homepage)
				}
			} else {
				assertStringPtr(t, "homepage", cfg.Repository.Homepage, *tt.wantHome)
			}
		})
	}
}

func TestMergeRepoSettingsPreservesMergeCommitFields(t *testing.T) {
	mergeTitle := "PR_TITLE"
	mergeMessage := "PR_BODY"
	live := &RepositorySettings{
		Description:        stringPtr("desc"),
		AllowMergeCommit:   boolPtr(false),
		MergeCommitTitle:   &mergeTitle,
		MergeCommitMessage: &mergeMessage,
	}

	cfg := &Config{License: "MIT"}
	MergeRepoSettings(cfg, live, "")

	assertStringPtr(t, "merge_commit_title", cfg.Repository.MergeCommitTitle, "PR_TITLE")
	assertStringPtr(t, "merge_commit_message", cfg.Repository.MergeCommitMessage, "PR_BODY")
}

func TestDefaultConfigLicenseValues(t *testing.T) {
	tests := []struct {
		name    string
		license string
	}{
		{name: "MIT", license: "MIT"},
		{name: "Apache-2.0", license: "Apache-2.0"},
		{name: "none", license: "none"},
		{name: "empty", license: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := DefaultConfig(tt.license)
			if err != nil {
				t.Fatalf("DefaultConfig(%q) error: %v", tt.license, err)
			}
			if cfg.License != tt.license {
				t.Errorf("License = %q, want %q", cfg.License, tt.license)
			}
		})
	}
}
