package packs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePackSkillsNestedDedup(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "base"))
	write(t, filepath.Join(root, "base", PackManifestName), []byte(`
schemaVersion: 1
id: base
skills:
  - a
includePacks:
  - child
`))
	mkdir(t, filepath.Join(root, "child"))
	write(t, filepath.Join(root, "child", PackManifestName), []byte(`
schemaVersion: 1
id: child
skills:
  - a
  - b
`))

	got, err := ResolvePackSkills(root, "base", map[string]struct{}{})
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got=%#v len=%d", got, len(got))
	}
}

func mkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func write(t *testing.T, path string, b []byte) {
	t.Helper()
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}
}
