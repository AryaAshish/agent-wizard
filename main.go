package main

import (
	"fmt"
	"os"

	"github.com/aryaashish/agent-wizard/internal/buildinfo"
	"github.com/aryaashish/agent-wizard/internal/cli"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v", "version":
			fmt.Printf("agent-wizard %s (commit=%s date=%s)\n", buildinfo.Version, buildinfo.Commit, buildinfo.Date)
			return
		}
	}
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "agent-wizard: %v\n", err)
		code := 1
		if ec, ok := err.(cli.ExitCoder); ok {
			code = ec.Code()
		}
		os.Exit(code)
	}
}
