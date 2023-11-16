package main

import (
	"log"
	"os"

	"github.com/jonathanhope/armaria/cmd/host/internal/messaging"
)

// Browser extensions can only access local resources (like files) with something called native messaging.
// This is a native messaging host that allows the Armaria extension to communicate with Armaria itself.

func main() {
	if err := messaging.Dispatch(os.Stdin, os.Stdout); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}
