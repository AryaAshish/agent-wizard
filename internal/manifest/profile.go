package manifest

type ProfileTarget struct {
	Kind string `yaml:"kind"`
	Path string `yaml:"path"`
}

type Profile struct {
	Name    string          `yaml:"name"`
	Targets []ProfileTarget `yaml:"targets,omitempty"`
}

type Hooks struct {
	PreSync  []string `yaml:"preSync,omitempty"`
	PostSync []string `yaml:"postSync,omitempty"`
}
