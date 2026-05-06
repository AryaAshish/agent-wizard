package drift

import (
	"fmt"
	"strings"

	"github.com/aryaashish/agent-wizard/internal/config"
	"github.com/aryaashish/agent-wizard/internal/engine"
	"github.com/aryaashish/agent-wizard/internal/hashutil"
	"github.com/aryaashish/agent-wizard/internal/lockfile"
	"github.com/aryaashish/agent-wizard/internal/manifest"
	"github.com/aryaashish/agent-wizard/internal/model"
)

// Evaluate compares lockfile entries with the live graph resolved from sources.
// Returns human-readable issues; ok is true when there are no issues.
func Evaluate(projectDir string, m manifest.Manifest, cfg config.Config, strictDigests bool) ([]string, bool, error) {
	lf, err := lockfile.Load(projectDir)
	if err != nil {
		return nil, false, fmt.Errorf("load lockfile: %w", err)
	}
	libRoot, err := engine.LibraryRoot(cfg, m.Sources)
	if err != nil {
		return nil, false, err
	}
	buckets, err := engine.BuildBuckets(cfg, m.Sources)
	if err != nil {
		return nil, false, err
	}
	exp, err := engine.ExpandSkillSelections(m, cfg, libRoot)
	if err != nil {
		return nil, false, err
	}

	selected := map[string]struct{}{}
	for _, ref := range exp {
		selected[ref.ID] = struct{}{}
	}

	var msgs []string
	for _, entry := range lf.Entries {
		if _, ok := selected[entry.SkillID]; !ok {
			msgs = append(msgs, fmt.Sprintf("locked skill %q not selected in manifest", entry.SkillID))
			continue
		}
		ref, err := model.ParseSkillRef(entry.SkillID)
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("%q: %v", entry.SkillID, err))
			continue
		}
		sk, err := engine.ResolveSkill(ref, buckets)
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("%q: %v", entry.SkillID, err))
			continue
		}
		if sk.SourceName != entry.SourceName {
			msgs = append(msgs, fmt.Sprintf("%q source changed (lock=%s live=%s)", entry.SkillID, entry.SourceName, sk.SourceName))
		}
		switch {
		case sk.ResolvedRef == "local":
			// no git sha to compare
		case strings.HasPrefix(sk.ResolvedRef, "archive:"):
			if entry.ResolvedRef != "" && entry.ResolvedRef != sk.ResolvedRef {
				msgs = append(msgs, fmt.Sprintf("%q archive ref drift lock=%s live=%s", entry.SkillID, entry.ResolvedRef, sk.ResolvedRef))
			}
		default:
			if entry.ResolvedRef != "" && entry.ResolvedRef != sk.ResolvedRef {
				msgs = append(msgs, fmt.Sprintf("%q ref drift lock=%s live=%s", entry.SkillID, entry.ResolvedRef, sk.ResolvedRef))
			}
		}
		if strictDigests && entry.LocalDigestSHA != "" && sk.ResolvedRef == "local" {
			d, err := hashutil.SkillMarkdownDigest(sk.Path)
			if err != nil {
				return nil, false, err
			}
			if d != entry.LocalDigestSHA {
				msgs = append(msgs, fmt.Sprintf("%q SKILL.md digest drift", entry.SkillID))
			}
		}
	}
	return msgs, len(msgs) == 0, nil
}
