package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/jonathanhope/armaria/cmd/cli/cmd"
)

func main() {
	rootCmd := cmd.RootCmdFactory()
	ctx := kong.Parse(&rootCmd)

	err := ctx.Run(&cmd.Context{
		DB:         rootCmd.DB,
		Formatter:  rootCmd.Formatter,
		Writer:     os.Stdout,
		ReturnCode: os.Exit})

	ctx.FatalIfErrorf(err)
}
