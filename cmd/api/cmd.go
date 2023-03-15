package main

import (
	"flag"
	"fmt"
)

func showHelp() {
	fmt.Printf("%s. Available Commands and Options:\n\n", AppName)
	flag.PrintDefaults()
}

// showVersion print app version and integrity
func showVersion() {
	fmt.Printf("%s\n"+
		"  Version: %s\n"+
		"  Build:   %s\n", AppName, AppVersion, CommitHash)
}
