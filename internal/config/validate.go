package config

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/wimpysworld/tailor/internal/swatch"
)

// ValidateSources checks that every swatch source in cfg matches a known
// embedded swatch. Returns an error listing the unrecognised source and all
// valid source names.
func ValidateSources(cfg *Config) error {
	valid := swatch.SourceNames()
	known := make(map[string]bool, len(valid))
	for _, name := range valid {
		known[name] = true
	}
	for _, s := range cfg.Swatches {
		if !known[s.Source] {
			return fmt.Errorf("unrecognised swatch source %q in config; valid sources: %s",
				s.Source, strings.Join(valid, ", "))
		}
	}
	return nil
}

// ValidateDuplicateDestinations checks that no two swatches share a
// destination. Returns an error identifying the conflicting entries.
func ValidateDuplicateDestinations(cfg *Config) error {
	seen := make(map[string]string, len(cfg.Swatches))
	for _, s := range cfg.Swatches {
		if prev, ok := seen[s.Destination]; ok {
			return fmt.Errorf("duplicate destination %q in config: sources %q and %q both target the same file",
				s.Destination, prev, s.Source)
		}
		seen[s.Destination] = s.Source
	}
	return nil
}

// ValidateRepoSettings checks that every field name in cfg.Repository
// matches the supported settings list. Returns an error identifying the
// unrecognised field and listing all valid field names.
func ValidateRepoSettings(cfg *Config) error {
	if cfg.Repository == nil {
		return nil
	}

	if len(cfg.Repository.Extra) > 0 {
		valid := repoSettingNames()
		for key := range cfg.Repository.Extra {
			return fmt.Errorf("unrecognised repository setting %q in config; valid settings: %s",
				key, strings.Join(valid, ", "))
		}
	}
	return nil
}

// repoSettingNames returns the sorted list of recognised yaml tag names from
// RepositorySettings, excluding the inline Extra field.
func repoSettingNames() []string {
	t := reflect.TypeOf(RepositorySettings{})
	var names []string
	for i := range t.NumField() {
		tag := t.Field(i).Tag.Get("yaml")
		if tag == "" || tag == ",inline" {
			continue
		}
		name, _, _ := strings.Cut(tag, ",")
		if name != "" {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}
