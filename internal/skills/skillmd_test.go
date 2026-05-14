package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSummarizeSkillMarkdown(t *testing.T) {
	tests := []struct {
		name string
		md   string
		want string
	}{
		{
			name: "first paragraph joined",
			md: `# My Skill

Two lines here
still same paragraph.

## When`,
			want: "Two lines here still same paragraph.",
		},
		{
			name: "heading right after title",
			md: `# Title

## When to use`,
			want: "-",
		},
		{
			name: "no h1",
			md:   "## No top title\n\nBody.",
			want: "-",
		},
		{
			name: "truncate long",
			md:   "# T\n\n" + strings.Repeat("x", 120),
			want: strings.Repeat("x", 80),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := summarizeSkillMarkdown(tt.md)
			if tt.name == "truncate long" {
				if got != tt.want {
					t.Fatalf("got %q, want %q", got, tt.want)
				}
				if len([]rune(got)) != 80 {
					t.Fatalf("got len %d runes, want 80", len([]rune(got)))
				}
				return
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSkillSummaryLine_File(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "alpha")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "# Alpha\n\nShort blurb for list.\n\n## When\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	got := SkillSummaryLine(Skill{ID: "alpha", Path: skillDir})
	want := "Short blurb for list."
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
