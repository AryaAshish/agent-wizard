package manifest

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const FileName = "agentskills.yaml"

type Manifest struct {
	SchemaVersion  int       `yaml:"schemaVersion"`
	TargetDir      string    `yaml:"targetDir"`
	InstallMode    string    `yaml:"installMode"`
	Sources        []string  `yaml:"sources"`
	Skills         []string  `yaml:"skills"`
	Packs          []string  `yaml:"packs,omitempty"`
	Profiles       []Profile `yaml:"profiles,omitempty"`
	Hooks          Hooks     `yaml:"hooks,omitempty"`
	AllowedSources []string  `yaml:"allowedSources,omitempty"`
}

func PathFromDir(projectDir string) string {
	return filepath.Join(projectDir, FileName)
}

func Load(projectDir string) (Manifest, error) {
	path := PathFromDir(projectDir)
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return Manifest{}, err
	}
	if m.SchemaVersion == 0 {
		m.SchemaVersion = 1
	}
	if m.TargetDir == "" {
		m.TargetDir = ".agents/skills"
	}
	if m.InstallMode == "" {
		m.InstallMode = "manifest-only"
	}
	return m, nil
}

func Save(projectDir string, m Manifest) error {
	if m.SchemaVersion == 0 {
		m.SchemaVersion = 1
	}
	if m.TargetDir == "" {
		m.TargetDir = ".agents/skills"
	}
	if m.InstallMode == "" {
		m.InstallMode = "manifest-only"
	}
	path := PathFromDir(projectDir)
	out, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}

func Init(projectDir string) (Manifest, error) {
	path := PathFromDir(projectDir)
	if _, err := os.Stat(path); err == nil {
		return Manifest{}, errors.New("manifest already exists")
	}
	m := Manifest{
		SchemaVersion: 1,
		TargetDir:     ".agents/skills",
		InstallMode:   "manifest-only",
		Sources:       []string{},
		Skills:        []string{},
	}
	return m, Save(projectDir, m)
}
