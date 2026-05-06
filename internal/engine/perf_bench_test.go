package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aryaashish/agent-wizard/internal/config"
	"github.com/aryaashish/agent-wizard/internal/manifest"
)

func BenchmarkSyncDryRunMediumLibrary(b *testing.B) {
	project := b.TempDir()
	lib := filepath.Join(b.TempDir(), "lib")
	if err := os.MkdirAll(lib, 0o755); err != nil {
		b.Fatal(err)
	}
	skills := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		id := fmt.Sprintf("skill-%03d", i)
		skills = append(skills, id)
		dir := filepath.Join(lib, id)
		_ = os.MkdirAll(dir, 0o755)
		_ = os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# "+id), 0o644)
	}

	m := manifest.Manifest{
		SchemaVersion: 1,
		InstallMode:   "manifest-only",
		TargetDir:     ".agents/skills",
		Sources:       []string{"local"},
		Skills:        skills,
	}
	cfg := config.Config{
		SchemaVersion: 1,
		Sources:       []config.Source{{Name: "local", Kind: "local", Path: lib}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Sync(project, m, cfg, ioDiscard{}, SyncOpts{DryRun: true}); err != nil {
			b.Fatal(err)
		}
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }
