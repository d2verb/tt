package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/d2verb/tt"
)

func main() {
	err := tt.Run(os.Args[1:], os.Stdout, os.Stderr)
	if err != nil && err != flag.ErrHelp {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
