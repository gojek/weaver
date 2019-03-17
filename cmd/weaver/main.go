package main

import (
	"fmt"
	"github.com/gojektech/weaver/internal/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Weaver"
	app.Version = fmt.Sprintf("%s built on %s (commit: %s)", Version, BuildDate, Commit)
	app.Description = "An Advanced HTTP Reverse Proxy with Dynamic Sharding Strategies"
	app.Usage = "Multi Utility CLI Tool"
	app.Run(os.Args)
}

// Build information (will be injected during build)
var (
	Version   = "1.0.0"
	Commit    = "n/a"
	BuildDate = "n/a"
)
