package engine

import (
	"strings"

	"github.com/aryaashish/agent-wizard/internal/community"
	"github.com/aryaashish/agent-wizard/internal/config"
)

const legacyCommunityGitCanonical = "https://github.com/aryaashish/agent-skills-community.git"

func normalizeGitURLForCompare(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.TrimSuffix(s, "/")
	return strings.ToLower(s)
}

// IsLegacyCommunityGitSource reports whether src is deprecated global config
// for the removed agent-skills-community repository; such entries should use
// the embedded starter library instead.
func IsLegacyCommunityGitSource(src config.Source) bool {
	if src.Name != community.SourceName || src.Kind != "git" {
		return false
	}
	return normalizeGitURLForCompare(src.GitURL) == legacyCommunityGitCanonical
}
