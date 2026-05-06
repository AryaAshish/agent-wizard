package catalog

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Entry struct {
	ID          string   `yaml:"id"`
	Kind        string   `yaml:"kind"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Source      string   `yaml:"source"`
	Tags        []string `yaml:"tags"`
	Maintainer  string   `yaml:"maintainer"`
	License     string   `yaml:"license"`
	Badge       string   `yaml:"badge"`
}

type IndexDoc struct {
	SchemaVersion int     `yaml:"schemaVersion"`
	Entries       []Entry `yaml:"entries"`
}

func ValidateFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var idx IndexDoc
	if err := yaml.Unmarshal(raw, &idx); err != nil {
		return err
	}
	if idx.SchemaVersion == 0 {
		return fmt.Errorf("schemaVersion missing or zero")
	}
	for _, e := range idx.Entries {
		if strings.TrimSpace(e.ID) == "" {
			return fmt.Errorf("entry missing id")
		}
		if strings.TrimSpace(e.Maintainer) == "" {
			return fmt.Errorf("entry %q missing maintainer", e.ID)
		}
		if strings.TrimSpace(e.License) == "" {
			return fmt.Errorf("entry %q missing license", e.ID)
		}
		for _, tag := range e.Tags {
			if !strings.Contains(tag, ":") {
				return fmt.Errorf("entry %q tag %q invalid (expected taxonomy like stack:android)", e.ID, tag)
			}
		}
	}
	return nil
}
