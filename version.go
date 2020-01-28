package main

import (
	"flag"
	"fmt"
)

// Version of the app
var Version = "snapshot"

var (
	version = flag.Bool("version", false, "Print version")
)

func printVersion() {
	fmt.Printf(`webhookd (%s)
Copyright (C) 2020 Nunux, Org.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Written by Nicolas Carlier.`, Version)
}
