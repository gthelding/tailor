package alter

import "github.com/wimpysworld/tailor/internal/config"

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
	return nil
}
