package main

import (
	"fmt"
)

// This is injected by CI
var (
	Build  = "SELFBUILD" // Will be the date/time for the build in UTC
	Commit = "N/A"       // Commit hash
)

// Version - Weaver Version
var Version = map[string]string{
	"MAJOR":  "0",
	"MINOR":  "0",
	"PATCH":  "1",
	"BUILD":  Build,
	"COMMIT": Commit,
}

// GenVersion - generates Weaver Version string
func GenVersion() string {
	return fmt.Sprintf("Weaver v%s.%s.%s %s %s",
		Version["MAJOR"],
		Version["MINOR"],
		Version["PATCH"],
		Version["COMMIT"],
		Version["BUILD"],
	)
}
