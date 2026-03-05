package alter

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/wimpysworld/tailor/internal/config"
)

// ApplyMode controls whether changes are written to disk.
type ApplyMode int

const (
	DryRun     ApplyMode = iota // preview only
	Apply                       // write if file is absent or alteration permits
	ForceApply                  // overwrite unconditionally
)

// Run executes the alter command. It validates the config, applies
// repository settings, fetches the licence, and processes swatches.
func Run(cfg *config.Config, dir string, mode ApplyMode) error {
	if err := config.ValidateSources(cfg); err != nil {
		return err
	}
	if err := config.ValidateDuplicateDestinations(cfg); err != nil {
		return err
	}
	if err := config.ValidateRepoSettings(cfg); err != nil {
		return err
	}

	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("creating GitHub API client: %w", err)
	}

	// Phase 4 (repository settings) not yet implemented; pass nil.
	var repoResults []RepoSettingResult

	// Licence processing.
	licenceResult, err := ProcessLicence(cfg, dir, mode, client)
	if err != nil {
		return err
	}

	// Swatch processing with placeholder token context.
	tokens := &TokenContext{}
	swatchResults, err := ProcessSwatches(cfg, dir, mode, tokens)
	if err != nil {
		return err
	}

	// Merge licence result into swatch results for unified output.
	if licenceResult != nil {
		swatchResults = append(swatchResults, *licenceResult)
	}

	output := FormatOutput(repoResults, swatchResults)
	if output != "" {
		fmt.Print(output)
	}

	return nil
}
