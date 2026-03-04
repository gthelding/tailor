package gh

import "testing"

func TestCheckAuth(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr string
	}{
		{
			name:  "valid token returns nil",
			token: "gho_validtoken123",
		},
		{
			name:    "empty token returns error",
			token:   "",
			wantErr: "tailor requires an authenticated GitHub CLI. Run 'gh auth login' to authenticate.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tokenForHost
			tokenForHost = func(host string) (string, string) {
				return tt.token, "oauth_token"
			}
			t.Cleanup(func() { tokenForHost = original })

			err := CheckAuth()

			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("CheckAuth() = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("CheckAuth() = nil, want error %q", tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Errorf("CheckAuth() error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}
