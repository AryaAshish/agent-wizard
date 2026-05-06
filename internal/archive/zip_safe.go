package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ExtractRemoteZip downloads zip to a temp file then extracts safely under dest.
func ExtractRemoteZip(url string, dest string, maxBytes int64) error {
	if maxBytes <= 0 {
		maxBytes = 50 * 1024 * 1024
	}
	client := &http.Client{Timeout: 2 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %s", resp.Status)
	}
	tmp, err := os.CreateTemp("", "agent-wizard-archive-*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()
	n, err := io.Copy(tmp, io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return err
	}
	if n > maxBytes {
		return fmt.Errorf("archive exceeds max size")
	}
	if _, err := tmp.Seek(0, 0); err != nil {
		return err
	}
	z, err := zip.NewReader(tmp, n)
	if err != nil {
		return err
	}
	return ExtractZip(z, dest)
}

func ExtractZip(z *zip.Reader, dest string) error {
	for _, f := range z.File {
		cleanName := filepath.ToSlash(filepath.Clean(f.Name))
		if cleanName == "." || strings.HasPrefix(cleanName, "../") || strings.Contains(cleanName, "/../") {
			return fmt.Errorf("unsafe zip path %q", f.Name)
		}
		target := filepath.Join(dest, filepath.FromSlash(cleanName))
		absDest, err := filepath.Abs(dest)
		if err != nil {
			return err
		}
		absTarget, err := filepath.Abs(target)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(absDest, absTarget)
		if err != nil || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
			return fmt.Errorf("path escapes dest: %q", f.Name)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(absTarget, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(absTarget), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		if err := writeFileAtomically(absTarget, rc); err != nil {
			rc.Close()
			return err
		}
		_ = rc.Close()
	}
	return nil
}

func writeFileAtomically(path string, r io.Reader) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".extract-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := io.Copy(tmp, r); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
