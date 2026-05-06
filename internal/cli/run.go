package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aryaashish/agent-wizard/internal/community"
	"github.com/aryaashish/agent-wizard/internal/config"
	"github.com/aryaashish/agent-wizard/internal/engine"
	"github.com/aryaashish/agent-wizard/internal/manifest"
)

func Run(args []string) error {
	return run(args, os.Stdout)
}

func run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		printHelp(stdout)
		return nil
	}

	if len(args) > 1 && isHelpFlag(args[1]) {
		if printCommandHelp(args[0], stdout) {
			return nil
		}
	}

	switch args[0] {
	case "--help", "-h", "help":
		if len(args) > 1 {
			if printCommandHelp(args[1], stdout) {
				return nil
			}
			return fmt.Errorf("unknown help topic %q", args[1])
		}
		printHelp(stdout)
		return nil
	case "init":
		return runInit(stdout)
	case "list":
		return runListExpanded(args[1:], stdout)
	case "add":
		return runAdd(args[1:], stdout)
	case "remove":
		return runRemove(args[1:], stdout)
	case "status":
		return runStatusExpanded(args[1:], stdout)
	case "sync":
		return runSync(args[1:], stdout)
	case "sources":
		return runSources(args[1:], stdout)
	case "icp":
		return runICP(args[1:], stdout)
	case "lock":
		return runLock(stdout)
	case "migrate":
		return runMigrate(stdout)
	case "cache":
		return runCache(args[1:], stdout)
	case "ci-check":
		return runCICheck(stdout)
	case "browse":
		return runBrowse(args[1:], stdout)
	case "watch":
		return runWatch(args[1:], stdout)
	case "import":
		return runImport(args[1:], stdout)
	case "pack":
		if len(args) >= 3 && args[1] == "add" && isHelpFlag(args[2]) {
			printCommandHelp("pack add", stdout)
			return nil
		}
		if len(args) < 2 || args[1] != "add" {
			return fmt.Errorf("unknown pack command (try: pack add <id>)")
		}
		return runPackAdd(args[2:], stdout)
	case "catalog":
		if len(args) < 2 || args[1] != "validate" {
			return fmt.Errorf("unknown catalog command (try: catalog validate <file>)")
		}
		return runCatalogValidate(args[2:], stdout)
	case "community":
		if len(args) < 2 || args[1] != "fetch" {
			return fmt.Errorf("unknown community command (try: community fetch)")
		}
		return runCommunityFetch(stdout)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func printHelp(stdout io.Writer) {
	fmt.Fprintln(stdout, "agent-wizard - manage reusable agent skills")
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "Usage:")
	fmt.Fprintln(stdout, "  agent-wizard <command> [flags]")
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "Commands:")
	fmt.Fprintln(stdout, "  init     Initialize project and launch community starter picker")
	fmt.Fprintln(stdout, "  list     List available skills from a source")
	fmt.Fprintln(stdout, "  add      Add skill to manifest (supports source qualifier)")
	fmt.Fprintln(stdout, "  remove   Remove skill from project manifest")
	fmt.Fprintln(stdout, "  status   Show manifest/source status (+ --json/--check-drifts)")
	fmt.Fprintln(stdout, "  sync     Sync selected skills to target dir")
	fmt.Fprintln(stdout, "  sources  Manage local source registry")
	fmt.Fprintln(stdout, "  lock     Write agentskills.lock pins")
	fmt.Fprintln(stdout, "  migrate  Backup + normalize manifest defaults")
	fmt.Fprintln(stdout, "  cache    Cache maintenance (status|prune)")
	fmt.Fprintln(stdout, "  ci-check Optional env policy gates for CI")
	fmt.Fprintln(stdout, "  browse   Minimal interactive picker (stdin)")
	fmt.Fprintln(stdout, "  watch    Poll-based auto sync loop")
	fmt.Fprintln(stdout, "  import   Copy discovered skills tree into library")
	fmt.Fprintln(stdout, "  pack     Pack helpers (pack add)")
	fmt.Fprintln(stdout, "  catalog  Curated index validation")
	fmt.Fprintln(stdout, "  community Refresh bundled community starter assets (community fetch)")
	fmt.Fprintln(stdout, "  icp      Set or validate target ICP")
	fmt.Fprintln(stdout, "  version  Show CLI version/build info")
	fmt.Fprintln(stdout, "  help     Show this help message")
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "Tip:")
	fmt.Fprintln(stdout, "  agent-wizard <command> --help")
	fmt.Fprintln(stdout, "  agent-wizard help <command>")
}

func runInit(stdout io.Writer) error {
	ui := newUI(stdout)
	ui.Header("project init")
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	var m manifest.Manifest
	m, err = manifest.Init(wd)
	if err != nil {
		if strings.Contains(err.Error(), "manifest already exists") {
			m, err = manifest.Load(wd)
			if err != nil {
				return err
			}
			ui.Warn("manifest already exists, reusing current config")
		} else {
			return err
		}
	}
	if err := ensureCommunitySourceConfigured(); err != nil {
		return err
	}
	m.Sources = engine.AddUnique(m.Sources, community.SourceName)
	if err := manifest.Save(wd, m); err != nil {
		return err
	}
	ui.OK(fmt.Sprintf("initialized %s with targetDir=%s", manifest.FileName, m.TargetDir))
	if !isInteractiveTerminal() {
		ui.Warn("non-interactive terminal detected; skipping picker")
		ui.Commands(
			"Browse community skills",
			"agent-wizard list --source-name community",
		)
		ui.Commands(
			"Install common starter skills",
			"agent-wizard add pr-review --source community",
			"agent-wizard pack add android-starter",
			"agent-wizard sync",
		)
		ui.NextActions("agent-wizard status", "agent-wizard list --installed")
		return nil
	}
	selection, err := runInitPicker(stdout)
	if err != nil {
		return err
	}
	installed := make([]string, 0, 4)
	for _, p := range selection.packs {
		m.Packs = engine.AddUnique(m.Packs, p)
		installed = append(installed, "pack:"+p)
	}
	for _, s := range selection.skills {
		m.Skills = engine.AddUnique(m.Skills, community.SourceName+"/"+s)
		installed = append(installed, "skill:"+community.SourceName+"/"+s)
	}
	if err := manifest.Save(wd, m); err != nil {
		return err
	}
	if err := runSync([]string{}, stdout); err != nil {
		return err
	}
	ui.Section("Summary")
	if len(installed) == 0 {
		fmt.Fprintln(stdout, "No skills selected; source configured only.")
	} else {
		fmt.Fprintf(stdout, "Installed: %s\n", strings.Join(installed, ", "))
	}
	fmt.Fprintf(stdout, "Target: %s\n", m.TargetDir)
	ui.Commands(
		"Browse community skills",
		"agent-wizard list --source-name community",
	)
	ui.Commands(
		"Install common starter skills",
		"agent-wizard add pr-review --source community",
		"agent-wizard pack add android-starter",
		"agent-wizard sync",
	)
	ui.NextActions("agent-wizard status", "agent-wizard list --installed")
	return nil
}

func runAdd(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("add requires a skill id")
	}
	skillID := args[0]
	sourceShorthand := ""
	remaining := args[1:]
	if len(remaining) > 0 && strings.HasPrefix(remaining[0], "-") && !strings.HasPrefix(remaining[0], "--") && remaining[0] != "-h" {
		sourceShorthand = strings.TrimPrefix(remaining[0], "-")
		remaining = remaining[1:]
	}
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var sourceName string
	fs.StringVar(&sourceName, "source", "", "source alias used to qualify skill id")
	fs.StringVar(&sourceName, "s", "", "source alias used to qualify skill id (shorthand)")
	if err := fs.Parse(remaining); err != nil {
		return err
	}
	if sourceShorthand != "" && sourceName != "" {
		return fmt.Errorf("use either shorthand -<source> or --source, not both")
	}
	if sourceName == "" {
		sourceName = sourceShorthand
	}
	if sourceName != "" && strings.Contains(skillID, "/") {
		return fmt.Errorf("skill %q is already source-qualified", skillID)
	}
	if sourceName != "" {
		skillID = sourceName + "/" + skillID
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	m, err := manifest.Load(wd)
	if err != nil {
		return err
	}
	m.Skills = engine.AddUnique(m.Skills, skillID)
	if err := manifest.Save(wd, m); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "added %s\n", skillID)
	return nil
}

func isHelpFlag(arg string) bool {
	return arg == "--help" || arg == "-h"
}

func printCommandHelp(command string, stdout io.Writer) bool {
	switch command {
	case "list":
		fmt.Fprintln(stdout, "Usage: agent-wizard list [--source PATH | --source-name NAME | --installed]")
		fmt.Fprintln(stdout, "Examples:")
		fmt.Fprintln(stdout, "  agent-wizard list --source ./examples/library")
		fmt.Fprintln(stdout, "  agent-wizard list --source-name community")
		fmt.Fprintln(stdout, "  agent-wizard list --installed")
		return true
	case "init":
		fmt.Fprintln(stdout, "Usage: agent-wizard init")
		fmt.Fprintln(stdout, "Creates manifest, wires community source, and shows starter picker.")
		fmt.Fprintln(stdout, "After init, browse skills with:")
		fmt.Fprintln(stdout, "  agent-wizard list --source-name community")
		return true
	case "add":
		fmt.Fprintln(stdout, "Usage: agent-wizard add <skill-id> [--source NAME]")
		fmt.Fprintln(stdout, "Shortcut:")
		fmt.Fprintln(stdout, "  agent-wizard add <skill-id> -<source>")
		fmt.Fprintln(stdout, "Examples:")
		fmt.Fprintln(stdout, "  agent-wizard add pr-review --source android")
		fmt.Fprintln(stdout, "  agent-wizard add pr-review -android")
		return true
	case "sources":
		fmt.Fprintln(stdout, "Usage:")
		fmt.Fprintln(stdout, "  agent-wizard sources list")
		fmt.Fprintln(stdout, "  agent-wizard sources add --name NAME --kind local|git|archive [--path PATH]")
		fmt.Fprintln(stdout, "    git extras: --git-url URL [--git-ref REF] [--subdir DIR]")
		fmt.Fprintln(stdout, "    archive extras: --archive-url URL")
		fmt.Fprintln(stdout, "    optional: --quiet")
		fmt.Fprintln(stdout, "  agent-wizard sources remove NAME")
		return true
	case "pack", "pack add":
		fmt.Fprintln(stdout, "Usage: agent-wizard pack add <pack-id>")
		fmt.Fprintln(stdout, "Example:")
		fmt.Fprintln(stdout, "  agent-wizard pack add android-starter")
		return true
	case "sync":
		fmt.Fprintln(stdout, "Usage: agent-wizard sync [--dry-run|--check] [--prune] [--strict-lock]")
		return true
	case "status":
		fmt.Fprintln(stdout, "Usage: agent-wizard status [--json] [--check-drifts] [--strict-digest]")
		return true
	case "community":
		fmt.Fprintln(stdout, "Usage: agent-wizard community fetch")
		return true
	default:
		return false
	}
}

func runRemove(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("remove requires a skill id")
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	m, err := manifest.Load(wd)
	if err != nil {
		return err
	}
	m.Skills = engine.RemoveValue(m.Skills, args[0])
	if err := manifest.Save(wd, m); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "removed %s\n", args[0])
	return nil
}

func runSync(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var dryRun bool
	var prune bool
	var strictLock bool
	fs.BoolVar(&dryRun, "dry-run", false, "show what would sync without writing")
	fs.BoolVar(&dryRun, "check", false, "alias for dry-run")
	fs.BoolVar(&prune, "prune", false, "remove synced skill dirs not in manifest (destructive)")
	fs.BoolVar(&strictLock, "strict-lock", false, "enforce agentskills.lock pins when present")
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
	opts := engine.SyncOpts{DryRun: dryRun, Prune: prune, StrictLock: strictLock}
	if err := engine.Sync(wd, m, cfg, stdout, opts); err != nil {
		return err
	}
	if !dryRun {
		newUI(stdout).NextActions("agent-wizard status", "agent-wizard list --installed")
	}
	return nil
}

func runSources(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("sources requires subcommand: list|add|remove")
	}
	cfgPath, err := config.DefaultPath()
	if err != nil {
		return err
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}
	switch args[0] {
	case "list":
		for _, s := range cfg.Sources {
			fmt.Fprintf(stdout, "%s\t%s\t%s\n", s.Name, s.Kind, s.Path)
		}
		return nil
	case "add":
		fs := flag.NewFlagSet("sources add", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		var name, kind, path string
		var gitURL, gitRef, subdir, archiveURL string
		var quiet bool
		fs.StringVar(&name, "name", "", "source name")
		fs.StringVar(&kind, "kind", "local", "source kind")
		fs.StringVar(&path, "path", "", "source path")
		fs.StringVar(&gitURL, "git-url", "", "git source URL")
		fs.StringVar(&gitURL, "gitUrl", "", "git source URL (camelCase alias)")
		fs.StringVar(&gitRef, "git-ref", "", "git ref/tag/branch")
		fs.StringVar(&gitRef, "gitRef", "", "git ref/tag/branch (camelCase alias)")
		fs.StringVar(&subdir, "subdir", "", "subdirectory inside git repository")
		fs.StringVar(&archiveURL, "archive-url", "", "zip archive URL")
		fs.StringVar(&archiveURL, "archiveUrl", "", "zip archive URL (camelCase alias)")
		fs.BoolVar(&quiet, "quiet", false, "suppress advisory output")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if name == "" {
			return fmt.Errorf("sources add requires --name")
		}
		switch kind {
		case "local":
			if path == "" {
				return fmt.Errorf("sources add local requires --path")
			}
		case "git":
			if gitURL == "" {
				return fmt.Errorf("sources add git requires --git-url")
			}
		case "archive":
			if archiveURL == "" {
				return fmt.Errorf("sources add archive requires --archive-url")
			}
		default:
			return fmt.Errorf("sources add unsupported --kind %q", kind)
		}
		if _, ok := cfg.GetSource(name); ok {
			return fmt.Errorf("source %q already exists", name)
		}
		cfg.Sources = append(cfg.Sources, config.Source{
			Name:       name,
			Kind:       kind,
			Path:       path,
			GitURL:     gitURL,
			GitRef:     gitRef,
			Subdir:     subdir,
			ArchiveURL: archiveURL,
		})
		if err := config.Save(cfgPath, cfg); err != nil {
			return err
		}
		if kind == "local" && !quiet {
			fmt.Fprintln(stdout, "warning: local paths are machine-specific and not team-shareable; use git/archive sources for team collaboration")
		}
		fmt.Fprintf(stdout, "added source %s\n", name)
		return nil
	case "remove":
		if len(args) < 2 {
			return fmt.Errorf("sources remove requires source name")
		}
		name := args[1]
		out := make([]config.Source, 0, len(cfg.Sources))
		for _, s := range cfg.Sources {
			if s.Name != name {
				out = append(out, s)
			}
		}
		cfg.Sources = out
		if err := config.Save(cfgPath, cfg); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "removed source %s\n", name)
		return nil
	default:
		return fmt.Errorf("unknown sources subcommand %q", args[0])
	}
}

type initSelection struct {
	packs  []string
	skills []string
}

func runInitPicker(stdout io.Writer) (initSelection, error) {
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "Select starter setup:")
	fmt.Fprintln(stdout, "  [1] Install Android starter pack (recommended)")
	fmt.Fprintln(stdout, "  [2] Pick individual skills")
	fmt.Fprintln(stdout, "  [3] Skip for now")
	fmt.Fprint(stdout, "Enter choice (1-3): ")
	sc := bufio.NewScanner(os.Stdin)
	if !sc.Scan() {
		return initSelection{}, sc.Err()
	}
	switch strings.TrimSpace(sc.Text()) {
	case "1":
		return initSelection{packs: []string{"android-starter"}}, nil
	case "2":
		fmt.Fprintln(stdout, "Pick a skill:")
		fmt.Fprintln(stdout, "  [1] pr-review")
		fmt.Fprintln(stdout, "  [2] plan-review")
		fmt.Fprintln(stdout, "  [3] launch-ready")
		fmt.Fprint(stdout, "Enter choice (1-3): ")
		if !sc.Scan() {
			return initSelection{}, sc.Err()
		}
		switch strings.TrimSpace(sc.Text()) {
		case "1":
			return initSelection{skills: []string{"pr-review"}}, nil
		case "2":
			return initSelection{skills: []string{"plan-review"}}, nil
		case "3":
			return initSelection{skills: []string{"launch-ready"}}, nil
		default:
			return initSelection{}, fmt.Errorf("invalid skill selection")
		}
	case "3":
		return initSelection{}, nil
	default:
		return initSelection{}, fmt.Errorf("invalid starter selection")
	}
}

func ensureCommunitySourceConfigured() error {
	cfgPath, err := config.DefaultPath()
	if err != nil {
		return err
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}
	for i := range cfg.Sources {
		if cfg.Sources[i].Name != community.SourceName {
			continue
		}
		if engine.IsLegacyCommunityGitSource(cfg.Sources[i]) {
			cfg.Sources[i] = config.Source{Name: community.SourceName, Kind: community.SourceKind}
			return config.Save(cfgPath, cfg)
		}
		return nil
	}
	cfg.Sources = append(cfg.Sources, config.Source{Name: community.SourceName, Kind: community.SourceKind})
	return config.Save(cfgPath, cfg)
}

func runCommunityFetch(stdout io.Writer) error {
	ui := newUI(stdout)
	if err := ensureCommunitySourceConfigured(); err != nil {
		return err
	}
	ui.Header("community fetch")
	root, err := community.Extract(true)
	if err != nil {
		return err
	}
	ui.OK("community starter assets refreshed")
	fmt.Fprintf(stdout, "path: %s\n", root)
	ui.NextActions("agent-wizard list --source-name community")
	return nil
}

func isInteractiveTerminal() bool {
	stdinInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	stdoutInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (stdinInfo.Mode()&os.ModeCharDevice) != 0 && (stdoutInfo.Mode()&os.ModeCharDevice) != 0
}

func runICP(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("icp requires mode: solo|team|enterprise")
	}
	if err := engine.ValidateICP(args[0]); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "icp validated: %s\n", args[0])
	return nil
}
