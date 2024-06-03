package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/alecthomas/kong"
	"github.com/jonathanhope/armaria/cmd/cli/internal"
	"github.com/jonathanhope/armaria/cmd/cli/internal/messaging"
	"github.com/jonathanhope/armaria/internal/manifest"
)

var version string

func main() {
	// When a browser sends a native message it will send the extension ID as the last argument.
	// When we see one of Armaria's extension IDs we switch to native messaging mode.
	extensionIds := []string{manifest.FirefoxExtension, manifest.ChromeExtension1, manifest.ChromeExtension2}
	hostMode := false
	for _, arg := range os.Args {
		if slices.Contains(extensionIds, arg) {
			hostMode = true
		}
	}

	if hostMode {
		if err := messaging.Dispatch(os.Stdin, os.Stdout); err != nil {
			fmt.Printf("Unexpected error: %s", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Otherwise we fall into the standard CLI.
	rootCmd := cmd.RootCmdFactory()
	ctx := kong.Parse(&rootCmd)
	err := ctx.Run(&cmd.Context{
		DB:         rootCmd.DB,
		Formatter:  rootCmd.Formatter,
		Writer:     os.Stdout,
		ReturnCode: os.Exit,
		Version:    version})

	ctx.FatalIfErrorf(err)
}
