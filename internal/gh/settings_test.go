package gh

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cli/go-gh/v2/pkg/api"
)

// testTransport redirects all requests to the test server, preserving the
// original request path so the test handler can route by path.
type testTransport struct {
	server *httptest.Server
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = t.server.Listener.Addr().String()
	return http.DefaultTransport.RoundTrip(req)
}

// newTestClient creates an api.RESTClient that sends all requests to the
// given test server.
func newTestClient(t *testing.T, server *httptest.Server) *api.RESTClient {
	t.Helper()
	client, err := api.NewRESTClient(api.ClientOptions{
		Host:      "github.com",
		AuthToken: "test-token",
		Transport: &testTransport{server: server},
	})
	if err != nil {
		t.Fatalf("NewRESTClient: %v", err)
	}
	return client
}

const fullRepoJSON = `{
	"description": "A tailor for your repos",
	"homepage": "https://tailor.dev",
	"has_wiki": false,
	"has_discussions": true,
	"has_projects": false,
	"has_issues": true,
	"allow_merge_commit": false,
	"allow_squash_merge": true,
	"allow_rebase_merge": true,
	"squash_merge_commit_title": "PR_TITLE",
	"squash_merge_commit_message": "PR_BODY",
	"merge_commit_title": "PR_TITLE",
	"merge_commit_message": "PR_BODY",
	"delete_branch_on_merge": true,
	"allow_update_branch": true,
	"allow_auto_merge": true,
	"web_commit_signoff_required": false
}`

const pvrEnabledJSON = `{"enabled": true}`

func TestReadRepoSettings(t *testing.T) {
	tests := []struct {
		name     string
		repoJSON string
		pvrJSON  string
		// expected field checks
		wantDesc    *string
		wantHome    *string
		wantWiki    bool
		wantDisc    bool
		wantProj    bool
		wantIssues  bool
		wantMerge   bool
		wantSquash  bool
		wantRebase  bool
		wantSqTitle string
		wantSqMsg   string
		wantMcTitle string
		wantMcMsg   string
		wantDelete  bool
		wantUpdate  bool
		wantAuto    bool
		wantSignoff bool
		wantPVR     bool
	}{
		{
			name:        "all fields populated",
			repoJSON:    fullRepoJSON,
			pvrJSON:     pvrEnabledJSON,
			wantDesc:    strPtr("A tailor for your repos"),
			wantHome:    strPtr("https://tailor.dev"),
			wantWiki:    false,
			wantDisc:    true,
			wantProj:    false,
			wantIssues:  true,
			wantMerge:   false,
			wantSquash:  true,
			wantRebase:  true,
			wantSqTitle: "PR_TITLE",
			wantSqMsg:   "PR_BODY",
			wantMcTitle: "PR_TITLE",
			wantMcMsg:   "PR_BODY",
			wantDelete:  true,
			wantUpdate:  true,
			wantAuto:    true,
			wantSignoff: false,
			wantPVR:     true,
		},
		{
			name: "empty description and homepage produce nil",
			repoJSON: `{
				"description": "",
				"homepage": "",
				"has_wiki": true,
				"has_discussions": false,
				"has_projects": true,
				"has_issues": false,
				"allow_merge_commit": true,
				"allow_squash_merge": false,
				"allow_rebase_merge": false,
				"squash_merge_commit_title": "COMMIT_OR_PR_TITLE",
				"squash_merge_commit_message": "COMMIT_MESSAGES",
				"merge_commit_title": "MERGE_MESSAGE",
				"merge_commit_message": "PR_TITLE",
				"delete_branch_on_merge": false,
				"allow_update_branch": false,
				"allow_auto_merge": false,
				"web_commit_signoff_required": true
			}`,
			pvrJSON:     `{"enabled": false}`,
			wantDesc:    nil,
			wantHome:    nil,
			wantWiki:    true,
			wantDisc:    false,
			wantProj:    true,
			wantIssues:  false,
			wantMerge:   true,
			wantSquash:  false,
			wantRebase:  false,
			wantSqTitle: "COMMIT_OR_PR_TITLE",
			wantSqMsg:   "COMMIT_MESSAGES",
			wantMcTitle: "MERGE_MESSAGE",
			wantMcMsg:   "PR_TITLE",
			wantDelete:  false,
			wantUpdate:  false,
			wantAuto:    false,
			wantSignoff: true,
			wantPVR:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/repos/testowner/testrepo":
					fmt.Fprint(w, tt.repoJSON)
				case "/repos/testowner/testrepo/private-vulnerability-reporting":
					fmt.Fprint(w, tt.pvrJSON)
				default:
					http.NotFound(w, r)
				}
			}))
			t.Cleanup(server.Close)

			client := newTestClient(t, server)
			settings, err := ReadRepoSettings(client, "testowner", "testrepo")
			if err != nil {
				t.Fatalf("ReadRepoSettings() error: %v", err)
			}

			// description and homepage
			assertStringPtrEqual(t, "description", settings.Description, tt.wantDesc)
			assertStringPtrEqual(t, "homepage", settings.Homepage, tt.wantHome)

			// bool fields
			assertBoolPtrVal(t, "has_wiki", settings.HasWiki, tt.wantWiki)
			assertBoolPtrVal(t, "has_discussions", settings.HasDiscussions, tt.wantDisc)
			assertBoolPtrVal(t, "has_projects", settings.HasProjects, tt.wantProj)
			assertBoolPtrVal(t, "has_issues", settings.HasIssues, tt.wantIssues)
			assertBoolPtrVal(t, "allow_merge_commit", settings.AllowMergeCommit, tt.wantMerge)
			assertBoolPtrVal(t, "allow_squash_merge", settings.AllowSquashMerge, tt.wantSquash)
			assertBoolPtrVal(t, "allow_rebase_merge", settings.AllowRebaseMerge, tt.wantRebase)
			assertBoolPtrVal(t, "delete_branch_on_merge", settings.DeleteBranchOnMerge, tt.wantDelete)
			assertBoolPtrVal(t, "allow_update_branch", settings.AllowUpdateBranch, tt.wantUpdate)
			assertBoolPtrVal(t, "allow_auto_merge", settings.AllowAutoMerge, tt.wantAuto)
			assertBoolPtrVal(t, "web_commit_signoff_required", settings.WebCommitSignoffRequired, tt.wantSignoff)
			assertBoolPtrVal(t, "private_vulnerability_reporting_enabled", settings.PrivateVulnerabilityReportEnabled, tt.wantPVR)

			// string fields (always non-nil)
			assertStringPtrVal(t, "squash_merge_commit_title", settings.SquashMergeCommitTitle, tt.wantSqTitle)
			assertStringPtrVal(t, "squash_merge_commit_message", settings.SquashMergeCommitMessage, tt.wantSqMsg)
			assertStringPtrVal(t, "merge_commit_title", settings.MergeCommitTitle, tt.wantMcTitle)
			assertStringPtrVal(t, "merge_commit_message", settings.MergeCommitMessage, tt.wantMcMsg)
		})
	}
}

func TestReadRepoSettingsRepoAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message": "Not Found"}`)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, server)
	_, err := ReadRepoSettings(client, "testowner", "testrepo")
	if err == nil {
		t.Fatal("ReadRepoSettings() expected error, got nil")
	}
}

func TestReadRepoSettingsPVRAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/testowner/testrepo":
			fmt.Fprint(w, fullRepoJSON)
		default:
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, `{"message": "Forbidden"}`)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, server)
	_, err := ReadRepoSettings(client, "testowner", "testrepo")
	if err == nil {
		t.Fatal("ReadRepoSettings() expected error for PVR failure, got nil")
	}
}

func strPtr(s string) *string { return &s }

func assertBoolPtrVal(t *testing.T, name string, got *bool, want bool) {
	t.Helper()
	if got == nil {
		t.Errorf("%s is nil, want %v", name, want)
		return
	}
	if *got != want {
		t.Errorf("%s = %v, want %v", name, *got, want)
	}
}

func assertStringPtrVal(t *testing.T, name string, got *string, want string) {
	t.Helper()
	if got == nil {
		t.Errorf("%s is nil, want %q", name, want)
		return
	}
	if *got != want {
		t.Errorf("%s = %q, want %q", name, *got, want)
	}
}

func assertStringPtrEqual(t *testing.T, name string, got, want *string) {
	t.Helper()
	if want == nil {
		if got != nil {
			t.Errorf("%s = %q, want nil", name, *got)
		}
		return
	}
	if got == nil {
		t.Errorf("%s is nil, want %q", name, *want)
		return
	}
	if *got != *want {
		t.Errorf("%s = %q, want %q", name, *got, *want)
	}
}
