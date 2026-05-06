package engine

import (
	"strings"

	"github.com/aryaashish/agent-wizard/internal/config"
	"github.com/aryaashish/agent-wizard/internal/hashutil"
	"github.com/aryaashish/agent-wizard/internal/lockfile"
	"github.com/aryaashish/agent-wizard/internal/manifest"
)

// GenerateLockfile writes agentskills.lock from the current resolved skill graph.
func GenerateLockfile(projectDir string, m manifest.Manifest, cfg config.Config) error {
	libRoot, err := LibraryRoot(cfg, m.Sources)
	if err != nil {
		return err
	}
	buckets, err := BuildBuckets(cfg, m.Sources)
	if err != nil {
		return err
	}
	exp, err := ExpandSkillSelections(m, cfg, libRoot)
	if err != nil {
		return err
	}
	var entries []lockfile.Entry
	for _, ref := range exp {
		sk, err := ResolveSkill(ref, buckets)
		if err != nil {
			return err
		}
		entry := lockfile.Entry{
			SkillID:     ref.ID,
			SourceName:  sk.SourceName,
			ResolvedRef: strings.TrimSpace(sk.ResolvedRef),
		}
		if sk.ResolvedRef == "local" {
			d, err := hashutil.SkillMarkdownDigest(sk.Path)
			if err != nil {
				return err
			}
			entry.LocalDigestSHA = d
		}
		entries = append(entries, entry)
	}
	lf := lockfile.Lockfile{SchemaVersion: 1, Entries: entries}
	return lockfile.Save(projectDir, lf)
}
