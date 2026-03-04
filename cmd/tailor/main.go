package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

var version = "dev"

var cli struct {
	Version bool `help:"Show version."`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("tailor"),
		kong.Description("Bespoke project templates for GitHub repositories."),
		kong.UsageOnError(),
		kong.Vars{"version": version},
	)

	if cli.Version {
		fmt.Printf("tailor %s\n", version)
		os.Exit(0)
	}

	ctx.PrintUsage(false)
}
