package gh

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
)

// licenceResponse holds the subset of fields returned by GET /licenses/{id}.
type licenceResponse struct {
	Body string `json:"body"`
}

// FetchLicence fetches licence text from the GitHub API.
// Returns the licence body text as returned by GET /licenses/{id}.
func FetchLicence(client *api.RESTClient, id string) (string, error) {
	var resp licenceResponse
	if err := client.Get(fmt.Sprintf("licenses/%s", id), &resp); err != nil {
		return "", fmt.Errorf("fetching licence %q: %w", id, err)
	}
	return resp.Body, nil
}
