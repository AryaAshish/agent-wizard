package cli

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aryaashish/agent-wizard/internal/cache"
	catalog "github.com/aryaashish/agent-wizard/internal/catalog"
	check "github.com/aryaashish/agent-wizard/internal/ci"
	"github.com/aryaashish/agent-wizard/internal/config"
	"github.com/aryaashish/agent-wizard/internal/drift"
	"github.com/aryaashish/agent-wizard/internal/engine"
	"github.com/aryaashish/agent-wizard/internal/manifest"
	"github.com/aryaashish/agent-wizard/internal/migrate"
	"github.com/aryaashish/agent-wizard/internal/skills"
)

func runListExpanded(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var sourcePath string
	var sourceName string
	var installed bool

	fs.StringVar(&sourcePath, "source", ".", "source directory to scan")
	fs.StringVar(&sourcePath, "S", ".", "source directory to scan (shorthand)")
	fs.StringVar(&sourceName, "source-name", "", "configured source name")
	fs.BoolVar(&installed, "installed", false, "list resolved skills from current manifest")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if installed {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		m, err := manifest.Load(wd)
		if err != nil {
			return err
		}
		cfgPath, err := config.DefaultPath()
		if err != nil {
			return err
		}
		cfg, err := config.Load(cfgPath)
		if err != nil {
			return err
		}
		libRoot, err := engine.LibraryRoot(cfg, m.Sources)
		if err != nil {
			return err
		}
		refs, err := engine.ExpandSkillSelections(m, cfg, libRoot)
		if err != nil {
			return err
		}
		for _, r := range refs {
			if r.SourceAlias == "" {
				fmt.Fprintln(stdout, r.ID)
				continue
			}
			fmt.Fprintf(stdout, "%s/%s\n", r.SourceAlias, r.ID)
		}
		return nil
	}

	if sourceName != "" {
		cfgPath, err := config.DefaultPath()
		if err != nil {
			return err
		}
		cfg, err := config.Load(cfgPath)
		if err != nil {
			return err
		}
		srcCfg, ok := cfg.GetSource(sourceName)
		if !ok {
			return fmt.Errorf("source %q not found", sourceName)
		}
		root, _, err := engine.MaterializeSource(srcCfg)
		if err != nil {
			return err
		}
		sourcePath = root
	}

	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return err
	}

	found, err := skills.NewLocalPathSource(absPath).Discover()
	if err != nil {
		return err
	}
	for _, s := range found {
		fmt.Fprintln(stdout, s.ID)
	}
	return nil
}

func runStatusExpanded(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var jsonOut bool
	var checkDrifts bool
	var strictDigest bool
	fs.BoolVar(&jsonOut, "json", false, "emit JSON")
	fs.BoolVar(&checkDrifts, "check-drifts", false, "compare agentskills.lock vs live sources")
	fs.BoolVar(&strictDigest, "strict-digest", false, "with --check-drifts, require SKILL.md digest match for local pins")
	if err := fs.Parse(args); err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	m, err := manifest.Load(wd)
	if err != nil {
		return err
	}
	cfgPath, err := config.DefaultPath()
	if err != nil {
		return err
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}

	if jsonOut {
		payload := map[string]any{
			"manifest":          manifest.PathFromDir(wd),
			"targetDir":         m.TargetDir,
			"installMode":       m.InstallMode,
			"sources":           m.Sources,
			"skills":            m.Skills,
			"packs":             m.Packs,
			"profiles":          m.Profiles,
			"schemaVersion":     m.SchemaVersion,
			"configuredSources": len(cfg.Sources),
		}
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(payload)
	}

	fmt.Fprintf(stdout, "manifest: %s\n", manifest.PathFromDir(wd))
	fmt.Fprintf(stdout, "targetDir: %s\n", m.TargetDir)
	fmt.Fprintf(stdout, "installMode: %s\n", m.InstallMode)
	fmt.Fprintf(stdout, "schemaVersion: %d\n", m.SchemaVersion)
	fmt.Fprintf(stdout, "sources: %s\n", strings.Join(m.Sources, ","))
	fmt.Fprintf(stdout, "packs: %s\n", strings.Join(m.Packs, ","))
	fmt.Fprintf(stdout, "skills: %s\n", strings.Join(m.Skills, ","))
	fmt.Fprintf(stdout, "configuredSources: %d\n", len(cfg.Sources))

	if checkDrifts {
		msgs, ok, err := drift.Evaluate(wd, m, cfg, strictDigest)
		if err != nil {
			return err
		}
		for _, msg := range msgs {
			fmt.Fprintln(stdout, "drift:", msg)
		}
		if !ok {
			return NewExit(3, fmt.Errorf("lockfile drift detected (%d issue(s))", len(msgs)))
		}
	}
	return nil
}

func runLock(stdout io.Writer) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	m, err := manifest.Load(wd)
	if err != nil {
		return err
	}
	cfgPath, err := config.DefaultPath()
	if err != nil {
		return err
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}
	if err := engine.GenerateLockfile(wd, m, cfg); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "wrote %s\n", filepath.Join(wd, "agentskills.lock"))
	return nil
}

func runMigrate(stdout io.Writer) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := migrate.Run(wd); err != nil {
		return err
	}
	fmt.Fprintln(stdout, "migrate completed (backup saved alongside manifest)")
	return nil
}

func runCache(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("cache requires subcommand: prune|status")
	}
	root, err := cache.Root()
	if err != nil {
		return err
	}
	switch args[0] {
	case "status":
		fmt.Fprintf(stdout, "cache: %s\n", root)
		return nil
	case "prune":
		if err := os.RemoveAll(filepath.Join(root, "git")); err != nil && !os.IsNotExist(err) {
			return err
		}
		if err := os.RemoveAll(filepath.Join(root, "archives")); err != nil && !os.IsNotExist(err) {
			return err
		}
		fmt.Fprintln(stdout, "cache pruned git/ and archives/")
		return nil
	default:
		return fmt.Errorf("unknown cache command %q", args[0])
	}
}

func runCICheck(stdout io.Writer) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	m, err := manifest.Load(wd)
	if err != nil {
		return err
	}
	if err := check.Check(m); err != nil {
		return err
	}
	fmt.Fprintln(stdout, "ci policy check passed")
	return nil
}

func runCatalogValidate(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("catalog validate requires path")
	}
	if err := catalog.ValidateFile(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "catalog OK: %s\n", args[0])
	return nil
}

func runBrowse(args []string, stdout io.Writer) error {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	found, err := skills.NewLocalPathSource(abs).Discover()
	if err != nil {
		return err
	}
	if len(found) == 0 {
		return fmt.Errorf("no skills in %s", abs)
	}
	fmt.Fprintln(stdout, "skills:")
	for i, s := range found {
		fmt.Fprintf(stdout, "  [%d] %s\n", i+1, s.ID)
	}
	fmt.Fprint(stdout, "pick number (enter): ")
	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		return sc.Err()
	}
	line := strings.TrimSpace(sc.Text())
	var idx int
	if _, err := fmt.Sscanf(line, "%d", &idx); err != nil || idx < 1 || idx > len(found) {
		return fmt.Errorf("invalid selection")
	}
	fmt.Fprintf(stdout, "selected: %s\n", found[idx-1].ID)
	return nil
}

func runWatch(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("watch", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var interval time.Duration
	fs.DurationVar(&interval, "interval", 2*time.Second, "polling interval")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if interval < 250*time.Millisecond {
		interval = 250 * time.Millisecond
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	cfgPath, err := config.DefaultPath()
	if err != nil {
		return err
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "watch syncing every %s (Ctrl+C stops)\n", interval)
	t := time.NewTicker(interval)
	defer t.Stop()
	runOnce := func() error {
		m, err := manifest.Load(wd)
		if err != nil {
			return err
		}
		return engine.Sync(wd, m, cfg, stdout, engine.SyncOpts{})
	}
	for {
		if err := runOnce(); err != nil {
			fmt.Fprintf(stdout, "watch sync error: %v\n", err)
		}
		<-t.C
	}
}

func runImport(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("import", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var from string
	var into string
	fs.StringVar(&from, "from", "", "source tree to scan")
	fs.StringVar(&into, "into", "", "destination library folder")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if from == "" || into == "" {
		return fmt.Errorf("import requires --from and --into")
	}
	absFrom, err := filepath.Abs(from)
	if err != nil {
		return err
	}
	absInto, err := filepath.Abs(into)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(absInto, 0o755); err != nil {
		return err
	}
	found, err := skills.NewLocalPathSource(absFrom).Discover()
	if err != nil {
		return err
	}
	for _, s := range found {
		destDir := filepath.Join(absInto, s.ID)
		if err := os.RemoveAll(destDir); err != nil {
			return err
		}
		if err := copyImportTree(s.Path, destDir); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "imported %s\n", s.ID)
	}
	return nil
}

func copyImportTree(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, b, info.Mode())
	})
}

func runPackAdd(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("pack add requires pack id")
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	m, err := manifest.Load(wd)
	if err != nil {
		return err
	}
	m.Packs = engine.AddUnique(m.Packs, args[0])
	if err := manifest.Save(wd, m); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "added pack %s\n", args[0])
	return nil
}
