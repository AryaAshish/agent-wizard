package projectdir

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aryaashish/agent-wizard/internal/manifest"
)

func TestFindManifestRoot(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, manifest.FileName)
	if err := os.WriteFile(manifestPath, []byte("schemaVersion: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, ok := FindManifestRoot(sub)
	if !ok || got != root {
		t.Fatalf("FindManifestRoot(%q) = %q, %v want %q, true", sub, got, ok, root)
	}
}

func TestFindGitRoot(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "pkg", "x")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	got, ok := FindGitRoot(sub)
	if !ok || got != root {
		t.Fatalf("FindGitRoot(%q) = %q, %v want %q, true", sub, got, ok, root)
	}
}

func TestResolveForProjectOps(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "s")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := ResolveForProjectOps(sub)
	if err != ErrNoManifest {
		t.Fatalf("want ErrNoManifest, got %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, manifest.FileName), []byte("schemaVersion: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := ResolveForProjectOps(sub)
	if err != nil || got != root {
		t.Fatalf("ResolveForProjectOps = %q, %v want %q, nil", got, err, root)
	}
}

func TestResolveForInitOrAdd(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "deep")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := ResolveForInitOrAdd(sub)
	if err != nil || got != root {
		t.Fatalf("ResolveForInitOrAdd(no manifest) = %q, %v want %q", got, err, root)
	}
	nogit := t.TempDir()
	nested := filepath.Join(nogit, "x")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	got, err = ResolveForInitOrAdd(nested)
	if err != nil || got != nested {
		t.Fatalf("ResolveForInitOrAdd(no git) = %q want %q", got, nested)
	}
}
