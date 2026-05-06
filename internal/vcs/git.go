package vcs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EnsureCheckout clones or updates repo and returns checkout directory and HEAD sha.
func EnsureCheckout(destDir, repoURL, gitRef, subdir string) (string, string, error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", "", err
	}
	gitMarker := filepath.Join(destDir, ".git")
	if _, err := os.Stat(gitMarker); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", repoURL, destDir)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", "", fmt.Errorf("git clone: %w", err)
		}
	} else if err != nil {
		return "", "", err
	} else {
		cmd := exec.Command("git", "-C", destDir, "fetch", "--all", "--prune")
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", "", fmt.Errorf("git fetch: %w", err)
		}
	}
	ref := strings.TrimSpace(gitRef)
	if ref == "" {
		ref = "HEAD"
	}
	co := exec.Command("git", "-C", destDir, "checkout", ref)
	co.Stderr = os.Stderr
	if err := co.Run(); err != nil {
		return "", "", fmt.Errorf("git checkout %q: %w", ref, err)
	}
	headOut, err := exec.Command("git", "-C", destDir, "rev-parse", "HEAD").Output()
	if err != nil {
		return "", "", fmt.Errorf("git rev-parse: %w", err)
	}
	head := strings.TrimSpace(string(headOut))
	finalPath := filepath.Join(destDir, filepath.FromSlash(subdir))
	return finalPath, head, nil
}
