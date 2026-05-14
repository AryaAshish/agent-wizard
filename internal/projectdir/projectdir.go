// Package projectdir resolves the agent-wizard project directory from cwd.
package projectdir

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/aryaashish/agent-wizard/internal/manifest"
)

// ErrNoManifest is returned when agentskills.yaml cannot be found walking up from start.
var ErrNoManifest = errors.New("agentskills.yaml not found")

// FindManifestRoot walks parents from start until manifest.FileName exists.
func FindManifestRoot(start string) (string, bool) {
	dir := filepath.Clean(start)
	for {
		p := filepath.Join(dir, manifest.FileName)
		if _, err := os.Stat(p); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}

// FindGitRoot walks parents until ".git" exists (file or directory).
func FindGitRoot(start string) (string, bool) {
	dir := filepath.Clean(start)
	for {
		gitPath := filepath.Join(dir, ".git")
		if fi, err := os.Stat(gitPath); err == nil && (fi.IsDir() || fi.Mode().IsRegular()) {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}

// ResolveForProjectOps returns the directory containing agentskills.yaml when walking up from start.
func ResolveForProjectOps(start string) (string, error) {
	if dir, ok := FindManifestRoot(start); ok {
		return dir, nil
	}
	return "", ErrNoManifest
}

// ResolveForInitOrAdd returns manifest dir if present, else nearest git root, else start (cleaned).
func ResolveForInitOrAdd(start string) (string, error) {
	start = filepath.Clean(start)
	if dir, ok := FindManifestRoot(start); ok {
		return dir, nil
	}
	if dir, ok := FindGitRoot(start); ok {
		return dir, nil
	}
	return start, nil
}
