package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aryaashish/agent-wizard/internal/community"
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

func TestRunSubcommandHelpList(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"help", "list"}, &out); err != nil {
		t.Fatalf("run(help list) error = %v", err)
	}
	s := out.String()
	if !strings.Contains(s, "aligned") || !strings.Contains(s, "awk") {
		t.Fatalf("run(help list) missing expected prose: %q", s)
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
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# plan-review skill\n\nTest blurb line for discovery.\n\n## More\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var out bytes.Buffer
	if err := run([]string{"list", "--source", root}, &out); err != nil {
		t.Fatalf("run(list) error = %v", err)
	}
	got := strings.TrimSpace(out.String())
	if !strings.HasPrefix(got, "plan-review") || !strings.Contains(got, "Test blurb line for discovery.") {
		t.Fatalf("run(list) output = %q", got)
	}
}

func TestRunListEmptyShowsHint(t *testing.T) {
	root := t.TempDir()
	var out bytes.Buffer
	if err := run([]string{"list", "--source", root}, &out); err != nil {
		t.Fatalf("run(list empty) error = %v", err)
	}
	s := out.String()
	if !strings.Contains(s, "No skills found") || !strings.Contains(s, "create-skill") {
		t.Fatalf("expected empty-state hints, got: %s", s)
	}
}

func TestRunListFilter(t *testing.T) {
	root := t.TempDir()
	for _, id := range []string{"alpha-one", "beta-two"} {
		dir := filepath.Join(root, id)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# s\n\nB.\n"), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
	}

	var out bytes.Buffer
	if err := run([]string{"list", "--source", root, "--filter", "beta"}, &out); err != nil {
		t.Fatalf("run(list --filter) error = %v", err)
	}
	got := strings.TrimSpace(out.String())
	if !strings.Contains(got, "beta-two") || !strings.Contains(got, "B.") {
		t.Fatalf("run(list --filter) output = %q", got)
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
	if err := run([]string{"add", "pr-review", "-android", "--no-sync"}, &out); err != nil {
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

func TestInitMigratesLegacyCommunityGitInGlobalConfig(t *testing.T) {
	home := t.TempDir()
	origHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", home); err != nil {
		t.Fatalf("Setenv(HOME) error = %v", err)
	}
	defer func() { _ = os.Setenv("HOME", origHome) }()

	cfgPath := filepath.Join(home, config.FileName)
	legacyYAML := "schemaVersion: 1\nsources:\n  - name: community\n    kind: git\n    gitUrl: https://github.com/AryaAshish/agent-skills-community.git\n"
	if err := os.WriteFile(cfgPath, []byte(legacyYAML), 0o644); err != nil {
		t.Fatalf("WriteFile config: %v", err)
	}

	proj := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(proj); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	var out bytes.Buffer
	if err := run([]string{"init"}, &out); err != nil {
		t.Fatalf("run(init) error = %v", err)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("Load(config): %v", err)
	}
	src, ok := cfg.GetSource(community.SourceName)
	if !ok {
		t.Fatal("expected community source after init")
	}
	if src.Kind != community.SourceKind {
		t.Fatalf("migrated kind = %q, want %q", src.Kind, community.SourceKind)
	}
	if src.GitURL != "" {
		t.Fatalf("expected empty gitUrl after migrate, got %q", src.GitURL)
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

func TestRunCreateSkill(t *testing.T) {
	root := t.TempDir()
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run([]string{"create-skill", "my-new-skill"}, &out); err != nil {
		t.Fatalf("create-skill: %v", err)
	}
	s := out.String()
	if !strings.Contains(s, "OK  created my-new-skill/SKILL.md") {
		t.Fatalf("unexpected out: %s", s)
	}
	if !strings.Contains(s, "list --source") || !strings.Contains(s, "internal/community/assets/my-new-skill") || !strings.Contains(s, "CONTRIBUTING.md") {
		t.Fatalf("missing hints: %s", s)
	}
	b, err := os.ReadFile(filepath.Join(root, "my-new-skill", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "## When to use") || !strings.Contains(string(b), "# my-new-skill") {
		t.Fatalf("template: %s", string(b))
	}
}

func TestRunCreateSkillConflict(t *testing.T) {
	root := t.TempDir()
	origWd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(origWd) })
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("dup-skill", 0o755); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	err := run([]string{"create-skill", "dup-skill"}, &out)
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("want exists error, got err=%v out=%q", err, out.String())
	}
}
