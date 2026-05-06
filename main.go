package main

import (
	"fmt"
	"os"

	"github.com/aryaashish/agent-wizard/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "agent-wizard: %v\n", err)
		code := 1
		if ec, ok := err.(cli.ExitCoder); ok {
			code = ec.Code()
		}
		os.Exit(code)
	}
}
