package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const FileName = ".agent-wizard-config.yaml"

type Source struct {
	Name string `yaml:"name"`
	Kind string `yaml:"kind"`
	Path string `yaml:"path"`
	// Git source fields (kind: git)
	GitURL string `yaml:"gitUrl,omitempty"`
	GitRef string `yaml:"gitRef,omitempty"`
	Subdir string `yaml:"subdir,omitempty"`
	// Archive source (kind: archive) — https URL to zip
	ArchiveURL string `yaml:"archiveUrl,omitempty"`
}

type Config struct {
	SchemaVersion int      `yaml:"schemaVersion"`
	Sources       []Source `yaml:"sources"`
}

func DefaultPath() (string, error) {
	// Tests set HOME for an isolated config dir. On Windows, os.UserHomeDir()
	// ignores HOME (it uses USERPROFILE), so prefer HOME when set.
	if h := os.Getenv("HOME"); h != "" {
		return filepath.Join(h, FileName), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, FileName), nil
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{SchemaVersion: 1, Sources: []Source{}}, nil
		}
		return Config{}, err
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return Config{}, err
	}
	if c.SchemaVersion == 0 {
		c.SchemaVersion = 1
	}
	return c, nil
}

func Save(path string, c Config) error {
	if c.SchemaVersion == 0 {
		c.SchemaVersion = 1
	}
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}

func (c Config) GetSource(name string) (Source, bool) {
	for _, s := range c.Sources {
		if s.Name == name {
			return s, true
		}
	}
	return Source{}, false
}
