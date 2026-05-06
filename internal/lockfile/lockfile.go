package lockfile

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const FileName = "agentskills.lock"

type Entry struct {
	SkillID        string `yaml:"skillId"`
	SourceName     string `yaml:"sourceName"`
	ResolvedRef    string `yaml:"resolvedRef,omitempty"`
	LocalDigestSHA string `yaml:"digest,omitempty"`
}

type Lockfile struct {
	SchemaVersion int       `yaml:"schemaVersion"`
	GeneratedAt   time.Time `yaml:"generatedAt"`
	Entries       []Entry   `yaml:"entries"`
}

func PathFromDir(projectDir string) string {
	return filepath.Join(projectDir, FileName)
}

func Load(projectDir string) (Lockfile, error) {
	path := PathFromDir(projectDir)
	data, err := os.ReadFile(path)
	if err != nil {
		return Lockfile{}, err
	}
	var lf Lockfile
	if err := yaml.Unmarshal(data, &lf); err != nil {
		return Lockfile{}, err
	}
	if lf.SchemaVersion == 0 {
		lf.SchemaVersion = 1
	}
	return lf, nil
}

func Save(projectDir string, lf Lockfile) error {
	if lf.SchemaVersion == 0 {
		lf.SchemaVersion = 1
	}
	if lf.GeneratedAt.IsZero() {
		lf.GeneratedAt = time.Now().UTC()
	}
	out, err := yaml.Marshal(lf)
	if err != nil {
		return err
	}
	return os.WriteFile(PathFromDir(projectDir), out, 0o644)
}

func EntryIndex(lf Lockfile) map[string]Entry {
	out := map[string]Entry{}
	for _, e := range lf.Entries {
		out[e.SkillID] = e
	}
	return out
}
