package main

import (
	"os"

	"github.com/preflightsh/preflight/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
