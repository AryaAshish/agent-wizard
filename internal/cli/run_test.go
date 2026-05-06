package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"help"}, &out); err != nil {
		t.Fatalf("run(help) error = %v", err)
	}
	if !strings.Contains(out.String(), "Usage:") {
		t.Fatalf("run(help) output missing Usage: %q", out.String())
	}
}

func TestRunList(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "plan-review")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# skill"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var out bytes.Buffer
	if err := run([]string{"list", "--source", root}, &out); err != nil {
		t.Fatalf("run(list) error = %v", err)
	}
	if got := strings.TrimSpace(out.String()); got != "plan-review" {
		t.Fatalf("run(list) output = %q, want %q", got, "plan-review")
	}
}

func TestRunInitAddRemove(t *testing.T) {
	project := t.TempDir()
	orig, _ := os.Getwd()
	defer func() { _ = os.Chdir(orig) }()
	if err := os.Chdir(project); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}

	var out bytes.Buffer
	if err := run([]string{"init"}, &out); err != nil {
		t.Fatalf("run(init) error = %v", err)
	}
	if err := run([]string{"add", "pr-review"}, &out); err != nil {
		t.Fatalf("run(add) error = %v", err)
	}
	if err := run([]string{"remove", "pr-review"}, &out); err != nil {
		t.Fatalf("run(remove) error = %v", err)
	}
}
