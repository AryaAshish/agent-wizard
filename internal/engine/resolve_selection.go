package engine

import (
	"fmt"
	"sort"

	"github.com/aryaashish/agent-wizard/internal/config"
	"github.com/aryaashish/agent-wizard/internal/manifest"
	"github.com/aryaashish/agent-wizard/internal/model"
	"github.com/aryaashish/agent-wizard/internal/packs"
	"github.com/aryaashish/agent-wizard/internal/skills"
)

// BuildBuckets aggregates discovered skills keyed by bare id across configured sources (manifest order).
func BuildBuckets(cfg config.Config, manifestSources []string) (map[string][]skills.Skill, error) {
	buckets := map[string][]skills.Skill{}
	for _, name := range manifestSources {
		srcCfg, ok := cfg.GetSource(name)
		if !ok {
			return nil, fmt.Errorf("source %q not found", name)
		}
		ms, err := materializeSource(srcCfg)
		if err != nil {
			return nil, fmt.Errorf("materialize source %q: %w", name, err)
		}
		found, err := discoverFromMaterialized(ms)
		if err != nil {
			return nil, fmt.Errorf("discover source %q: %w", name, err)
		}
		for _, s := range found {
			buckets[s.ID] = append(buckets[s.ID], s)
		}
	}
	return buckets, nil
}

func ResolveSkill(ref model.SkillRef, buckets map[string][]skills.Skill) (skills.Skill, error) {
	candidates := buckets[ref.ID]
	if ref.SourceAlias != "" {
		for _, cand := range candidates {
			if cand.SourceName == ref.SourceAlias {
				return cand, nil
			}
		}
		return skills.Skill{}, fmt.Errorf("skill %q not found in source %q", ref.ID, ref.SourceAlias)
	}
	if len(candidates) == 0 {
		return skills.Skill{}, fmt.Errorf("skill %q not found", ref.ID)
	}
	if len(candidates) != 1 {
		return skills.Skill{}, fmt.Errorf("ambiguous skill %q — qualify as source/id", ref.ID)
	}
	return candidates[0], nil
}

// LibraryRoot selects the filesystem root used to resolve `.agent-wizard-pack.yaml` manifests.
// It prefers the first materialized source listed in manifest.Sources order.
func LibraryRoot(cfg config.Config, manifestSources []string) (string, error) {
	for _, name := range manifestSources {
		srcCfg, ok := cfg.GetSource(name)
		if !ok {
			continue
		}
		ms, err := materializeSource(srcCfg)
		if err != nil {
			return "", err
		}
		return ms.Root, nil
	}
	return "", fmt.Errorf("no sources configured")
}

// ExpandSkillSelections combines manifest.skill entries and skills declared by packs.
func ExpandSkillSelections(m manifest.Manifest, cfg config.Config, libraryRoot string) ([]model.SkillRef, error) {
	var refs []model.SkillRef
	seen := map[string]struct{}{}

	for _, p := range m.Packs {
		skillsInPack, err := packs.ResolvePackSkills(libraryRoot, p, map[string]struct{}{})
		if err != nil {
			return nil, fmt.Errorf("pack %q: %w", p, err)
		}
		for _, s := range skillsInPack {
			r, err := model.ParseSkillRef(s)
			if err != nil {
				return nil, err
			}
			key := skillKey(r)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			refs = append(refs, r)
		}
	}

	for _, s := range m.Skills {
		r, err := model.ParseSkillRef(s)
		if err != nil {
			return nil, err
		}
		key := skillKey(r)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		refs = append(refs, r)
	}

	sort.Slice(refs, func(i, j int) bool {
		return skillSortKey(refs[i]) < skillSortKey(refs[j])
	})
	return refs, nil
}

func skillKey(r model.SkillRef) string {
	return r.SourceAlias + "|" + r.ID
}

func skillSortKey(r model.SkillRef) string {
	if r.SourceAlias == "" {
		return r.ID
	}
	return r.SourceAlias + "/" + r.ID
}
