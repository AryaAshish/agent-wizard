package engine

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aryaashish/agent-wizard/internal/config"
	"github.com/aryaashish/agent-wizard/internal/hashutil"
	"github.com/aryaashish/agent-wizard/internal/lockfile"
	"github.com/aryaashish/agent-wizard/internal/manifest"
)

// SyncOpts controls sync behaviour.
type SyncOpts struct {
	DryRun     bool
	Prune      bool
	StrictLock bool
}

// Sync resolves selections, optionally validates lockfile, and copies skills into each configured profile target.
func Sync(projectDir string, m manifest.Manifest, cfg config.Config, stdout io.Writer, opts SyncOpts) error {
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

	lf, lfErr := lockfile.Load(projectDir)
	lockIndex := map[string]lockfile.Entry{}
	if lfErr == nil {
		lockIndex = lockfile.EntryIndex(lf)
	}

	selectedIDsForPrune := map[string]struct{}{}

	runHooks := func(cmds []string, phase string) error {
		for _, c := range cmds {
			cmd := strings.TrimSpace(c)
			if cmd == "" {
				continue
			}
			if opts.DryRun {
				fmt.Fprintf(stdout, "would run %s hook: %s\n", phase, cmd)
				continue
			}
			h := execShell(cmd, projectDir)
			h.Stdout = stdout
			h.Stderr = os.Stderr
			if err := h.Run(); err != nil {
				return fmt.Errorf("%s hook %q: %w", phase, cmd, err)
			}
		}
		return nil
	}

	if err := runHooks(m.Hooks.PreSync, "preSync"); err != nil {
		return err
	}

	for _, ref := range exp {
		sk, err := ResolveSkill(ref, buckets)
		if err != nil {
			return err
		}
		selectedIDsForPrune[sk.ID] = struct{}{}

		if opts.StrictLock && lfErr == nil {
			e, ok := lockIndex[ref.ID]
			if !ok {
				return fmt.Errorf("strict-lock: skill %q missing from lockfile", ref.ID)
			}
			if e.SourceName != sk.SourceName {
				return fmt.Errorf("strict-lock: source mismatch for %q", ref.ID)
			}
			switch {
			case strings.HasPrefix(sk.ResolvedRef, "archive:"):
				if e.ResolvedRef != sk.ResolvedRef {
					return fmt.Errorf("strict-lock: archive ref mismatch for %q", ref.ID)
				}
			case sk.ResolvedRef != "local":
				if e.ResolvedRef != sk.ResolvedRef {
					return fmt.Errorf("strict-lock: resolved ref mismatch for %q", ref.ID)
				}
			}
			if sk.ResolvedRef == "local" && e.LocalDigestSHA != "" {
				digest, err := hashutil.SkillMarkdownDigest(sk.Path)
				if err != nil {
					return err
				}
				if digest != e.LocalDigestSHA {
					return fmt.Errorf("strict-lock: SKILL.md digest mismatch for %q", ref.ID)
				}
			}
		}

		for _, prof := range manifest.EffectiveProfiles(m) {
			for _, t := range prof.Targets {
				destRel := filepath.Join(t.Path, sk.ID)
				destAbs := filepath.Join(projectDir, destRel)
				if opts.DryRun {
					fmt.Fprintf(stdout, "would sync %s (%s) to %s [%s/%s]\n", sk.ID, sk.SourceName, destAbs, prof.Name, t.Kind)
					continue
				}
				if err := atomicCopySkillDir(sk.Path, destAbs); err != nil {
					return fmt.Errorf("sync skill %q: %w", sk.ID, err)
				}
				fmt.Fprintf(stdout, "synced %s → %s\n", sk.ID, destRel)
			}
		}
	}

	if opts.Prune {
		for _, prof := range manifest.EffectiveProfiles(m) {
			for _, t := range prof.Targets {
				base := filepath.Join(projectDir, t.Path)
				if err := pruneUnknown(base, selectedIDsForPrune, stdout, opts.DryRun); err != nil {
					return err
				}
			}
		}
	}

	if err := runHooks(m.Hooks.PostSync, "postSync"); err != nil {
		return err
	}
	return nil
}

func execShell(cmd string, dir string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		c := exec.Command("cmd", "/C", cmd)
		c.Dir = dir
		return c
	}
	c := exec.Command("sh", "-c", cmd)
	c.Dir = dir
	return c
}

func pruneUnknown(dir string, keep map[string]struct{}, stdout io.Writer, dry bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, e := range entries {
		name := e.Name()
		if name == "" || strings.HasPrefix(name, ".") {
			continue
		}
		if _, ok := keep[name]; ok {
			continue
		}
		target := filepath.Join(dir, name)
		if dry {
			fmt.Fprintf(stdout, "would prune %s\n", target)
			continue
		}
		if err := os.RemoveAll(target); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "pruned %s\n", target)
	}
	return nil
}

func atomicCopySkillDir(src string, dest string) error {
	if err := os.RemoveAll(dest); err != nil {
		return err
	}
	parent := filepath.Dir(dest)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	tmp, err := os.MkdirTemp(parent, ".skill-tmp-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	if err := copyDir(src, tmp); err != nil {
		return err
	}
	return os.Rename(tmp, dest)
}
