package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aryaashish/agent-wizard/internal/config"
)

func TestEndUserFlow_EmbeddedCommunity_ListFilterAddSync(t *testing.T) {
	project := t.TempDir()
	home := t.TempDir()
	restore := setEnvAndCwd(t, map[string]string{"HOME": home}, project)
	defer restore()

	var out bytes.Buffer
	mustRun(t, []string{"init"}, &out)
	out.Reset()

	mustRun(t, []string{"list", "--source-name", "community", "--filter", "pr-review"}, &out)
	if !strings.Contains(out.String(), "pr-review") {
		t.Fatalf("expected pr-review in list output, got:\n%s", out.String())
	}
	out.Reset()

	mustRun(t, []string{"add", "pr-review", "--source", "community"}, &out)
	out.Reset()
	mustRun(t, []string{"sync"}, &out)

	skillPath := filepath.Join(project, ".agents", "skills", "pr-review", "SKILL.md")
	b, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("synced community skill missing: %v", err)
	}
	if !strings.Contains(string(b), "When to use") {
		t.Fatalf("expected launch-ready SKILL sections, got head: %.200q", string(b))
	}
}

func TestEndUserFlow_EmbeddedCommunity_AddColdStartNoInit(t *testing.T) {
	project := t.TempDir()
	home := t.TempDir()
	restore := setEnvAndCwd(t, map[string]string{"HOME": home}, project)
	defer restore()

	var out bytes.Buffer
	if err := run([]string{"add", "pr-review", "--source", "community"}, &out); err != nil {
		t.Fatalf("add (cold start): %v; out=%s", err, out.String())
	}
	s := out.String()
	if !strings.Contains(s, "✔ Installed pr-review") || !strings.Contains(s, "→ Open:") || !strings.Contains(s, ".agents/skills/pr-review/SKILL.md") {
		t.Fatalf("unexpected add output:\n%s", s)
	}
	skillPath := filepath.Join(project, ".agents", "skills", "pr-review", "SKILL.md")
	b, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("synced community skill missing: %v", err)
	}
	if !strings.Contains(string(b), "When to use") {
		t.Fatalf("expected launch-ready SKILL sections, got head: %.200q", string(b))
	}
}

func TestEndUserFlow_EmbeddedCommunity_PackAddSync(t *testing.T) {
	project := t.TempDir()
	home := t.TempDir()
	restore := setEnvAndCwd(t, map[string]string{"HOME": home}, project)
	defer restore()

	var out bytes.Buffer
	mustRun(t, []string{"init"}, &out)
	out.Reset()
	mustRun(t, []string{"pack", "add", "android-starter"}, &out)
	out.Reset()
	mustRun(t, []string{"sync"}, &out)

	want := []string{"pr-review", "plan-review", "launch-ready", "cursor-rules-hooks", "github-actions-matrix"}
	for _, id := range want {
		p := filepath.Join(project, ".agents", "skills", id, "SKILL.md")
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("missing synced skill %q: %v", id, err)
		}
	}
}

func TestCLI_IdempotentAddSync(t *testing.T) {
	project := t.TempDir()
	home := t.TempDir()
	restore := setEnvAndCwd(t, map[string]string{"HOME": home}, project)
	defer restore()

	var out bytes.Buffer
	mustRun(t, []string{"init"}, &out)
	out.Reset()
	mustRun(t, []string{"add", "pr-review", "--source", "community"}, &out)
	out.Reset()
	mustRun(t, []string{"add", "pr-review", "--source", "community"}, &out)
	out.Reset()
	mustRun(t, []string{"sync"}, &out)
	out.Reset()
	mustRun(t, []string{"sync"}, &out)
	p := filepath.Join(project, ".agents", "skills", "pr-review", "SKILL.md")
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("skill missing after idempotent ops: %v", err)
	}
}

func TestCLI_ErrorHints(t *testing.T) {
	tests := []struct {
		name          string
		prep          func(t *testing.T, buf *bytes.Buffer)
		args          []string
		wantSubstring string
	}{
		{
			name: "sync_unknown_skill",
			prep: func(t *testing.T, buf *bytes.Buffer) {
				mustRun(t, []string{"init"}, buf)
				buf.Reset()
				mustRun(t, []string{"add", "definitely-nonexistent-skill-xyz", "--source", "community", "--no-sync"}, buf)
			},
			args:          []string{"sync"},
			wantSubstring: "hint:",
		},
		{
			name: "sync_corrupt_manifest_yaml",
			prep: func(t *testing.T, buf *bytes.Buffer) {
				mustRun(t, []string{"init"}, buf)
				writeFile(t, filepath.Join(mustDetectWd(t), "agentskills.yaml"), "not: yaml: [[[ broken")
			},
			args:          []string{"sync"},
			wantSubstring: "hint:",
		},
		{
			name: "sync_unknown_manifest_source_second_in_list",
			prep: func(t *testing.T, buf *bytes.Buffer) {
				mustRun(t, []string{"init"}, buf)
				wd := mustDetectWd(t)
				raw := readFile(t, filepath.Join(wd, "agentskills.yaml"))
				updated := strings.Replace(raw, "sources:\n    - community", "sources:\n    - community\n    - no-such-registry-zzz", 1)
				if updated == raw {
					t.Fatalf("could not inject unknown source — got:\n%s", raw)
				}
				writeFile(t, filepath.Join(wd, "agentskills.yaml"), updated)
			},
			args:          []string{"sync"},
			wantSubstring: "hint:",
		},
		{
			name: "add_corrupt_manifest_yaml",
			prep: func(t *testing.T, buf *bytes.Buffer) {
				writeFile(t, filepath.Join(mustDetectWd(t), "agentskills.yaml"), "not: yaml: [[[ broken")
			},
			args:          []string{"add", "pr-review", "--source", "community"},
			wantSubstring: "hint:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := t.TempDir()
			home := t.TempDir()
			restore := setEnvAndCwd(t, map[string]string{"HOME": home}, project)
			defer restore()
			var setup bytes.Buffer
			tt.prep(t, &setup)
			setup.Reset()
			err := run(tt.args, &setup)
			if err == nil {
				t.Fatalf("expected error, got nil; out=%q", setup.String())
			}
			msg := err.Error()
			if tt.wantSubstring != "" && !strings.Contains(msg, tt.wantSubstring) {
				t.Fatalf("error %q missing %q", msg, tt.wantSubstring)
			}
		})
	}
}

func mustDetectWd(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	return wd
}

func TestEndUserFlow_LocalSource_Pack_Lock_Sync_StatusJSON(t *testing.T) {
	project := t.TempDir()
	home := t.TempDir()
	library := filepath.Join(t.TempDir(), "library")

	mkSkill(t, library, "pr-review")
	mkSkill(t, library, "plan-review")
	mkPack(t, library, "android-starter", "pr-review", "plan-review")

	restore := setEnvAndCwd(t, map[string]string{"HOME": home}, project)
	defer restore()

	var out bytes.Buffer
	mustRun(t, []string{"init"}, &out)
	out.Reset()
	mustRun(t, []string{"sources", "add", "--name", "local-lib", "--kind", "local", "--path", library}, &out)
	out.Reset()

	// Wire source into manifest (sources command is user-level config by design).
	m := readFile(t, filepath.Join(project, "agentskills.yaml"))
	m = replaceSourcesList(m, "sources:\n    - local-lib")
	writeFile(t, filepath.Join(project, "agentskills.yaml"), m)

	mustRun(t, []string{"pack", "add", "android-starter"}, &out)
	out.Reset()
	mustRun(t, []string{"lock"}, &out)
	out.Reset()
	mustRun(t, []string{"sync"}, &out)
	out.Reset()
	mustRun(t, []string{"status", "--json"}, &out)

	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("status --json parse failed: %v\npayload=%s", err, out.String())
	}
	if got := payload["installMode"]; got != "manifest-only" {
		t.Fatalf("installMode=%v want manifest-only", got)
	}
	if _, err := os.Stat(filepath.Join(project, ".agents", "skills", "pr-review", "SKILL.md")); err != nil {
		t.Fatalf("synced skill missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(project, ".agents", "skills", "plan-review", "SKILL.md")); err != nil {
		t.Fatalf("synced skill missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(project, "agentskills.lock")); err != nil {
		t.Fatalf("lockfile missing: %v", err)
	}
}

func TestNegative_AmbiguousBareSkillNeedsNamespace(t *testing.T) {
	project := t.TempDir()
	home := t.TempDir()
	libA := filepath.Join(t.TempDir(), "lib-a")
	libB := filepath.Join(t.TempDir(), "lib-b")
	mkSkill(t, libA, "pr-review")
	mkSkill(t, libB, "pr-review")

	restore := setEnvAndCwd(t, map[string]string{"HOME": home}, project)
	defer restore()

	var out bytes.Buffer
	mustRun(t, []string{"init"}, &out)
	out.Reset()
	mustRun(t, []string{"sources", "add", "--name", "a", "--kind", "local", "--path", libA}, &out)
	out.Reset()
	mustRun(t, []string{"sources", "add", "--name", "b", "--kind", "local", "--path", libB}, &out)
	out.Reset()
	m := readFile(t, filepath.Join(project, "agentskills.yaml"))
	m = replaceSourcesList(m, "sources:\n    - a\n    - b")
	writeFile(t, filepath.Join(project, "agentskills.yaml"), m)
	mustRun(t, []string{"add", "pr-review", "--no-sync"}, &out)
	out.Reset()

	err := run([]string{"sync", "--dry-run"}, &out)
	if err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("expected ambiguous error, got=%v out=%q", err, out.String())
	}
}

func TestNegative_StrictLockDigestMismatchAndDriftExitCode(t *testing.T) {
	project := t.TempDir()
	home := t.TempDir()
	lib := filepath.Join(t.TempDir(), "lib")
	mkSkill(t, lib, "pr-review")

	restore := setEnvAndCwd(t, map[string]string{"HOME": home}, project)
	defer restore()

	var out bytes.Buffer
	mustRun(t, []string{"init"}, &out)
	mustRun(t, []string{"sources", "add", "--name", "local-lib", "--kind", "local", "--path", lib}, &out)
	m := readFile(t, filepath.Join(project, "agentskills.yaml"))
	m = replaceSourcesList(m, "sources:\n    - local-lib")
	writeFile(t, filepath.Join(project, "agentskills.yaml"), m)
	mustRun(t, []string{"add", "pr-review"}, &out)
	mustRun(t, []string{"lock"}, &out)

	// Mutate source skill markdown after lock to force digest mismatch.
	writeFile(t, filepath.Join(lib, "pr-review", "SKILL.md"), "# changed\n")

	out.Reset()
	err := run([]string{"sync", "--dry-run", "--strict-lock"}, &out)
	if err == nil || !strings.Contains(err.Error(), "digest") {
		t.Fatalf("expected strict-lock digest mismatch, got=%v out=%q", err, out.String())
	}

	out.Reset()
	err = run([]string{"status", "--check-drifts", "--strict-digest"}, &out)
	if err == nil {
		t.Fatalf("expected drift error")
	}
	ec, ok := err.(ExitCoder)
	if !ok || ec.Code() != 3 {
		t.Fatalf("expected exit code 3 for drift, got err=%v", err)
	}
}

func setEnvAndCwd(t *testing.T, env map[string]string, cwd string) func() {
	t.Helper()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(cwd); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	oldEnv := map[string]string{}
	for k, v := range env {
		oldEnv[k] = os.Getenv(k)
		if err := os.Setenv(k, v); err != nil {
			t.Fatalf("Setenv %s: %v", k, err)
		}
	}
	return func() {
		_ = os.Chdir(oldWd)
		for k, v := range oldEnv {
			_ = os.Setenv(k, v)
		}
	}
}

func mkPack(t *testing.T, root, id string, skills ...string) {
	t.Helper()
	dir := filepath.Join(root, id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	var b strings.Builder
	b.WriteString("schemaVersion: 1\n")
	b.WriteString("id: ")
	b.WriteString(id)
	b.WriteString("\nskills:\n")
	for _, s := range skills {
		b.WriteString("  - ")
		b.WriteString(s)
		b.WriteString("\n")
	}
	if err := os.WriteFile(filepath.Join(dir, ".agent-wizard-pack.yaml"), []byte(b.String()), 0o644); err != nil {
		t.Fatalf("WriteFile pack: %v", err)
	}
}

func mkSkill(t *testing.T, root, id string) {
	t.Helper()
	dir := filepath.Join(root, id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# "+id+"\n"), 0o644); err != nil {
		t.Fatalf("WriteFile SKILL.md: %v", err)
	}
}

func mustRun(t *testing.T, args []string, out *bytes.Buffer) {
	t.Helper()
	if err := run(args, out); err != nil {
		t.Fatalf("run(%v) err=%v out=%q", args, err, out.String())
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", path, err)
	}
	return string(b)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}
}

func init() {
	// compile-time check that config still supports source kinds.
	_ = config.Source{}
}

func replaceSourcesList(in, replacement string) string {
	out := strings.Replace(in, "sources:\n    - community", replacement, 1)
	if out != in {
		return out
	}
	return strings.Replace(in, "sources: []", replacement, 1)
}
