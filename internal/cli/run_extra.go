package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/aryaashish/agent-wizard/internal/cache"
	catalog "github.com/aryaashish/agent-wizard/internal/catalog"
	check "github.com/aryaashish/agent-wizard/internal/ci"
	"github.com/aryaashish/agent-wizard/internal/config"
	"github.com/aryaashish/agent-wizard/internal/drift"
	"github.com/aryaashish/agent-wizard/internal/engine"
	"github.com/aryaashish/agent-wizard/internal/manifest"
	"github.com/aryaashish/agent-wizard/internal/migrate"
	"github.com/aryaashish/agent-wizard/internal/model"
	"github.com/aryaashish/agent-wizard/internal/projectdir"
	"github.com/aryaashish/agent-wizard/internal/skills"
)

func runListExpanded(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var sourcePath string
	var sourceName string
	var filterNeedle string
	var installed bool

	fs.StringVar(&sourcePath, "source", ".", "source directory to scan")
	fs.StringVar(&sourcePath, "S", ".", "source directory to scan (shorthand)")
	fs.StringVar(&sourceName, "source-name", "", "configured source name")
	fs.StringVar(&filterNeedle, "filter", "", "when listing skill ids: include only ids containing substring (case-insensitive)")
	fs.BoolVar(&installed, "installed", false, "list resolved skills from current manifest")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if installed {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		pd, err := projectdir.ResolveForProjectOps(cwd)
		if err != nil {
			if errors.Is(err, projectdir.ErrNoManifest) {
				return fmt.Errorf("%w\nhint: run `agent-wizard add` or `agent-wizard init` first", err)
			}
			return err
		}
		printProjectIfDifferent(stdout, cwd, pd)
		m, err := manifest.Load(pd)
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
		if strings.TrimSpace(filterNeedle) != "" {
			var filtered []model.SkillRef
			for _, r := range refs {
				if skills.IDContainsFold(r.ID, filterNeedle) {
					filtered = append(filtered, r)
				}
			}
			refs = filtered
		}
		buckets, err := engine.BuildBuckets(cfg, m.Sources)
		if err != nil {
			return err
		}
		var rows []listSkillRow
		for _, r := range refs {
			sk, err := engine.ResolveSkill(r, buckets)
			display := r.ID
			if r.SourceAlias != "" {
				display = r.SourceAlias + "/" + r.ID
			}
			desc := "-"
			if err == nil {
				desc = skills.SkillSummaryLine(sk)
			}
			rows = append(rows, listSkillRow{id: display, summary: desc})
		}
		writeSkillTable(stdout, rows)
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
		srcCfg, ok := engine.ResolveSource(cfg, sourceName)
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
	found = skills.FilterSkillsByIDSubstring(found, filterNeedle)
	sort.Slice(found, func(i, j int) bool {
		return found[i].ID < found[j].ID
	})
	var rows []listSkillRow
	for _, s := range found {
		rows = append(rows, listSkillRow{id: s.ID, summary: skills.SkillSummaryLine(s)})
	}
	writeSkillTable(stdout, rows)
	return nil
}

// listSkillRow is one line of `list` tabular output (id and summary from SKILL.md).
type listSkillRow struct {
	id      string
	summary string
}

func writeSkillTable(stdout io.Writer, rows []listSkillRow) {
	if len(rows) == 0 {
		fmt.Fprintln(stdout, "No skills found for this source or filter.")
		fmt.Fprintln(stdout, "Try:")
		fmt.Fprintln(stdout, "  agent-wizard list --source-name community")
		fmt.Fprintln(stdout, "  (widen --filter or fix --source PATH)")
		fmt.Fprintln(stdout, "Create a local scaffold:")
		fmt.Fprintln(stdout, "  agent-wizard create-skill <skill-id>")
		return
	}
	w := tabwriter.NewWriter(stdout, 0, 0, 2, ' ', 0)
	for _, r := range rows {
		fmt.Fprintf(w, "%s\t%s\n", r.id, r.summary)
	}
	_ = w.Flush()
}

const createSkillMarkdownTemplate = `# %[1]s

Replace this paragraph with one or two sentences. This is the first paragraph under the title — it appears as the short summary in agent-wizard list.

## When to use

-

## When not to use

-

## Inputs

-

## Outputs

-

## Steps

1.

## Safety

-
`

var createSkillIDRegexp = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func runCreateSkill(args []string, stdout io.Writer) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: agent-wizard create-skill <skill-id>")
	}
	skillID := strings.TrimSpace(args[0])
	if skillID == "." || skillID == ".." || strings.ContainsAny(skillID, `/\`) {
		return fmt.Errorf("invalid skill id %q", skillID)
	}
	if !createSkillIDRegexp.MatchString(skillID) {
		return fmt.Errorf(`skill id %q is invalid — use lowercase letters, digits, and "-" only (e.g. my-review-checklist)`, skillID)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	pd, err := projectdir.ResolveForInitOrAdd(cwd)
	if err != nil {
		return err
	}
	printProjectIfDifferent(stdout, cwd, pd)
	dir := filepath.Join(pd, skillID)
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("./%s already exists — remove it or choose another id", skillID)
	}
	if err := os.Mkdir(dir, 0o755); err != nil {
		return err
	}
	body := fmt.Sprintf(createSkillMarkdownTemplate, skillID)
	skillPath := filepath.Join(dir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(body), 0o644); err != nil {
		return err
	}
	ui := newUI(stdout)
	ui.OK(fmt.Sprintf("created %s/SKILL.md", skillID))
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Local testing:")
	fmt.Fprintf(stdout, "  cd %s && agent-wizard list --source . --filter %s\n", pd, skillID)
	fmt.Fprintln(stdout, "  Or register that directory as a local source:")
	fmt.Fprintln(stdout, "    agent-wizard sources add --name NAME --kind local --path ABS_PATH_TO_PARENT")
	fmt.Fprintf(stdout, "    agent-wizard add %s --source NAME && agent-wizard sync\n", skillID)
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "To contribute this skill to the bundled library:")
	fmt.Fprintln(stdout, "  • Fork the agent-wizard repository")
	fmt.Fprintf(stdout, "  • Copy this folder to: internal/community/assets/%s/\n", skillID)
	fmt.Fprintln(stdout, "  • Open a pull request")
	fmt.Fprintln(stdout, "  Guide: CONTRIBUTING.md — structure and quality bar.")
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

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	pd, err := projectdir.ResolveForProjectOps(cwd)
	if err != nil {
		if errors.Is(err, projectdir.ErrNoManifest) {
			return fmt.Errorf("%w\nhint: run from your project tree", err)
		}
		return err
	}
	printProjectIfDifferent(stdout, cwd, pd)
	m, err := manifest.Load(pd)
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
			"manifest":          manifest.PathFromDir(pd),
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

	fmt.Fprintf(stdout, "manifest: %s\n", manifest.PathFromDir(pd))
	fmt.Fprintf(stdout, "targetDir: %s\n", m.TargetDir)
	fmt.Fprintf(stdout, "installMode: %s\n", m.InstallMode)
	fmt.Fprintf(stdout, "schemaVersion: %d\n", m.SchemaVersion)
	fmt.Fprintf(stdout, "sources: %s\n", strings.Join(m.Sources, ","))
	fmt.Fprintf(stdout, "packs: %s\n", strings.Join(m.Packs, ","))
	fmt.Fprintf(stdout, "skills: %s\n", strings.Join(m.Skills, ","))
	fmt.Fprintf(stdout, "configuredSources: %d\n", len(cfg.Sources))

	if checkDrifts {
		msgs, ok, err := drift.Evaluate(pd, m, cfg, strictDigest)
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
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	pd, err := projectdir.ResolveForProjectOps(cwd)
	if err != nil {
		if errors.Is(err, projectdir.ErrNoManifest) {
			return fmt.Errorf("%w\nhint: run from your project tree", err)
		}
		return err
	}
	printProjectIfDifferent(stdout, cwd, pd)
	m, err := manifest.Load(pd)
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
	if err := engine.GenerateLockfile(pd, m, cfg); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "wrote %s\n", filepath.Join(pd, "agentskills.lock"))
	return nil
}

func runMigrate(stdout io.Writer) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	pd, err := projectdir.ResolveForProjectOps(cwd)
	if err != nil {
		if errors.Is(err, projectdir.ErrNoManifest) {
			return fmt.Errorf("%w\nhint: run from your project tree", err)
		}
		return err
	}
	printProjectIfDifferent(stdout, cwd, pd)
	if err := migrate.Run(pd); err != nil {
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
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	pd, err := projectdir.ResolveForProjectOps(cwd)
	if err != nil {
		if errors.Is(err, projectdir.ErrNoManifest) {
			return fmt.Errorf("%w\nhint: run from your project tree", err)
		}
		return err
	}
	printProjectIfDifferent(stdout, cwd, pd)
	m, err := manifest.Load(pd)
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
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	pd, err := projectdir.ResolveForProjectOps(cwd)
	if err != nil {
		if errors.Is(err, projectdir.ErrNoManifest) {
			return fmt.Errorf("%w\nhint: run from your project tree", err)
		}
		return err
	}
	printProjectIfDifferent(stdout, cwd, pd)
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
		m, err := manifest.Load(pd)
		if err != nil {
			return err
		}
		return engine.Sync(pd, m, cfg, stdout, engine.SyncOpts{})
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
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	pd, err := projectdir.ResolveForProjectOps(cwd)
	if err != nil {
		if errors.Is(err, projectdir.ErrNoManifest) {
			return fmt.Errorf("%w\nhint: run from your project tree", err)
		}
		return err
	}
	printProjectIfDifferent(stdout, cwd, pd)
	m, err := manifest.Load(pd)
	if err != nil {
		return err
	}
	m.Packs = engine.AddUnique(m.Packs, args[0])
	if err := manifest.Save(pd, m); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "added pack %s\n", args[0])
	return nil
}
