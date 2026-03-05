package alter

import (
	"fmt"
	"slices"
	"strings"
)

// RepoSettingCategory classifies the outcome of processing a single repository setting.
type RepoSettingCategory string

const (
	WouldSet   RepoSettingCategory = "would set"
	RSNoChange RepoSettingCategory = "no change"
)

// RepoSettingResult records the field name, category, and display value for one
// repository setting.
type RepoSettingResult struct {
	Field    string
	Category RepoSettingCategory
	Value    string
}

// labelWidth is the fixed column width for status labels in formatted output.
// Sized to accommodate "skipped (first-fit, exists): " (29 characters).
const labelWidth = 29

// FormatOutput produces the alter command output from repo settings results
// and swatch results (including licence).
func FormatOutput(repoResults []RepoSettingResult, swatchResults []SwatchResult) string {
	if len(repoResults) == 0 && len(swatchResults) == 0 {
		return ""
	}

	var b strings.Builder

	for _, r := range sortRepoResults(repoResults) {
		label := string(r.Category) + ":"
		switch r.Category {
		case WouldSet:
			fmt.Fprintf(&b, "%-*srepository.%s = %s\n", labelWidth, label, r.Field, r.Value)
		case RSNoChange:
			fmt.Fprintf(&b, "%-*srepository.%s (already %s)\n", labelWidth, label, r.Field, r.Value)
		}
	}

	for _, r := range sortSwatchResults(swatchResults) {
		label := string(r.Category) + ":"
		fmt.Fprintf(&b, "%-*s%s\n", labelWidth, label, r.Destination)
	}

	return b.String()
}

// sortRepoResults returns a sorted copy: actionable (WouldSet) before
// informational (RSNoChange), lexicographic by field within each group.
func sortRepoResults(results []RepoSettingResult) []RepoSettingResult {
	if len(results) == 0 {
		return nil
	}
	sorted := make([]RepoSettingResult, len(results))
	copy(sorted, results)
	slices.SortStableFunc(sorted, func(a, b RepoSettingResult) int {
		ao, bo := repoOrder(a.Category), repoOrder(b.Category)
		if ao != bo {
			if ao < bo {
				return -1
			}
			return 1
		}
		return strings.Compare(a.Field, b.Field)
	})
	return sorted
}

// repoOrder returns the sort priority for a RepoSettingCategory.
func repoOrder(c RepoSettingCategory) int {
	switch c {
	case WouldSet:
		return 0
	default:
		return 1
	}
}

// sortSwatchResults returns a sorted copy: actionable (WouldCopy, WouldOverwrite)
// before informational (NoChange, Skipped), lexicographic by destination within
// each group.
func sortSwatchResults(results []SwatchResult) []SwatchResult {
	if len(results) == 0 {
		return nil
	}
	sorted := make([]SwatchResult, len(results))
	copy(sorted, results)
	slices.SortStableFunc(sorted, func(a, b SwatchResult) int {
		ao, bo := swatchOrder(a.Category), swatchOrder(b.Category)
		if ao != bo {
			if ao < bo {
				return -1
			}
			return 1
		}
		return strings.Compare(a.Destination, b.Destination)
	})
	return sorted
}

// swatchOrder returns the sort priority for a SwatchCategory.
func swatchOrder(c SwatchCategory) int {
	switch c {
	case WouldCopy:
		return 0
	case WouldOverwrite:
		return 1
	case NoChange:
		return 2
	default:
		return 3
	}
}
