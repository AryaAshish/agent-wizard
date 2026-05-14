package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func runWizard(stdout io.Writer, stdin io.Reader) error {
	return runWizardInteractive(stdout, stdin, true)
}

func runWizardInteractive(stdout io.Writer, stdin io.Reader, requireTTY bool) error {
	if requireTTY && !isInteractiveTerminal() {
		fmt.Fprintln(stdout, "wizard requires an interactive terminal (TTY). Example: agent-wizard add pr-review --source community")
		return nil
	}
	sc := bufio.NewScanner(stdin)
	for {
		fmt.Fprintln(stdout, "")
		fmt.Fprintln(stdout, "1) Install a community skill")
		fmt.Fprintln(stdout, "2) Add team repo (Git)")
		fmt.Fprintln(stdout, "3) Exit")
		fmt.Fprint(stdout, "Choose [1-3]: ")
		if !sc.Scan() {
			if err := sc.Err(); err != nil {
				return err
			}
			return nil
		}
		line := strings.TrimSpace(strings.ToLower(sc.Text()))
		if line == "" || line == "q" || line == "3" {
			return nil
		}
		switch line {
		case "1":
			wizardInstallLoop(stdout, sc)
		case "2":
			wizardAddTeamRepo(stdout, sc)
		default:
			fmt.Fprintln(stdout, "Enter 1, 2, 3, q, or empty to exit.")
		}
	}
}

func wizardInstallLoop(stdout io.Writer, sc *bufio.Scanner) {
	for {
		fmt.Fprint(stdout, "Install a skill (press Enter for pr-review): ")
		if !sc.Scan() {
			return
		}
		skill := strings.TrimSpace(sc.Text())
		if skill == "" {
			skill = "pr-review"
		}
		if err := runAdd([]string{skill, "--source", "community"}, stdout); err != nil {
			fmt.Fprintf(stdout, "error: %v\n", err)
			fmt.Fprintln(stdout, "hint: agent-wizard list --source-name community")
			fmt.Fprint(stdout, "Press Enter to continue... ")
			_ = sc.Scan()
			return
		}
		fmt.Fprint(stdout, "Install another? (y/n): ")
		if !sc.Scan() {
			return
		}
		a := strings.TrimSpace(strings.ToLower(sc.Text()))
		if a != "y" && a != "yes" {
			return
		}
	}
}

func wizardAddTeamRepo(stdout io.Writer, sc *bufio.Scanner) {
	fmt.Fprint(stdout, "Source name (e.g. my-team): ")
	if !sc.Scan() {
		return
	}
	name := strings.TrimSpace(sc.Text())
	if name == "" {
		fmt.Fprintln(stdout, "name required")
		return
	}
	fmt.Fprint(stdout, "Git URL: ")
	if !sc.Scan() {
		return
	}
	url := strings.TrimSpace(sc.Text())
	if url == "" {
		fmt.Fprintln(stdout, "URL required")
		return
	}
	if err := runSources([]string{"add", "--name", name, "--kind", "git", "--git-url", url}, stdout); err != nil {
		fmt.Fprintf(stdout, "error: %v\n", err)
		fmt.Fprint(stdout, "Press Enter to continue... ")
		_ = sc.Scan()
	}
}
