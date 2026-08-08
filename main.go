package main

import (
	"os"

	"github.com/ElJoker63/ipa-downloader/v2/cmd"
)

func main() {
	// If CLI subcommands or flags are supplied, execute the CLI frontend.
	// Otherwise, launch the modern Wails desktop application.
	if len(os.Args) > 1 {
		os.Exit(cmd.Execute())
		return
	}

	if err := runDesktopApp(); err != nil {
		println("Error running desktop application: " + err.Error())
		os.Exit(1)
	}
}
