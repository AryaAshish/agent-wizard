package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
}

func AddUnique(items []string, value string) []string {
	for _, it := range items {
		if it == value {
			return items
		}
	}
	return append(items, value)
}

func RemoveValue(items []string, value string) []string {
	out := make([]string, 0, len(items))
	for _, it := range items {
		if it != value {
			out = append(out, it)
		}
	}
	return out
}

func ValidateICP(mode string) error {
	switch strings.ToLower(mode) {
	case "solo", "team", "enterprise":
		return nil
	default:
		return fmt.Errorf("unsupported icp %q, expected solo|team|enterprise", mode)
	}
}
