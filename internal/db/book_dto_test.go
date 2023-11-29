package db

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseTags(t *testing.T) {
	type test struct {
		input string
		want  []string
	}

	tests := []test{
		{input: "", want: []string{}},
		{input: "one,two,three", want: []string{"one", "two", "three"}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%+v", tc.input), func(t *testing.T) {
			got := parseTags(tc.input)
			equal := reflect.DeepEqual(got, tc.want)
			if !equal {
				t.Errorf("got %+v; want %+v", got, tc.want)
			}
		})
	}
}
