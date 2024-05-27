package cmd

import (
	"io"
)

// Context is the context for an invocation of Armaria.
type Context struct {
	DB         *string   // bookmarks database to use
	Formatter  Formatter // how to format the output
	Writer     io.Writer // where to write output
	ReturnCode func(int) // set the return code
	Version    string    // the current version of Armaria
}
