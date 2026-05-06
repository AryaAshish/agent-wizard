package manifest

// EffectiveProfiles returns explicit profiles or a synthetic single-profile default using TargetDir.
func EffectiveProfiles(m Manifest) []Profile {
	if len(m.Profiles) > 0 {
		return m.Profiles
	}
	return []Profile{{
		Name: "default",
		Targets: []ProfileTarget{{
			Kind: "agents",
			Path: m.TargetDir,
		}},
	}}
}
