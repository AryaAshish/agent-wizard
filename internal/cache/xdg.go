package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

func Root() (string, error) {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".cache")
	}
	dir := filepath.Join(base, "agent-wizard")
	return dir, os.MkdirAll(dir, 0o755)
}

func GitCheckoutDir(remoteURL string) (string, error) {
	root, err := Root()
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(remoteURL))
	h := hex.EncodeToString(sum[:])
	dir := filepath.Join(root, "git", h)
	return dir, os.MkdirAll(dir, 0o755)
}

func ArchiveDir() (string, error) {
	root, err := Root()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(root, "archives")
	return dir, os.MkdirAll(dir, 0o755)
}

func ArchiveExtractDir(url string) (string, error) {
	base, err := ArchiveDir()
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(url))
	h := hex.EncodeToString(sum[:])
	dir := filepath.Join(base, h)
	return dir, os.MkdirAll(dir, 0o755)
}
