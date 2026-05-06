package migrate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aryaashish/agent-wizard/internal/manifest"
)

const CurrentSchemaVersion = 1

// Run creates a backup and re-saves the manifest with defaults applied (no-op upgrade for schema 1).
func Run(projectDir string) error {
	src := manifest.PathFromDir(projectDir)
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	backup := filepath.Join(projectDir, manifest.FileName+".bak")
	if err := os.WriteFile(backup, data, 0o644); err != nil {
		return fmt.Errorf("backup manifest: %w", err)
	}
	m, err := manifest.Load(projectDir)
	if err != nil {
		return err
	}
	if m.SchemaVersion < CurrentSchemaVersion {
		m.SchemaVersion = CurrentSchemaVersion
	}
	return manifest.Save(projectDir, m)
}
