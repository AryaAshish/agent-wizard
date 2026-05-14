package skills

import (
	"strings"
)

// IDContainsFold reports whether id contains needle as a case-insensitive substring.
// Empty or whitespace needle matches every id.
func IDContainsFold(id, needle string) bool {
	if strings.TrimSpace(needle) == "" {
		return true
	}
	n := strings.ToLower(strings.TrimSpace(needle))
	return strings.Contains(strings.ToLower(id), n)
}

// FilterSkillsByIDSubstring keeps skills whose ID contains needle (case-insensitive).
// Empty or whitespace needle returns skills unchanged.
func FilterSkillsByIDSubstring(skills []Skill, needle string) []Skill {
	if strings.TrimSpace(needle) == "" {
		return skills
	}
	var out []Skill
	for _, s := range skills {
		if IDContainsFold(s.ID, needle) {
			out = append(out, s)
		}
	}
	return out
}
