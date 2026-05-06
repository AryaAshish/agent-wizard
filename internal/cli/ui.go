package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type ui struct {
	out io.Writer
	tty bool
}

func newUI(out io.Writer) ui {
	stdoutInfo, err := os.Stdout.Stat()
	isTTY := err == nil && (stdoutInfo.Mode()&os.ModeCharDevice) != 0
	return ui{out: out, tty: isTTY}
}

func (u ui) Header(title string) {
	if !u.tty {
		fmt.Fprintln(u.out, title)
		return
	}
	fmt.Fprintln(u.out, "╔══════════════════════════════════════════════╗")
	fmt.Fprintln(u.out, "║                agent-wizard                 ║")
	fmt.Fprintln(u.out, "╠══════════════════════════════════════════════╣")
	fmt.Fprintf(u.out, "║ %-44s ║\n", strings.ToUpper(title))
	fmt.Fprintln(u.out, "╚══════════════════════════════════════════════╝")
}

func (u ui) Section(title string) {
	if u.tty {
		fmt.Fprintf(u.out, "\n-- %s --\n", title)
		return
	}
	fmt.Fprintln(u.out, title)
}

func (u ui) OK(msg string) {
	if u.tty {
		fmt.Fprintf(u.out, "\x1b[32mOK\x1b[0m  %s\n", msg)
		return
	}
	fmt.Fprintf(u.out, "OK  %s\n", msg)
}

func (u ui) Warn(msg string) {
	if u.tty {
		fmt.Fprintf(u.out, "\x1b[33mWARN\x1b[0m  %s\n", msg)
		return
	}
	fmt.Fprintf(u.out, "WARN  %s\n", msg)
}

func (u ui) NextActions(actions ...string) {
	if len(actions) == 0 {
		return
	}
	u.Section("Next actions")
	for _, a := range actions {
		fmt.Fprintf(u.out, "  - %s\n", a)
	}
}

func (u ui) Commands(title string, commands ...string) {
	if len(commands) == 0 {
		return
	}
	u.Section(title)
	if u.tty {
		fmt.Fprintln(u.out, "  copy/paste:")
	}
	for _, c := range commands {
		fmt.Fprintf(u.out, "  %s\n", c)
	}
}
