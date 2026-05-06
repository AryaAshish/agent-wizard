package skills

type Skill struct {
	ID          string
	Path        string
	SourceName  string
	ResolvedRef string
}

type Source interface {
	Discover() ([]Skill, error)
}
