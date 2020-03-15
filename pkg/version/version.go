package version

import (
	"flag"
	"fmt"
)

// Version of the app
var Version = "snapshot"

// GitCommit is the GIT commit revision
var GitCommit = "n/a"

// Built is the built date
var Built = "n/a"

// ShowVersion is the flag used to print version
var ShowVersion = flag.Bool("version", false, "Print version")

// Print version to stdout
func Print() {
	fmt.Printf(`Version:    %s
Git commit: %s
Built:      %s

Copyright (C) 2020 Nicolas Carlier
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
`, Version, GitCommit, Built)
}
