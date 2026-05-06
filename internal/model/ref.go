package model

import (
	"fmt"
	"strings"
)

// SkillRef parses optional source qualification: "mysource/pr-review" vs "pr-review".
type SkillRef struct {
	SourceAlias string // empty if bare id
	ID          string
}

func ParseSkillRef(s string) (SkillRef, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return SkillRef{}, fmt.Errorf("empty skill reference")
	}
	if i := strings.IndexByte(s, '/'); i >= 0 {
		alias := strings.TrimSpace(s[:i])
		id := strings.TrimSpace(s[i+1:])
		if alias == "" || id == "" {
			return SkillRef{}, fmt.Errorf("invalid qualified skill reference %q", s)
		}
		return SkillRef{SourceAlias: alias, ID: id}, nil
	}
	return SkillRef{SourceAlias: "", ID: s}, nil
}
