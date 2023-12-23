package utils

import (
	"strings"
)

// Substr cuts a string down to substring.
// if the strings length is reduced ellipsis are added.
func Substr(s string, from, length int) string {
	if s == "" {
		return s
	}

	wb := strings.Split(s, "")

	to := from + length
	if to > len(wb) {
		to = len(wb)
	}

	if from > len(wb) {
		from = len(wb)
	}

	if to < len(wb) {
		to -= 1
	}

	out := strings.Join(wb[from:to], "")
	if s == out {
		return s
	}

	if out != "" {
		return out + "…"
	} else {
		return out
	}
}
