package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/jonathanhope/armaria/cmd/cli/internal"
)

var version string

func main() {
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
