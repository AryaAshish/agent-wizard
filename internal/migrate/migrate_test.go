package migrate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aryaashish/agent-wizard/internal/manifest"
)

func TestRunCreatesBackupAndPreservesManifest(t *testing.T) {
	project := t.TempDir()
	initial := []byte("schemaVersion: 1\ntargetDir: .agents/skills\ninstallMode: manifest-only\nsources: []\nskills: []\n")
	if err := os.WriteFile(filepath.Join(project, manifest.FileName), initial, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Run(project); err != nil {
		t.Fatalf("Run() err=%v", err)
	}

	if _, err := os.Stat(filepath.Join(project, manifest.FileName+".bak")); err != nil {
		t.Fatalf("backup missing: %v", err)
	}
	m, err := manifest.Load(project)
	if err != nil {
		t.Fatalf("manifest load err=%v", err)
	}
	if m.SchemaVersion != CurrentSchemaVersion {
		t.Fatalf("schemaVersion=%d want %d", m.SchemaVersion, CurrentSchemaVersion)
	}
}
