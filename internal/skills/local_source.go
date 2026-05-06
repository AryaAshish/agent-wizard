package skills

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

type LocalPathSource struct {
	root string
}

func NewLocalPathSource(root string) LocalPathSource {
	return LocalPathSource{root: root}
}

func (s LocalPathSource) Discover() ([]Skill, error) {
	var discovered []Skill
	seen := map[string]struct{}{}

	err := filepath.WalkDir(s.root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() != "SKILL.md" {
			return nil
		}

		id := filepath.Base(filepath.Dir(path))
		if id == "." || strings.TrimSpace(id) == "" {
			return fmt.Errorf("invalid skill id from path %s", path)
		}
		if _, exists := seen[id]; exists {
			return fmt.Errorf("duplicate skill id %q", id)
		}
		seen[id] = struct{}{}

		discovered = append(discovered, Skill{
			ID:   id,
			Path: filepath.Dir(path),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(discovered, func(i, j int) bool {
		return discovered[i].ID < discovered[j].ID
	})
	return discovered, nil
}
