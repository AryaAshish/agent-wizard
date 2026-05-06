package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocalPathSourceDiscoverFindsSkills(t *testing.T) {
	root := t.TempDir()
	mkSkill(t, root, "android-review")
	mkSkill(t, root, "backend-checks")

	source := NewLocalPathSource(root)
	got, err := source.Discover()
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("Discover() count = %d, want 2", len(got))
	}
	if got[0].ID != "android-review" || got[1].ID != "backend-checks" {
		t.Fatalf("Discover() ids = %#v", got)
	}
}

func TestLocalPathSourceDiscoverDuplicateIDFails(t *testing.T) {
	root := t.TempDir()
	mkSkillAt(t, filepath.Join(root, "source-a", "my-skill"))
	mkSkillAt(t, filepath.Join(root, "source-b", "my-skill"))

	source := NewLocalPathSource(root)
	_, err := source.Discover()
	if err == nil {
		t.Fatal("Discover() expected duplicate id error, got nil")
	}
}

func mkSkill(t *testing.T, root string, id string) {
	t.Helper()
	mkSkillAt(t, filepath.Join(root, id))
}

func mkSkillAt(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	p := filepath.Join(dir, "SKILL.md")
	if err := os.WriteFile(p, []byte("# Skill\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}
