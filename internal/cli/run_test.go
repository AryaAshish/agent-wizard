package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aryaashish/agent-wizard/internal/config"
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

func TestRunSubcommandHelp(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"help", "add"}, &out); err != nil {
		t.Fatalf("run(help add) error = %v", err)
	}
	if !strings.Contains(out.String(), "Usage: agent-wizard add") {
		t.Fatalf("run(help add) output missing add usage: %q", out.String())
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
	home := t.TempDir()
	origHome := os.Getenv("HOME")
	defer func() { _ = os.Setenv("HOME", origHome) }()
	if err := os.Setenv("HOME", home); err != nil {
		t.Fatalf("Setenv(HOME) error = %v", err)
	}
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
	if err := run([]string{"add", "pr-review", "-android"}, &out); err != nil {
		t.Fatalf("run(add -android) error = %v", err)
	}
	if err := run([]string{"remove", "pr-review"}, &out); err != nil {
		t.Fatalf("run(remove) error = %v", err)
	}
	mPath := filepath.Join(project, "agentskills.yaml")
	b, err := os.ReadFile(mPath)
	if err != nil {
		t.Fatalf("ReadFile(manifest) error = %v", err)
	}
	if !strings.Contains(string(b), "community") {
		t.Fatalf("manifest missing community source: %q", string(b))
	}
	if !strings.Contains(out.String(), "agent-wizard list --source-name community") {
		t.Fatalf("init output missing browse guidance: %q", out.String())
	}
}

func TestRunSourcesAddGitURLFlag(t *testing.T) {
	home := t.TempDir()
	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", home); err != nil {
		t.Fatalf("Setenv(HOME) error = %v", err)
	}
	defer func() { _ = os.Setenv("HOME", origHome) }()

	var out bytes.Buffer
	err := run([]string{
		"sources", "add",
		"--name", "community",
		"--kind", "git",
		"--git-url", "https://github.com/AryaAshish/agent-skills-community.git",
	}, &out)
	if err != nil {
		t.Fatalf("run(sources add git) error = %v", err)
	}

	cfgPath, err := config.DefaultPath()
	if err != nil {
		t.Fatalf("DefaultPath() error = %v", err)
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("Load(config) error = %v", err)
	}
	src, ok := cfg.GetSource("community")
	if !ok {
		t.Fatalf("expected source community to exist")
	}
	if src.GitURL != "https://github.com/AryaAshish/agent-skills-community.git" {
		t.Fatalf("git URL = %q, want %q", src.GitURL, "https://github.com/AryaAshish/agent-skills-community.git")
	}
}

func TestRunCommunityFetch(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"community", "fetch"}, &out); err != nil {
		t.Fatalf("run(community fetch) error = %v", err)
	}
	if !strings.Contains(out.String(), "community starter assets refreshed") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestRunSourcesAddLocalWarnsByDefault(t *testing.T) {
	home := t.TempDir()
	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", home); err != nil {
		t.Fatalf("Setenv(HOME) error = %v", err)
	}
	defer func() { _ = os.Setenv("HOME", origHome) }()

	var out bytes.Buffer
	dir := t.TempDir()
	if err := run([]string{"sources", "add", "--name", "localdev", "--kind", "local", "--path", dir}, &out); err != nil {
		t.Fatalf("run(sources add local) error = %v", err)
	}
	if !strings.Contains(out.String(), "not team-shareable") {
		t.Fatalf("expected shareability warning, got %q", out.String())
	}
}

func TestRunSourcesAddLocalQuietSuppressesWarning(t *testing.T) {
	home := t.TempDir()
	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", home); err != nil {
		t.Fatalf("Setenv(HOME) error = %v", err)
	}
	defer func() { _ = os.Setenv("HOME", origHome) }()

	var out bytes.Buffer
	dir := t.TempDir()
	if err := run([]string{"sources", "add", "--name", "localdev", "--kind", "local", "--path", dir, "--quiet"}, &out); err != nil {
		t.Fatalf("run(sources add local --quiet) error = %v", err)
	}
	if strings.Contains(out.String(), "not team-shareable") {
		t.Fatalf("warning should be suppressed, got %q", out.String())
	}
}
