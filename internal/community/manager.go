package community

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	SourceName = "community"
	SourceKind = "community"
)

//go:embed assets/** assets/**/.agent-wizard-pack.yaml
var assets embed.FS

func cacheRoot() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "agent-wizard", "community"), nil
}

func Extract(force bool) (string, error) {
	root, err := cacheRoot()
	if err != nil {
		return "", err
	}
	dst := filepath.Join(root, "library")
	if force {
		if err := os.RemoveAll(dst); err != nil {
			return "", err
		}
	}
	// Fast path when already extracted.
	if _, err := os.Stat(filepath.Join(dst, "pr-review", "SKILL.md")); err == nil {
		return dst, nil
	}
	if err := os.RemoveAll(dst); err != nil {
		return "", err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return "", err
	}
	if err := fs.WalkDir(assets, "assets", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == "assets" {
			return nil
		}
		rel, err := filepath.Rel("assets", path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := assets.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, b, 0o644)
	}); err != nil {
		return "", err
	}
	return dst, nil
}
