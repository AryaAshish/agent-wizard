package skills

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

const skillSummaryMaxRunes = 80

// SkillSummaryLine reads SKILL.md next to the skill folder and returns a short summary,
// or "-" if missing or unreadable.
func SkillSummaryLine(s Skill) string {
	p := filepath.Join(s.Path, "SKILL.md")
	data, err := os.ReadFile(p)
	if err != nil {
		return "-"
	}
	return summarizeSkillMarkdown(string(data))
}

// summarizeSkillMarkdown returns the first paragraph under the first H1, truncated,
// or "-" if none (plain Markdown only — no YAML frontmatter).
func summarizeSkillMarkdown(md string) string {
	s := strings.TrimPrefix(md, "\ufeff")
	var lines []string
	for _, part := range strings.Split(s, "\n") {
		lines = append(lines, strings.TrimSuffix(part, "\r"))
	}
	h1 := -1
	for i, ln := range lines {
		t := strings.TrimSpace(ln)
		if strings.HasPrefix(t, "# ") && len(t) > 2 {
			h1 = i
			break
		}
	}
	if h1 < 0 {
		return "-"
	}
	i := h1 + 1
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}
	var para []string
	for ; i < len(lines); i++ {
		t := strings.TrimSpace(lines[i])
		if t == "" {
			break
		}
		if strings.HasPrefix(t, "#") {
			break
		}
		para = append(para, t)
	}
	if len(para) == 0 {
		return "-"
	}
	joined := collapseSpace(strings.Join(para, " "))
	return truncateRunes(joined, skillSummaryMaxRunes)
}

func collapseSpace(s string) string {
	var b strings.Builder
	var prevSpace bool
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !prevSpace {
				b.WriteRune(' ')
				prevSpace = true
			}
			continue
		}
		prevSpace = false
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

func truncateRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	var b strings.Builder
	n := 0
	for _, r := range s {
		if n >= max {
			break
		}
		b.WriteRune(r)
		n++
	}
	return b.String()
}
