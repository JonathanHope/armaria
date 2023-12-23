package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSubstr(t *testing.T) {
	type test struct {
		input  string
		from   int
		length int
		want   string
	}

	tests := []test{
		{input: "", from: 0, length: 0, want: ""},
		{input: "a", from: 2, length: 0, want: ""},
		{input: "a", from: 0, length: 2, want: "a"},
		{input: "a", from: 0, length: 1, want: "a"},
		{input: "aaa", from: 0, length: 2, want: "a…"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := Substr(tc.input, tc.from, tc.length)
			diff := cmp.Diff(got, tc.want)
			if diff != "" {
				t.Errorf("Expected and actual strings different:\n%s", diff)
			}
		})
	}
}
