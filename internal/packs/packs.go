package packs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const PackManifestName = ".agent-wizard-pack.yaml"

type Pack struct {
	ID      string   `yaml:"id"`
	Version int      `yaml:"schemaVersion"`
	Skills  []string `yaml:"skills"`
	Nested  []string `yaml:"includePacks"`
}

func ResolvePackSkills(rootLibrary string, packID string, seen map[string]struct{}) ([]string, error) {
	if _, ok := seen[packID]; ok {
		return nil, fmt.Errorf("cycle detected in packs at %q", packID)
	}
	seen[packID] = struct{}{}

	dir := filepath.Join(rootLibrary, packID)
	data, err := os.ReadFile(filepath.Join(dir, PackManifestName))
	if err != nil {
		return nil, err
	}
	var p Pack
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	if p.Version == 0 {
		p.Version = 1
	}
	collected := make([]string, 0, len(p.Skills))
	for _, skill := range p.Skills {
		collected = append(collected, strings.TrimSpace(skill))
	}
	for _, nested := range p.Nested {
		sk, err := ResolvePackSkills(rootLibrary, strings.TrimSpace(nested), seen)
		if err != nil {
			return nil, err
		}
		collected = append(collected, sk...)
	}
	return dedupePreserveOrder(collected), nil
}

func dedupePreserveOrder(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// LocatePackRoots walks the library root and maps pack id → directory containing pack manifest.
func LocatePackRoots(libraryRoot string) (map[string]string, error) {
	out := map[string]string{}
	err := filepath.WalkDir(libraryRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() != PackManifestName {
			return nil
		}
		dir := filepath.Dir(path)
		id := filepath.Base(dir)
		if id == "" || id == "." {
			return fmt.Errorf("invalid pack id path for %s", path)
		}
		if existing, exists := out[id]; exists && existing != dir {
			return fmt.Errorf("duplicate pack id %q", id)
		}
		out[id] = dir
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
