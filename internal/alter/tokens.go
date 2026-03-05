package alter

// TokenContext holds resolved values for template substitution.
// Placeholder until Phase 3 implements full token resolution.
type TokenContext struct{}

// Substitute is a no-op placeholder. Returns content unchanged.
func (tc *TokenContext) Substitute(content []byte, source string) []byte {
	return content
}
