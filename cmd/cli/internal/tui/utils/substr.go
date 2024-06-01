package utils

import (
	"strings"

	"github.com/muesli/reflow/ansi"
)

// Substr cuts a string down to substring.
// if the strings length is reduced ellipsis are added.
func Substr(s string, length int) string {
	// If the input string is empty or the target length is too short return an empty string.
	if s == "" || length <= 0 {
		return s
	}

	// Measure the actual width of the string.
	width := ansi.PrintableRuneWidth(s)

	// If the width of the string fits in the max length then return it unchanged.
	if length >= width {
		return s
	}

	// Account for the ellipsis in the max length.
	if length < width {
		length -= 1
	}

	// Trim the the string down to the desired length.
	// Start with the obvious case of every char being single width.
	// If that doesn't work remove one char at a time until the width is reached.
	runes := strings.Split(s, "")
	if length > len(runes) {
		length = len(runes)
	}

	out := strings.Join(runes[0:length], "")
	additional := 0
	for ansi.PrintableRuneWidth(out) > length {
		additional += 1
		out = strings.Join(runes[0:length-additional], "")
	}

	return out + "…"
}
