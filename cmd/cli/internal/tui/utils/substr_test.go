package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSubstr(t *testing.T) {
	type test struct {
		input  string
		length int
		want   string
	}

	tests := []test{
		{input: "", length: 0, want: ""},
		{input: "a", length: 2, want: "a"},
		{input: "a", length: 1, want: "a"},
		{input: "aaa", length: 2, want: "a…"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := Substr(tc.input, tc.length)
			diff := cmp.Diff(got, tc.want)
			if diff != "" {
				t.Errorf("Expected and actual strings different:\n%s", diff)
			}
		})
	}
}
