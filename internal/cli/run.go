package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

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

	switch args[0] {
	case "--help", "-h", "help":
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
		if len(args) < 2 || args[1] != "add" {
			return fmt.Errorf("unknown pack command (try: pack add <id>)")
		}
		return runPackAdd(args[2:], stdout)
	case "catalog":
		if len(args) < 2 || args[1] != "validate" {
			return fmt.Errorf("unknown catalog command (try: catalog validate <file>)")
		}
		return runCatalogValidate(args[2:], stdout)
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
	fmt.Fprintln(stdout, "  init     Initialize agentskills.yaml in current project")
	fmt.Fprintln(stdout, "  list     List discovered skills from a source path")
	fmt.Fprintln(stdout, "  add      Add skill to project manifest")
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
	fmt.Fprintln(stdout, "  icp      Set or validate target ICP")
	fmt.Fprintln(stdout, "  help     Show this help message")
}

func runInit(stdout io.Writer) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	m, err := manifest.Init(wd)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "initialized %s with targetDir=%s\n", manifest.FileName, m.TargetDir)
	return nil
}

func runAdd(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("add requires a skill id")
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	m, err := manifest.Load(wd)
	if err != nil {
		return err
	}
	m.Skills = engine.AddUnique(m.Skills, args[0])
	if err := manifest.Save(wd, m); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "added %s\n", args[0])
	return nil
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
	return engine.Sync(wd, m, cfg, stdout, opts)
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
		fs.StringVar(&name, "name", "", "source name")
		fs.StringVar(&kind, "kind", "local", "source kind")
		fs.StringVar(&path, "path", "", "source path")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if name == "" || path == "" {
			return fmt.Errorf("sources add requires --name and --path")
		}
		if _, ok := cfg.GetSource(name); ok {
			return fmt.Errorf("source %q already exists", name)
		}
		cfg.Sources = append(cfg.Sources, config.Source{Name: name, Kind: kind, Path: path})
		if err := config.Save(cfgPath, cfg); err != nil {
			return err
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
